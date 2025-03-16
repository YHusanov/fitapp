package services

import (
	"diplomIshi/internal/models"
	"gorm.io/gorm"
	"math"
	"time"
)

type CalculatorService struct {
	db *gorm.DB // Bazaga ulanish
}

func NewCalculatorService(db *gorm.DB) *CalculatorService {
	return &CalculatorService{db: db}
}

func (s *CalculatorService) CalculateBMI(weight, height float64) float64 {
	heightInMeters := height / 100
	return weight / math.Pow(heightInMeters, 2)
}

func (s *CalculatorService) CalculateBMR(weight, height float64, age int, gender string) float64 {
	if gender == "male" {
		return 10*weight + 6.25*height - 5*float64(age) + 5
	}
	return 10*weight + 6.25*height - 5*float64(age) - 161 // female
}

func (s *CalculatorService) CalculateTDEE(bmr float64, activityLevel string) float64 {
	activityCoefficients := map[string]float64{
		"low":     1.2,
		"average": 1.375,
		"high":    1.55,
		"athlete": 1.9,
	}
	coeff := activityCoefficients[activityLevel]
	return bmr * coeff
}

func (s *CalculatorService) RecommendCalories(tdee float64, goal string) int {
	switch goal {
	case "weight_loss":
		return int(tdee * 0.85) // 15% deficit
	case "weight_gain":
		return int(tdee * 1.15) // 15% surplus
	default:
		return int(tdee) // maintenance
	}
}

func (s *CalculatorService) GenerateMealPlan(user *models.User, tdee float64) []models.Meal {
	calories := s.RecommendCalories(tdee, user.Goal)
	mealsPerDay := 6
	caloriesPerMeal := calories / mealsPerDay

	var allMeals []models.Meal
	// Vazn tashlash uchun kaloriyasi past ovqatlarni tanlash
	if user.Goal == "weight_loss" {
		s.db.Where("calories <= ?", 250).Find(&allMeals) // 250 kcal dan past ovqatlar
	} else {
		s.db.Find(&allMeals)
	}

	if len(allMeals) == 0 {
		return []models.Meal{}
	}

	var userMeals []models.Meal
	for i := 0; i < mealsPerDay && i < len(allMeals); i++ {
		meal := allMeals[i%len(allMeals)]
		meal.Calories = caloriesPerMeal // Kaloriyani moslashtirish
		userMeals = append(userMeals, meal)
	}
	return userMeals
}

func (s *CalculatorService) GenerateExercisePlan(user *models.User) []models.Exercise {
	var allExercises []models.Exercise
	s.db.Find(&allExercises)

	if len(allExercises) == 0 {
		return []models.Exercise{}
	}

	var userExercises []models.Exercise
	for _, ex := range allExercises {
		if user.Goal == "weight_loss" && ex.CaloriesBurned > 100 { // Kardio uchun yuqori kaloriya sarfi
			userExercises = append(userExercises, ex)
		} else if len(userExercises) < 3 {
			userExercises = append(userExercises, ex)
		}
	}
	return userExercises
}

func (s *CalculatorService) AwardPoints(userID uint, action string) models.Point {
	points := 0
	reason := ""
	switch action {
	case "meal":
		points = 10
		reason = "Completed meal"
	case "exercise":
		points = 20
		reason = "Completed exercise"
	}
	return models.Point{
		UserID: userID,
		Points: points,
		Reason: reason,
	}
}

func (s *CalculatorService) CheckAchievements(userID uint, recentDays []time.Time) []models.Achievement {
	var achievements []models.Achievement
	now := time.Now()

	// Check for 5-day workout streak
	if len(recentDays) >= 5 {
		streak := true
		for i := 0; i < 4; i++ {
			if recentDays[i].Sub(recentDays[i+1]).Hours() > 48 { // More than 2 days gap
				streak = false
				break
			}
		}
		if streak && recentDays[0].Truncate(24*time.Hour).Equal(now.Truncate(24*time.Hour)) {
			achievements = append(achievements, models.Achievement{
				UserID: userID,
				Name:   "5-Day Workout Streak",
			})
		}
	}

	return achievements
}

func (s *CalculatorService) UpdateUserPlan(user *models.User) (float64, []models.Meal, []models.Exercise) {
	bmr := s.CalculateBMR(user.Weight, user.Height, user.Age, user.Gender)
	tdee := s.CalculateTDEE(bmr, user.ActivityLevel)
	mealPlan := s.GenerateMealPlan(user, tdee)
	exercisePlan := s.GenerateExercisePlan(user)
	return tdee, mealPlan, exercisePlan
}
