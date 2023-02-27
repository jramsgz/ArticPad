package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/utils"
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
		AllowOrigins: config.Config("APP_URL", "http://localhost:3000"),
	}))
	app.Use(csrf.New(csrf.Config{
		KeyLookup:      "header:X-Csrf-Token",
		CookieName:     "articpad_csrf_",
		CookieSameSite: "Strict",
		Expiration:     3 * time.Hour,
		KeyGenerator:   utils.UUID,
	}))
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))
}
