package middlewares

import (
	"ekycapp/app/controllers/helpers"
	"ekycapp/config"

	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware is a middleware function to verify the token
func AdminTokenAuth(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	//return c.Next()
	if token == "" {

		return helpers.Error(c, "EG0000", "Unauthorized Admin Token, Empty", 401)
	}
	secretKey := []byte(config.Denv("ADMIN_TOKEN_SECRET"))
	claims, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil || claims == nil {

		return helpers.Error(c, "EG0000", "Unauthorized Admin Token, Error", 401)
	}
	_, ok := claims.Claims.(jwt.MapClaims)
	if !ok {
		return helpers.Error(c, "EG0000", "Unauthorized Admin Token", 401)
	}

	return c.Next()
}
