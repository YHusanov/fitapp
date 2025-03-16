package models

import "time"

type Meal struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	UserID   uint   `json:"user_id"`
	Name     string `json:"name"` // e.g., "Pilaf", "Soup"
	Calories int    `json:"calories"`
	Protein  int    `json:"protein"`
	Carbs    int    `json:"carbs"`
	Fat      int    `json:"fat"`
	MealType string `json:"meal_type"` // e.g., "breakfast", "snack"
}

type Exercise struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	UserID         uint   `json:"user_id"`
	Name           string `json:"name"`     // e.g., "Push-ups", "Running"
	Duration       int    `json:"duration"` // in minutes
	CaloriesBurned int    `json:"calories_burned"`
	ImageURL       string `json:"image_url"`
}

type Point struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	Points    int       `json:"points"`
	Reason    string    `json:"reason"` // e.g., "Completed meal", "Completed exercise"
	CreatedAt time.Time `json:"created_at"`
}

type Achievement struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	Name      string    `json:"name"` // e.g., "5-Day Streak"
	CreatedAt time.Time `json:"created_at"`
}

type Reminder struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	Type      string    `json:"type"`    // "meal" or "exercise"
	Time      string    `json:"time"`    // Cron format, e.g., "0 8 * * *" (8:00 AM daily)
	Message   string    `json:"message"` // e.g., "Time for breakfast!"
	Active    bool      `json:"active"`  // Enable/disable reminder
	CreatedAt time.Time `json:"created_at"`
}

type Post struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	PostID    uint      `json:"post_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
