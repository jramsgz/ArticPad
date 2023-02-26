package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jramsgz/articpad/config"
)

// SetupWebRoutes setup router for the web app
func SetupWebRoutes(app *fiber.App) {
	// Serve Single Page application on "/"
	// assume static file at static folder
	app.Static("/", config.Config("static_dir", "static"), fiber.Static{
		Compress: true,
		MaxAge:   3600,
	})

	// Panic test route, this brings up an error
	app.Get("/panic", func(ctx *fiber.Ctx) error {
		panic("Hi, I'm a panic error!")
	})

	app.Get("/*", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("./" + config.Config("static_dir", "static") + "/index.html")
	})
}
