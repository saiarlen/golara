package routes

import (
	"os"
	"strings"
	"github.com/test/myapp/framework"

	"github.com/gofiber/fiber/v2"
)

// RegisterWebRoutes registers all web routes
func RegisterWebRoutes(app *framework.Golara) {

	// Home page
	app.App.Get("/", func(c *fiber.Ctx) error {
		// Read template and replace placeholders
		content, err := os.ReadFile("resources/views/welcome.html")
		if err != nil {
			return c.Status(500).SendString("Template not found")
		}
		
		// Replace placeholders with actual values
		html := string(content)
		html = strings.ReplaceAll(html, "{{.AppName}}", "Golara Framework")
		html = strings.ReplaceAll(html, "{{.version}}", "1.0.0")
		html = strings.ReplaceAll(html, "{{.title}}", "Welcome to Golara")
		
		c.Set("Content-Type", "text/html")
		return c.SendString(html)
	})

	// About page
	app.App.Get("/about", func(c *fiber.Ctx) error {
		// Read template and replace placeholders
		content, err := os.ReadFile("resources/views/about.html")
		if err != nil {
			return c.Status(500).SendString("Template not found")
		}
		
		// Replace placeholders with actual values
		html := string(content)
		html = strings.ReplaceAll(html, "{{.AppName}}", "Golara Framework")
		html = strings.ReplaceAll(html, "{{.title}}", "About Golara")
		
		c.Set("Content-Type", "text/html")
		return c.SendString(html)
	})

	// Documentation redirect
	app.App.Get("/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/docs/html/index.html")
	})
}