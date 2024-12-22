package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var logLevel string

func setLogLevel() {
	logLevel = GetEnv("LOG_LEVEL", "INFO")
}

func GetLogLevel() string {
	return logLevel
}

func LoadEnv(provided_path string) {
	err := godotenv.Load(provided_path)
	if err != nil {
		log.Println("No .env file found, using default environment variables")
	}
	setLogLevel()
}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
