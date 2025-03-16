package handlers

import (
	"diplomIshi/internal/models"
	"diplomIshi/internal/repository"
	"diplomIshi/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strconv"
	"time"
)

type UserHandler struct {
	repo        *repository.UserRepository
	calcSvc     *services.CalculatorService
	reminderSvc *services.ReminderService
}

func NewUserHandler(repo *repository.UserRepository, calcSvc *services.CalculatorService, reminderSvc *services.ReminderService) *UserHandler {
	return &UserHandler{repo: repo, calcSvc: calcSvc, reminderSvc: reminderSvc}
}

func (h *UserHandler) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	bmi := h.calcSvc.CalculateBMI(user.Weight, user.Height)
	bmr := h.calcSvc.CalculateBMR(user.Weight, user.Height, user.Age, user.Gender)
	tdee := h.calcSvc.CalculateTDEE(bmr, user.ActivityLevel)
	calories := h.calcSvc.RecommendCalories(tdee, user.Goal)

	c.JSON(http.StatusCreated, gin.H{
		"user":     user,
		"bmi":      bmi,
		"tdee":     tdee,
		"calories": calories,
	})
}

func (h *UserHandler) AddProgress(c *gin.Context) {
	var progress models.Progress
	if err := c.ShouldBindJSON(&progress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure the user exists
	_, err := h.repo.FindByID(progress.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := h.repo.Db.Create(&progress).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save progress"})
		return
	}

	c.JSON(http.StatusCreated, progress)
}

func (h *UserHandler) GetProgress(c *gin.Context) {
	userID := c.Param("user_id")
	var progressRecords []models.Progress
	if err := h.repo.Db.Where("user_id = ?", userID).Find(&progressRecords).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch progress"})
		return
	}

	c.JSON(http.StatusOK, progressRecords)
}

func (h *UserHandler) GetPlan(c *gin.Context) {
	userID := c.Param("user_id")

	// Stringni int ga aylantirish
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// int ni uint ga aylantirish
	user, err := h.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	bmr := h.calcSvc.CalculateBMR(user.Weight, user.Height, user.Age, user.Gender)
	tdee := h.calcSvc.CalculateTDEE(bmr, user.ActivityLevel)

	mealPlan := h.calcSvc.GenerateMealPlan(user, tdee)
	exercisePlan := h.calcSvc.GenerateExercisePlan(user)

	// Save plans to database
	for _, meal := range mealPlan {
		h.repo.Db.Create(&meal)
	}
	for _, ex := range exercisePlan {
		h.repo.Db.Create(&ex)
	}

	c.JSON(http.StatusOK, gin.H{
		"meal_plan":     mealPlan,
		"exercise_plan": exercisePlan,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var loginData struct {
		FullName string `json:"full_name"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.repo.Db.Where("full_name = ? AND password = ?", loginData.FullName, loginData.Password).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (h *UserHandler) CompleteAction(c *gin.Context) {
	userID := c.GetUint("user_id") // From AuthMiddleware
	var actionData struct {
		Action string `json:"action"` // "meal" or "exercise"
	}
	if err := c.ShouldBindJSON(&actionData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Award points
	point := h.calcSvc.AwardPoints(userID, actionData.Action)
	if err := h.repo.AddPoint(&point); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to award points"})
		return
	}

	// Check achievements (only for exercise)
	if actionData.Action == "exercise" {
		recentDays, err := h.repo.GetRecentExerciseDays(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check achievements"})
			return
		}
		achievements := h.calcSvc.CheckAchievements(userID, recentDays)
		for _, ach := range achievements {
			h.repo.AddAchievement(&ach)
		}
	}

	c.JSON(http.StatusOK, gin.H{"points": point})
}

func (h *UserHandler) GetGamificationData(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr) // Convert string to int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	points, err := h.repo.GetPoints(uint(userID)) // Cast to uint
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch points"})
		return
	}

	achievements, err := h.repo.GetAchievements(uint(userID)) // Cast to uint
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch achievements"})
		return
	}

	totalPoints := 0
	for _, p := range points {
		totalPoints += p.Points
	}

	c.JSON(http.StatusOK, gin.H{
		"total_points": totalPoints,
		"points":       points,
		"achievements": achievements,
	})
}

func (h *UserHandler) CreateReminder(c *gin.Context) {
	userID := c.GetUint("user_id")
	var reminder models.Reminder
	if err := c.ShouldBindJSON(&reminder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	reminder.UserID = userID
	reminder.Active = true

	if err := h.repo.CreateReminder(&reminder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reminder"})
		return
	}

	if err := h.reminderSvc.ScheduleReminder(&reminder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to schedule reminder"})
		return
	}

	c.JSON(http.StatusCreated, reminder)
}

func (h *UserHandler) GetReminders(c *gin.Context) {
	userID := c.GetUint("user_id")
	reminders, err := h.repo.GetReminders(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reminders"})
		return
	}
	c.JSON(http.StatusOK, reminders)
}

func (h *UserHandler) DisableReminder(c *gin.Context) {
	userID := c.GetUint("user_id")
	reminderID := c.Param("reminder_id")
	var reminder models.Reminder
	if err := h.repo.Db.Where("id = ? AND user_id = ?", reminderID, userID).First(&reminder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reminder not found"})
		return
	}

	reminder.Active = false
	if err := h.repo.UpdateReminder(&reminder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disable reminder"})
		return
	}

	h.reminderSvc.RemoveReminder(reminder.ID)
	c.JSON(http.StatusOK, gin.H{"message": "Reminder disabled"})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.GetUint("user_id")
	var updateData struct {
		Height             float64 `json:"height"`
		Weight             float64 `json:"weight"`
		WaistCircumference float64 `json:"waist_circumference"`
		ActivityLevel      string  `json:"activity_level,omitempty"`
		EatingHabits       string  `json:"eating_habits,omitempty"`
		Goal               string  `json:"goal,omitempty"`
	}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.repo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if a month has passed since last update
	if !user.LastUpdated.IsZero() && time.Since(user.LastUpdated).Hours() < 720 { // 720 hours = 30 days
		c.JSON(http.StatusForbidden, gin.H{"error": "Can only update once per month"})
		return
	}

	// Update user data
	user.Height = updateData.Height
	user.Weight = updateData.Weight
	user.WaistCircumference = updateData.WaistCircumference
	if updateData.ActivityLevel != "" {
		user.ActivityLevel = updateData.ActivityLevel
	}
	if updateData.EatingHabits != "" {
		user.EatingHabits = updateData.EatingHabits
	}
	if updateData.Goal != "" {
		user.Goal = updateData.Goal
	}
	user.LastUpdated = time.Now()

	if err := h.repo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Recalculate TDEE and generate new plans
	tdee, mealPlan, exercisePlan := h.calcSvc.UpdateUserPlan(user)

	// Clear old plans
	if err := h.repo.DeleteOldPlans(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear old plans"})
		return
	}

	// Save new plans
	for _, meal := range mealPlan {
		h.repo.Db.Create(&meal)
	}
	for _, ex := range exercisePlan {
		h.repo.Db.Create(&ex)
	}

	c.JSON(http.StatusOK, gin.H{
		"user":          user,
		"tdee":          tdee,
		"meal_plan":     mealPlan,
		"exercise_plan": exercisePlan,
	})
}
