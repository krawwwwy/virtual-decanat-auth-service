package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	} //загружаем .env
}

func GetConnectionString() string {
	fmt.Println(os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", //достаем значения из .env
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
}
