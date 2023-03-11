package health

import "github.com/gofiber/fiber/v2"

// Create a handler. Leave this empty, as we have no domains nor use-cases.
type HealthHandler struct{}

// Represents a new handler.
func NewHealthHandler(healthRoute fiber.Router) {
	handler := &HealthHandler{}

	// Declare routing.
	healthRoute.Get("", handler.healthCheck)
}

// Check for the health of the API.
func (h *HealthHandler) healthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"message": "OK",
	})
}