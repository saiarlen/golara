package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

// MiddlewareRegistry manages middleware registration and execution
type MiddlewareRegistry struct {
	global []fiber.Handler
	groups map[string][]fiber.Handler
}

// NewMiddlewareRegistry creates a new middleware registry
func NewMiddlewareRegistry() *MiddlewareRegistry {
	return &MiddlewareRegistry{
		global: make([]fiber.Handler, 0),
		groups: make(map[string][]fiber.Handler),
	}
}

// RegisterGlobal registers a global middleware
func (mr *MiddlewareRegistry) RegisterGlobal(handler fiber.Handler) {
	mr.global = append(mr.global, handler)
}

// RegisterGroup registers middleware for a specific group
func (mr *MiddlewareRegistry) RegisterGroup(group string, handler fiber.Handler) {
	if mr.groups[group] == nil {
		mr.groups[group] = make([]fiber.Handler, 0)
	}
	mr.groups[group] = append(mr.groups[group], handler)
}

// ApplyGlobal applies all global middleware to the app
func (mr *MiddlewareRegistry) ApplyGlobal(app *fiber.App) {
	for _, handler := range mr.global {
		app.Use(handler)
	}
}

// GetGroupMiddleware returns middleware for a specific group
func (mr *MiddlewareRegistry) GetGroupMiddleware(group string) []fiber.Handler {
	return mr.groups[group]
}

// Built-in middleware functions

// RequestID adds a unique request ID to each request
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = fmt.Sprintf("%d", time.Now().UnixNano())
			c.Set("X-Request-ID", requestID)
		}
		c.Locals("requestID", requestID)
		return c.Next()
	}
}

// RateLimiter implements basic rate limiting
func RateLimiter(maxRequests int, window time.Duration) fiber.Handler {
	requests := make(map[string][]time.Time)
	
	return func(c *fiber.Ctx) error {
		ip := c.IP()
		now := time.Now()
		
		// Clean old requests
		if reqs, exists := requests[ip]; exists {
			filtered := make([]time.Time, 0)
			for _, req := range reqs {
				if now.Sub(req) < window {
					filtered = append(filtered, req)
				}
			}
			requests[ip] = filtered
		}
		
		// Check rate limit
		if len(requests[ip]) >= maxRequests {
			return c.Status(429).JSON(fiber.Map{
				"error": "Rate limit exceeded",
				"retry_after": window.Seconds(),
			})
		}
		
		// Add current request
		requests[ip] = append(requests[ip], now)
		
		return c.Next()
	}
}

// APIKeyAuth validates API key from header
func APIKeyAuth(validKeys []string) fiber.Handler {
	keyMap := make(map[string]bool)
	for _, key := range validKeys {
		keyMap[key] = true
	}
	
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "API key required",
			})
		}
		
		if !keyMap[apiKey] {
			return c.Status(401).JSON(fiber.Map{
				"error": "Invalid API key",
			})
		}
		
		return c.Next()
	}
}

// CORS with configurable options
func CORS(origins []string, methods []string, headers []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		origin := c.Get("Origin")
		
		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range origins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}
		
		if allowed {
			c.Set("Access-Control-Allow-Origin", origin)
		}
		
		if len(methods) > 0 {
			c.Set("Access-Control-Allow-Methods", joinStrings(methods, ", "))
		}
		
		if len(headers) > 0 {
			c.Set("Access-Control-Allow-Headers", joinStrings(headers, ", "))
		}
		
		c.Set("Access-Control-Allow-Credentials", "true")
		
		if c.Method() == "OPTIONS" {
			return c.SendStatus(204)
		}
		
		return c.Next()
	}
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}