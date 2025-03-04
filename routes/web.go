package routes

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func WebRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).Format("Northing to see here | Contact Saiarlen | https://www.saiarlen.in")
	})
}
