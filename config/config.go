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

// EnvVar type for environment variables
type EnvVar string

// Available environment variables
const (
	DBDriver       EnvVar = "DB_DRIVER"
	DBHost         EnvVar = "DB_HOST"
	DBUsername     EnvVar = "DB_USERNAME"
	DBPassword     EnvVar = "DB_PASSWORD"
	DBPort         EnvVar = "DB_PORT"
	DBDatabase     EnvVar = "DB_DATABASE"
	RedisHost      EnvVar = "REDIS_HOST"
	RedisPort      EnvVar = "REDIS_PORT"
	RedisUsername  EnvVar = "REDIS_USERNAME"
	RedisPassword  EnvVar = "REDIS_PASSWORD"
	RedisDB        EnvVar = "REDIS_DB"
	MailHost       EnvVar = "MAIL_HOST"
	MailPort       EnvVar = "MAIL_PORT"
	MailUser       EnvVar = "MAIL_USER"
	MailPass       EnvVar = "MAIL_PASS"
	MailFrom       EnvVar = "MAIL_FROM"
	MailForceTLS   EnvVar = "MAIL_FORCE_TLS"
	MailSMTPAuth   EnvVar = "MAIL_SMTP_AUTH"
	EnableMail     EnvVar = "ENABLE_MAIL"
	Debug          EnvVar = "DEBUG"
	LogLevel       EnvVar = "LOG_LEVEL"
	LogDir         EnvVar = "LOG_DIR"
	AppAddr        EnvVar = "APP_ADDR"
	StaticDir      EnvVar = "STATIC_DIR"
	AppURL         EnvVar = "APP_URL"
	Secret         EnvVar = "SECRET"
	TrustedProxies EnvVar = "TRUSTED_PROXIES"
	TemplatesDir   EnvVar = "TEMPLATES_DIR"
	LocalesDir     EnvVar = "LOCALES_DIR"
	RateLimit      EnvVar = "RATE_LIMIT_AUTH"
)

// Map of default values
var defaults = map[EnvVar]string{
	DBDriver:       "sqlite",
	DBHost:         "localhost",
	DBUsername:     "root",
	DBPassword:     "",
	DBPort:         "5432",
	DBDatabase:     "config/articpad.db",
	RedisHost:      "localhost",
	RedisPort:      "6379",
	RedisUsername:  "",
	RedisPassword:  "",
	RedisDB:        "0",
	MailHost:       "localhost",
	MailPort:       "25",
	MailUser:       "",
	MailPass:       "",
	MailFrom:       "ArticPad",
	MailForceTLS:   "false",
	MailSMTPAuth:   "login",
	EnableMail:     "false",
	Debug:          "false",
	LogLevel:       "debug",
	LogDir:         "./logs",
	AppAddr:        ":8080",
	StaticDir:      "static",
	AppURL:         "http://localhost:8080",
	Secret:         "MyRandomSecureSecret",
	TrustedProxies: "",
	TemplatesDir:   "templates",
	LocalesDir:     "locales",
	RateLimit:      "true",
}

// LoadEnv loads the .env file, this should be called before using GetString or GetInt functions
func LoadEnv() error {
	loadDefaults()
	return godotenv.Overload("./config/.env")
}

// GetString func to get string value from environment variable
func GetString(key EnvVar, defaultValue ...string) string {
	if value, ok := os.LookupEnv(string(key)); ok {
		if value != "" {
			return value
		}
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

// GetInt func to get int value from environment variable
func GetInt(key EnvVar, defaultValue ...int) int {
	if value, ok := os.LookupEnv(string(key)); ok {
		if value != "" {
			if intValue, ok := strconv.Atoi(value); ok == nil {
				return intValue
			}
		}
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

func loadDefaults() {
	for key, value := range defaults {
		if _, ok := os.LookupEnv(string(key)); !ok {
			os.Setenv(string(key), value)
		}
	}
}
