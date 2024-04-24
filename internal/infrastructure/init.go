package infrastructure

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/redis/v3"
	"github.com/jramsgz/articpad/config"
	"github.com/jramsgz/articpad/internal/user"
	"github.com/jramsgz/articpad/pkg/i18n"
	"github.com/jramsgz/articpad/pkg/mail"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type App struct {
	fiber  *fiber.App
	logger zerolog.Logger
	db     *gorm.DB
	mail   *mail.Mailer
	i18n   *i18n.I18n
	redis  *redis.Storage
}

// Run ArticPad API & Static Server
func Run() {
	if err := config.LoadEnv(); err != nil {
		log.Fatal(err)
	}

	logger, _, logFile := startLogger(&LoggerConfig{
		Level: config.GetString("LOG_LEVEL"),
		Dir:   config.GetString("LOG_DIR"),
	})

	i18n, err := startI18n(config.GetString("LOCALES_DIR"))
	if err != nil {
		logger.Fatal().Msgf("failed to start i18n service: %s", err.Error())
	}

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
		err = db.AutoMigrate(&user.User{})
		if err != nil {
			logger.Fatal().Msgf("failed to automigrate models: %s", err.Error())
			return
		}
	}

	mailClient, err := mail.NewMailer(&mail.MailConfig{
		Host:     config.GetString("MAIL_HOST"),
		Port:     config.GetInt("MAIL_PORT"),
		Username: config.GetString("MAIL_USERNAME"),
		Password: config.GetString("MAIL_PASSWORD"),
		From:     config.GetString("MAIL_FROM"),
		ForceTLS: config.GetString("MAIL_FORCE_TLS") == "true",
	})
	if err != nil || mailClient == nil {
		logger.Fatal().Msgf("Mail server connection error: %s", err)
	}

	var redisDB *redis.Storage
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Error().Msgf("Redis connection error: %s. Some features may not be available.", err)
			}
		}()
		redisDB = redis.New(redis.Config{
			Host:     config.GetString("REDIS_HOST"),
			Port:     config.GetInt("REDIS_PORT"),
			Username: config.GetString("REDIS_USERNAME"),
			Password: config.GetString("REDIS_PASSWORD"),
			Database: config.GetInt("REDIS_DB"),
		})
	}()

	app := &App{
		logger: logger,
		db:     db,
		mail:   mailClient,
		i18n:   i18n,
		redis:  redisDB,
	}
	app.fiber = app.startFiberServer()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	var serverShutdown sync.WaitGroup

	go func() {
		<-c
		serverShutdown.Add(1)
		defer serverShutdown.Done()
		_ = app.fiber.ShutdownWithTimeout(60 * time.Second)
	}()

	if !fiber.IsChild() {
		logger.Info().Msgf("Starting ArticPad %s with isProduction: %t", config.Version, config.GetString("DEBUG") == "false")
		logger.Info().Msgf("BuildTime: %s | Commit: %s", config.BuildTime, config.Commit)
		logger.Info().Msgf("Listening on %s", config.GetString("APP_ADDR"))
	}
	if err := app.fiber.Listen(config.GetString("APP_ADDR")); err != nil {
		logger.Fatal().Err(err).Msg("Error starting server")
	}

	if !fiber.IsChild() {
		logger.Info().Msg("Shutting down server...")
		serverShutdown.Wait()
	}

	sqlDB, _ := db.DB()
	_ = sqlDB.Close()
	_ = logFile.Close()
	_ = mailClient.Close()
	_ = redisDB.Close()
}
