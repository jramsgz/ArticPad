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

// Map of default values
var defaults = map[string]string{
	"DB_DRIVER":       "sqlite",
	"DB_HOST":         "localhost",
	"DB_USERNAME":     "root",
	"DB_PASSWORD":     "",
	"DB_PORT":         "5432",
	"DB_DATABASE":     "config/articpad.db",
	"MAIL_HOST":       "localhost",
	"MAIL_PORT":       "25",
	"MAIL_USER":       "",
	"MAIL_PASS":       "",
	"MAIL_FROM":       "ArticPad",
	"MAIL_FORCE_TLS":  "false",
	"ENABLE_MAIL":     "false",
	"DEBUG":           "false",
	"LOG_LEVEL":       "debug",
	"LOG_DIR":         "./logs",
	"APP_ADDR":        ":8080",
	"STATIC_DIR":      "static",
	"APP_URL":         "http://localhost:8080",
	"SECRET":          "MyRandomSecureSecret",
	"TRUSTED_PROXIES": "",
	"TEMPLATES_DIR":   "templates",
}

// LoadEnv loads the .env file, this should be called before using GetString or GetInt functions
func LoadEnv() error {
	loadDefaults()
	return godotenv.Overload("./config/.env")
}

// GetString func to get string value from environment variable
func GetString(key string, defaultValue ...string) string {
	if value, ok := os.LookupEnv(key); ok {
		if value != "" {
			return value
		}
	}
	return defaultValue[0]
}

// GetInt func to get int value from environment variable
func GetInt(key string, defaultValue ...int) int {
	if value, ok := os.LookupEnv(key); ok {
		if value != "" {
			if intValue, ok := strconv.Atoi(value); ok == nil {
				return intValue
			}
		}
	}
	return defaultValue[0]
}

func loadDefaults() {
	for key, value := range defaults {
		if _, ok := os.LookupEnv(key); !ok {
			os.Setenv(key, value)
		}
	}
}
