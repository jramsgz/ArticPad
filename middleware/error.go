package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandling() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Status code defaults to 500
		code := fiber.StatusInternalServerError
		// Message defaults to "Internal Server Error"
		message := "Internal Server Error"

		// Get error
		err := c.Next()
		if err != nil {
			// Check if it's a fiber.Error
			if e, ok := err.(*fiber.Error); ok {
				// Override status code if fiber.Error type
				code = e.Code
				message = e.Message
			}

			// Log error
			log.Printf("[ERROR] %d - %s", code, err.Error())

			// Send custom error page
			return c.Status(code).JSON(fiber.Map{
				"success":   false,
				"error":     message,
				"requestId": c.Locals("requestid"),
			})
		}
		return nil
	}
}
