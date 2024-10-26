package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv(provided_path string) {
	err := godotenv.Load(provided_path)
	if err != nil {
		log.Println("No .env file found, using default environment variables")
	}
}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
