package infrastructure

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jramsgz/articpad/config"
	"github.com/jramsgz/articpad/internal/user"
	"github.com/jramsgz/articpad/pkg/mail"
)

// Run ArticPad API & Static Server
func Run() {
	// Load configuration from .env file.
	if err := config.LoadEnv(); err != nil {
		log.Fatal(err)
	}

	// Start logger.
	logger, _, logFile := startLogger(&LoggerConfig{
		Level: config.GetString("LOG_LEVEL"),
		Dir:   config.GetString("LOG_DIR"),
	})

	// Start i18n service.
	i18n, err := startI18n(config.GetString("LOCALES_DIR"))
	if err != nil {
		logger.Fatal().Msgf("failed to start i18n service: %s", err.Error())
	}

	// Try to connect to the specified database.
	db, err := connectToDB(&DatabaseConfig{
		Driver:   config.GetString("DB_DRIVER"),
		Host:     config.GetString("DB_HOST"),
		Username: config.GetString("DB_USERNAME"),
		Password: config.GetString("DB_PASSWORD"),
		Port:     config.GetInt("DB_PORT"),
		Database: config.GetString("DB_DATABASE"),
	})
	if err != nil || db == nil {
		logger.Fatal().Msgf("Database connection error: %s", err)
	}

	if !fiber.IsChild() {
		logger.Info().Msg("Running migrations...")
		// Auto-migrate database models
		err = db.AutoMigrate(&user.User{})
		if err != nil {
			logger.Fatal().Msgf("failed to automigrate models: %s", err.Error())
			return
		}
	}

	// Connect to mail server.
	mailClient, err := mail.NewMailer(&mail.MailConfig{
		Host:     config.GetString("MAIL_HOST"),
		Port:     config.GetInt("MAIL_PORT"),
		Username: config.GetString("MAIL_USERNAME"),
		Password: config.GetString("MAIL_PASSWORD"),
		From:     config.GetString("MAIL_FROM", "ArticPad <"+config.GetString("MAIL_USERNAME")+">"),
		ForceTLS: config.GetString("MAIL_FORCE_TLS") == "true",
	})
	if err != nil || mailClient == nil {
		logger.Fatal().Msgf("Mail server connection error: %s", err)
	}

	// Start Fiber.
	app := startFiberServer(logger, db, mailClient, i18n)

	// Setup graceful shutdown.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	var serverShutdown sync.WaitGroup

	go func() {
		<-c
		serverShutdown.Add(1)
		defer serverShutdown.Done()
		_ = app.ShutdownWithTimeout(60 * time.Second)
	}()

	if !fiber.IsChild() {
		logger.Info().Msgf("Starting ArticPad %s with isProduction: %t", config.Version, true)
		logger.Info().Msgf("BuildTime: %s | Commit: %s", config.BuildTime, config.Commit)
		logger.Info().Msgf("Listening on %s", config.GetString("APP_ADDR"))
	}
	if err := app.Listen(config.GetString("APP_ADDR")); err != nil {
		logger.Fatal().Err(err).Msg("Error starting server")
	}

	if !fiber.IsChild() {
		logger.Info().Msg("Shutting down server...")
		serverShutdown.Wait()
	}

	// TODO: Only for main process or every process?
	if true {
		// Close resources.
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
		_ = logFile.Close()
		_ = mailClient.Close()
	}
}
