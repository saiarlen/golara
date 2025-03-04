package middlewares

import (
	"ekycapp/app/controllers/helpers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func LicenseValidator(c *fiber.Ctx) error {
	store := session.New()
	s, err := store.Get(c)

	if err != nil {
		panic(err)
	}

	status := s.Get("app_status")

	if status != "active" {

		return helpers.Error(c, "EG0001", "Invalid license contact administrator.", 401)
	}
	return c.Next()
}
