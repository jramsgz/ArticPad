package infrastructure

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/redis/v3"
	"github.com/jmoiron/sqlx"
	"github.com/jramsgz/articpad/config"
	"github.com/jramsgz/articpad/internal/migrations"
	"github.com/jramsgz/articpad/pkg/i18n"
	"github.com/jramsgz/articpad/pkg/mail"
	"github.com/rs/zerolog"
)

type App struct {
	fiber  *fiber.App
	logger zerolog.Logger
	db     *sqlx.DB
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
		Level: config.GetString(config.LogLevel),
		Dir:   config.GetString(config.LogDir),
	})

	i18n, err := startI18n(config.GetString(config.LocalesDir))
	if err != nil {
		logger.Fatal().Msgf("failed to start i18n service: %s", err.Error())
	}

	db, err := connectToDB(&DatabaseConfig{
		Driver:   config.GetString(config.DBDriver),
		Host:     config.GetString(config.DBHost),
		Username: config.GetString(config.DBUsername),
		Password: config.GetString(config.DBPassword),
		Port:     config.GetInt(config.DBPort),
		Database: config.GetString(config.DBDatabase),
	})
	if err != nil || db == nil {
		logger.Fatal().Msgf("Database connection error: %s", err)
	}

	if !fiber.IsChild() {
		logger.Info().Msg("Running migrations...")
		err = migrations.RunMigrations(db.DB, config.GetString(config.DBDriver), logger)
		if err != nil {
			logger.Fatal().Msgf("failed to automigrate models: %s", err.Error())
			return
		}
	}

	mailClient, err := mail.NewMailer(&mail.MailConfig{
		Host:     config.GetString(config.MailHost),
		Port:     config.GetInt(config.MailPort),
		Username: config.GetString(config.MailUser),
		Password: config.GetString(config.MailPass),
		From:     config.GetString(config.MailFrom),
		ForceTLS: config.GetString(config.MailForceTLS) == "true",
		SMTPAuth: config.GetString(config.MailSMTPAuth),
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
			Host:     config.GetString(config.RedisHost),
			Port:     config.GetInt(config.RedisPort),
			Username: config.GetString(config.RedisUsername),
			Password: config.GetString(config.RedisPassword),
			Database: config.GetInt(config.RedisDB),
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
		logger.Info().Msgf("Starting ArticPad %s with isProduction: %t", config.Version, config.GetString(config.Debug) == "false")
		logger.Info().Msgf("BuildTime: %s | Commit: %s", config.BuildTime, config.Commit)
		logger.Info().Msgf("Listening on %s", config.GetString(config.AppAddr))
	}
	if err := app.fiber.Listen(config.GetString(config.AppAddr)); err != nil {
		logger.Fatal().Err(err).Msg("Error starting server")
	}

	if !fiber.IsChild() {
		logger.Info().Msg("Shutting down server...")
		serverShutdown.Wait()
	}

	_ = db.Close()
	_ = logFile.Close()
	_ = mailClient.Close()
	_ = redisDB.Close()
}
