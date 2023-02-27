package config

import (
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv func to load .env file, this should be called before using Config func
func LoadEnv() error {
	return godotenv.Load("./config/.env")
}

// Config func to get env value
func Config(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
