package routes

import (
	"diplomIshi/internal/config"
	"diplomIshi/internal/handlers"
	"diplomIshi/internal/repository"
	"diplomIshi/internal/services"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	userRepo := repository.NewUserRepository(cfg.DB)
	calcSvc := services.NewCalculatorService(cfg.DB)
	reminderSvc := services.NewReminderService(userRepo)
	userHandler := handlers.NewUserHandler(userRepo, calcSvc, reminderSvc)
	communityHandler := handlers.NewCommunityHandler(userRepo)

	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	auth := r.Group("/").Use(handlers.AuthMiddleware())
	{
		auth.POST("/progress", userHandler.AddProgress)
		auth.GET("/progress/:user_id", userHandler.GetProgress)
		auth.GET("/plan/:user_id", userHandler.GetPlan)
		auth.POST("/complete", userHandler.CompleteAction)
		auth.GET("/gamification/:user_id", userHandler.GetGamificationData)
		auth.POST("/reminders", userHandler.CreateReminder)
		auth.GET("/reminders", userHandler.GetReminders)
		auth.PUT("/reminders/:reminder_id/disable", userHandler.DisableReminder)
		auth.POST("/community/posts", communityHandler.CreatePost)
		auth.GET("/community/posts", communityHandler.GetPosts)
		auth.POST("/community/posts/:post_id/comments", communityHandler.CreateComment)
		auth.GET("/community/posts/:post_id/comments", communityHandler.GetComments)
		auth.PUT("/update", userHandler.UpdateUser)
		auth.GET("/users", userHandler.GetAllUsers)
	}

	reminderSvc.Start()
	users := []uint{1} // Replace with dynamic user list if needed
	for _, userID := range users {
		reminderSvc.LoadReminders(userID)
	}

	return r
}
