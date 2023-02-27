package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jramsgz/articpad/config"
)

// SetupWebRoutes setup router for the web app
func SetupWebRoutes(app *fiber.App) {
	// Serve Single Page application on "/"
	// assume static file at static folder
	app.Static("/", config.Config("STATIC_DIR", "static"), fiber.Static{
		Compress: true,
		MaxAge:   3600,
	})

	// Health check route
	app.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "OK",
		})
	})

	// Panic test route, this brings up an error
	app.Get("/panic", func(ctx *fiber.Ctx) error {
		panic("Hi, I'm a panic error!")
	})

	app.Get("/*", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("./" + config.Config("STATIC_DIR", "static") + "/index.html")
	})
}
