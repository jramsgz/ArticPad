package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var (
	// Version of the application
	Version = "dev"
	// BuildTime of the application
	BuildTime = time.Now().Format(time.RFC3339)
	// Commit of the application
	Commit = "dev build"
)

// LoadEnv func to load .env file, this should be called before using GetString or GetInt functions
func LoadEnv() error {
	return godotenv.Load("./config/.env")
}

// GetString func to get string value from environment variable
func GetString(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

// GetInt func to get int value from environment variable
func GetInt(key string, defaultValue int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intValue, ok := strconv.Atoi(value); ok == nil {
			return intValue
		}
	}
	return defaultValue
}
