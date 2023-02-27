package router

import (
	"github.com/jramsgz/articpad/handler"
	"github.com/jramsgz/articpad/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupApiRoutes setup router for api
func SetupApiRoutes(app *fiber.App) {
	// Middleware
	api := app.Group("/api")
	apiv1 := api.Group("/v1")
	apiv1.Get("/", handler.Hello)

	// Auth
	auth := apiv1.Group("/auth")
	auth.Post("/login", handler.Login)

	// User
	user := apiv1.Group("/user")
	user.Get("/:id", handler.GetUser)
	user.Post("/", handler.CreateUser)
	user.Patch("/:id", middleware.Protected(), handler.UpdateUser)
	user.Delete("/:id", middleware.Protected(), handler.DeleteUser)

	// Product
	product := apiv1.Group("/product")
	product.Get("/", handler.GetAllProducts)
	product.Get("/:id", handler.GetProduct)
	product.Post("/", middleware.Protected(), handler.CreateProduct)
	product.Delete("/:id", middleware.Protected(), handler.DeleteProduct)

	// 404 Handler
	api.Use(func(c *fiber.Ctx) error {
		// JSON api response
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"error":   "Not Found",
		})
	})
}
