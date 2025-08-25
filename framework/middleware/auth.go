package middleware

import (
	"github.com/test/myapp/framework/errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWTConfig holds JWT middleware configuration
type JWTConfig struct {
	SecretKey    string
	TokenLookup  string // "header:Authorization,query:token,cookie:jwt"
	AuthScheme   string // "Bearer"
	ContextKey   string // "user"
	ErrorHandler fiber.ErrorHandler
}

// DefaultJWTConfig returns default JWT configuration
func DefaultJWTConfig() JWTConfig {
	return JWTConfig{
		TokenLookup: "header:Authorization",
		AuthScheme:  "Bearer",
		ContextKey:  "user",
	}
}

// JWT returns JWT middleware
func JWT(config ...JWTConfig) fiber.Handler {
	cfg := DefaultJWTConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	return func(c *fiber.Ctx) error {
		token := extractToken(c, cfg)
		if token == "" {
			return errors.AuthErr("Missing or invalid token")
		}

		// Parse and validate token
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.SecretKey), nil
		})

		if err != nil || !parsedToken.Valid {
			return errors.AuthErr("Invalid token")
		}

		// Store user claims in context
		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
			c.Locals(cfg.ContextKey, claims)
		}

		return c.Next()
	}
}

// extractToken extracts JWT token from request
func extractToken(c *fiber.Ctx, cfg JWTConfig) string {
	parts := strings.Split(cfg.TokenLookup, ",")

	for _, part := range parts {
		authParts := strings.Split(strings.TrimSpace(part), ":")
		if len(authParts) != 2 {
			continue
		}

		switch authParts[0] {
		case "header":
			auth := c.Get(authParts[1])
			if auth != "" {
				if cfg.AuthScheme != "" {
					prefix := cfg.AuthScheme + " "
					if strings.HasPrefix(auth, prefix) {
						return strings.TrimPrefix(auth, prefix)
					}
				} else {
					return auth
				}
			}
		case "query":
			return c.Query(authParts[1])
		case "cookie":
			return c.Cookies(authParts[1])
		}
	}

	return ""
}

// RequireRole middleware for role-based access control
func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user")
		if user == nil {
			return errors.AuthErr("Authentication required")
		}

		claims, ok := user.(jwt.MapClaims)
		if !ok {
			return errors.AuthErr("Invalid user claims")
		}

		userRole, exists := claims["role"].(string)
		if !exists {
			return errors.AuthErr("User role not found")
		}

		// Check if user has required role
		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}

		return errors.NewAppError(errors.AuthError, "Insufficient permissions", 403)
	}
}
