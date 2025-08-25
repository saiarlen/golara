package routes

import (
	"github.com/test/myapp/app/controllers"
	"github.com/test/myapp/framework"

	"github.com/gofiber/fiber/v2"
)

// RegisterAPIRoutes registers all API routes
func RegisterAPIRoutes(app *framework.Golara) {
	// Create controllers
	userController := &controllers.UserController{
		DB:     app.DB,
		Events: app.Events,
	}

	// API group
	api := app.Group("/api")

	// User routes
	api.Get("/users", userController.Index)
	api.Post("/users", userController.Store)

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "healthy",
			"framework": "Golara",
			"version":   "1.0.0",
		})
	})
}
