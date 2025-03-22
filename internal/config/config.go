package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DB   *gorm.DB
	Port string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Agar DATABASE_URL mavjud bo‘lsa, undan foydalanamiz
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Agar DATABASE_URL yo‘q bo‘lsa, eski usul bilan DSN yasaymiz
		dsn = "host=" + os.Getenv("DB_HOST") +
			" user=" + os.Getenv("DB_USER") +
			" password=" + os.Getenv("DB_PASSWORD") +
			" dbname=" + os.Getenv("DB_NAME") +
			" port=" + os.Getenv("DB_PORT") +
			" sslmode=require"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	return &Config{
		DB:   db,
		Port: os.Getenv("PORT"),
	}
}
