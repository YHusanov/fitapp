package main

import (
	"diplomIshi/internal/config"
	"diplomIshi/internal/models"
	"diplomIshi/internal/routes"
	"log"
)

func main() {
	cfg := config.LoadConfig()

	err := cfg.DB.AutoMigrate(&models.User{}, &models.Progress{}, &models.Meal{}, &models.Exercise{}, &models.Point{}, &models.Achievement{}, &models.Reminder{}, &models.Post{}, &models.Comment{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	r := routes.SetupRoutes(cfg)
	r.Run(":" + cfg.Port)
}
