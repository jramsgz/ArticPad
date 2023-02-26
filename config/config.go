package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config func to get env value
func Config(key string, defaultValue string) string {
	// load .env file
	err := godotenv.Load("./config/.env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
