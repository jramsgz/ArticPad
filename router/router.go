package router

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	SetupApiRoutes(app)
	SetupWebRoutes(app)
}
