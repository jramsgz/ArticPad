package main

import (
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jramsgz/articpad/config"
	"github.com/jramsgz/articpad/database"
	"github.com/jramsgz/articpad/middleware"
	"github.com/jramsgz/articpad/model"
	"github.com/jramsgz/articpad/router"
	"github.com/rs/zerolog"

	"github.com/gofiber/fiber/v2"
)

var (
	// Version of the application
	Version = "dev"
	// BuildTime of the application
	BuildTime = time.Now().Format(time.RFC3339)
	// Commit of the application
	Commit = "dev build"

	// Build information. Populated at build-time.
	// go build -ldflags "-X main.Version=1.0.0 -X main.BuildTime=2020-01-01T00:00:00Z -X main.Commit=abcdef"
)

// App contains the "global" components that are
// passed around.
type App struct {
	fiber  *fiber.App
	logger zerolog.Logger
	DB     *database.Database
}

func main() {
	// MultiWriter to log to both console and file
	if _, err := os.Stat("./logs"); os.IsNotExist(err) {
		os.MkdirAll("./logs", 0755)
	}
	file, err := os.OpenFile("./logs/articpad.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()
	mw := io.MultiWriter(os.Stdout, file)

	// Load .env file
	if err := config.LoadEnv(); err != nil {
		log.Fatal(err)
	}

	var isProduction bool = config.GetString("DEBUG", "false") == "false"
	// Set log level
	var logLevel zerolog.Level = zerolog.DebugLevel
	desiredLevel, err := zerolog.ParseLevel(config.GetString("LOG_LEVEL", "debug"))
	if err == nil {
		logLevel = desiredLevel
	}

	// Set trusted proxies
	var trustedProxies []string
	if config.GetString("TRUSTED_PROXIES", "") != "" {
		trustedProxies = strings.Split(config.GetString("TRUSTED_PROXIES", ""), ",")
	}
	var enableProxy bool = len(trustedProxies) > 0

	// App initialization
	app := App{
		fiber: fiber.New(fiber.Config{
			Prefork:                 isProduction,
			ServerHeader:            "ArticPad Server " + Version,
			AppName:                 "ArticPad",
			DisableStartupMessage:   isProduction,
			EnableTrustedProxyCheck: enableProxy,
			TrustedProxies:          trustedProxies,
		}),
		logger: zerolog.New(mw).With().Timestamp().Logger().Level(logLevel),
	}

	// Database initialization
	db, err := database.New(&database.DatabaseConfig{
		Driver:   config.GetString("DB_DRIVER", "sqlite"),
		Host:     config.GetString("DB_HOST", "localhost"),
		Username: config.GetString("DB_USERNAME", "root"),
		Password: config.GetString("DB_PASSWORD", ""),
		Port:     config.GetInt("DB_PORT", 3306),
		Database: config.GetString("DB_DATABASE", "config/articpad.db"),
	})

	// Auto-migrate database models
	if !fiber.IsChild() {
		if err != nil {
			log.Fatal("failed to connect to database:", err.Error())
		} else {
			if db == nil {
				log.Fatal("failed to connect to database: db is nil")
			} else {
				app.DB = db
				err := app.DB.AutoMigrate(&model.User{})
				if err != nil {
					log.Fatal("failed to automigrate models:", err.Error())
					return
				}
			}
		}
	}

	// Middleware registration
	middleware.RegisterMiddlewares(app.fiber, app.logger)

	// Route registration
	router.SetupRoutes(app.fiber, app.logger, app.DB)

	if !fiber.IsChild() {
		log.Printf("Starting ArticPad %s with isProduction: %t", Version, isProduction)
		log.Printf("BuildTime: %s | Commit: %s", BuildTime, Commit)
		log.Printf("Listening on %s", config.GetString("APP_ADDR", ":3000"))
	}
	if err := app.fiber.Listen(config.GetString("APP_ADDR", ":3000")); err != nil {
		app.logger.Fatal().Err(err).Msg("Error starting server")
	}
}
