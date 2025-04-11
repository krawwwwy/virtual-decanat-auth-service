package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func init() {
	// Загрузка переменных окружения из файла .env
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
}

func GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
}
