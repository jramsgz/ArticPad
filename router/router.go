package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jramsgz/articpad/database"
	"github.com/rs/zerolog"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App, logger zerolog.Logger, db *database.Database) {
	SetupApiRoutes(app, logger, db)
	SetupWebRoutes(app)
}
