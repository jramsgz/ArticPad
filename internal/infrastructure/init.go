package infrastructure

import (
	"log"

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
	logger, _, logFile := StartLogger(&LoggerConfig{
		Level: config.GetString("LOG_LEVEL", "debug"),
		Dir:   config.GetString("LOG_DIR", "./logs"),
	})

	// Try to connect to the specified database.
	db, err := ConnectToDB(&DatabaseConfig{
		Driver:   config.GetString("DB_DRIVER", "sqlite"),
		Host:     config.GetString("DB_HOST", "localhost"),
		Username: config.GetString("DB_USERNAME", "root"),
		Password: config.GetString("DB_PASSWORD", ""),
		Port:     config.GetInt("DB_PORT", 3306),
		Database: config.GetString("DB_DATABASE", "config/articpad.db"),
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
		Host:     config.GetString("MAIL_HOST", "localhost"),
		Port:     config.GetInt("MAIL_PORT", 25),
		Username: config.GetString("MAIL_USERNAME", ""),
		Password: config.GetString("MAIL_PASSWORD", ""),
		From:     config.GetString("MAIL_FROM", "ArticPad <"+config.GetString("MAIL_USERNAME", "")+">"),
		ForceTLS: config.GetString("MAIL_FORCE_TLS", "false") == "true",
	})
	if err != nil || mailClient == nil {
		logger.Fatal().Msgf("Mail server connection error: %s", err)
	}

	// Start Fiber.
	StartFiberServer(logger, db, mailClient)

	// TODO: Only for main process or every process?
	if true {
		// Close resources.
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
		_ = logFile.Close()
		_ = mailClient.Close()
	}
}
