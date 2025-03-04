package middlewares

import (
	"ekycapp/app/controllers/helpers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// ErrorHandler handles errors that occur in middleware.
//
// @param c - Context to pass to next middleware
//
// @return Error or nil to continue the chain without error or return to the next middleware in the chain with no
// @Usage: Simply "return err" in anywhere in the code it ill catch those errors
// Any 2001 errors use return err
func ErrorHandler(c *fiber.Ctx) error {
	err := c.Next()
	if err != nil {
		log.Errorf("Caught in middleware: %s", err)
		return helpers.Error(c, "EC2001", fiber.Map{"0": err.Error()}, 500)
	}
	return nil
}

// PanicHandler handles errors that occur go core level.
//
// @param c - Context to pass to next middleware
//
// @Usage: It automatically catches any panic errors
func PanicHandler(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Recovered from panic: %s", r)
			helpers.Error(c, "EC2001", fiber.Map{"0": r}, 500)
		}
	}()

	return c.Next()
}
