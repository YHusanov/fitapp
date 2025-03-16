package models

import (
	"time"
)

type User struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	FullName           string    `json:"full_name"`
	Height             float64   `json:"height"`
	Weight             float64   `json:"weight"`
	Age                int       `json:"age"`
	Gender             string    `json:"gender"`
	ActivityLevel      string    `json:"activity_level"`
	EatingHabits       string    `json:"eating_habits"`
	Goal               string    `json:"goal"`
	WaistCircumference float64   `json:"waist_circumference"`
	Password           string    `json:"password" gorm:"not null"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	LastUpdated        time.Time `json:"last_updated"` // New field to track last update
}

type Progress struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	Weight    float64   `json:"weight"`
	Date      time.Time `json:"date"`
	Calories  int       `json:"calories"`
	PhotoURL  string    `json:"photo_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
