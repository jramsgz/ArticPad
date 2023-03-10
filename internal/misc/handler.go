package misc

import "github.com/gofiber/fiber/v2"

// Create a handler. Leave this empty, as we have no domains nor use-cases.
type MiscHandler struct{}

// Represents a new handler.
func NewMiscHandler(healthRoute fiber.Router) {
	handler := &MiscHandler{}

	// Declare routing.
	healthRoute.Get("", handler.defaultResponse)
}

// Defeault API response.
func (h *MiscHandler) defaultResponse(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"message": "Welcome to the API!",
	})
}
