package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/jramsgz/articpad/config"
)

func RegisterMiddlewares(app *fiber.App) {
	app.Use(ErrorHandling())
	app.Use(requestid.New(requestid.Config{
		// Request ID disabled for static files (not prefixed with /api)
		Next: func(c *fiber.Ctx) bool {
			return c.Path() == "/api"
		},
	}))
	app.Use(recover.New())
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
