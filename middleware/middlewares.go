package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jramsgz/articpad/config"
	"github.com/rs/zerolog"
)

// RegisterMiddlewares register middlewares for the app
func RegisterMiddlewares(app *fiber.App, logger zerolog.Logger) {
	app.Use(Logger(logger, func(c *fiber.Ctx) bool {
		return c.Path() == "/health" // skip logging for health check
	}))
	app.Use(cors.New(cors.Config{
		MaxAge:       1800,
		AllowOrigins: config.GetString("APP_URL", "http://localhost:3000"),
	}))
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))
}
