package middleware

import (
	"fmt"
	"github.com/test/myapp/framework/errors"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Max          int                     // Maximum requests
	Expiration   time.Duration           // Time window
	KeyFunc      func(*fiber.Ctx) string // Function to generate key
	LimitReached fiber.Handler           // Handler when limit is reached
	Storage      RateLimitStorage        // Storage interface
}

// RateLimitStorage interface for rate limit storage
type RateLimitStorage interface {
	Get(key string) (int, time.Time, bool)
	Set(key string, count int, expiration time.Time)
	Delete(key string)
}

// MemoryStorage implements in-memory rate limit storage
type MemoryStorage struct {
	data map[string]rateLimitEntry
	mu   sync.RWMutex
}

type rateLimitEntry struct {
	count      int
	expiration time.Time
}

// NewMemoryStorage creates new memory storage
func NewMemoryStorage() *MemoryStorage {
	storage := &MemoryStorage{
		data: make(map[string]rateLimitEntry),
	}

	// Cleanup expired entries every minute
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			storage.cleanup()
		}
	}()

	return storage
}

func (m *MemoryStorage) Get(key string) (int, time.Time, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, exists := m.data[key]
	if !exists || time.Now().After(entry.expiration) {
		return 0, time.Time{}, false
	}

	return entry.count, entry.expiration, true
}

func (m *MemoryStorage) Set(key string, count int, expiration time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = rateLimitEntry{
		count:      count,
		expiration: expiration,
	}
}

func (m *MemoryStorage) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data, key)
}

func (m *MemoryStorage) cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for key, entry := range m.data {
		if now.After(entry.expiration) {
			delete(m.data, key)
		}
	}
}

// DefaultRateLimitConfig returns default rate limit configuration
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		Max:        100,
		Expiration: time.Hour,
		KeyFunc: func(c *fiber.Ctx) string {
			return c.IP()
		},
		Storage: NewMemoryStorage(),
	}
}

// RateLimit returns rate limiting middleware
func RateLimit(config ...RateLimitConfig) fiber.Handler {
	cfg := DefaultRateLimitConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	return func(c *fiber.Ctx) error {
		key := cfg.KeyFunc(c)

		count, expiration, exists := cfg.Storage.Get(key)
		now := time.Now()

		if !exists || now.After(expiration) {
			// First request or expired window
			cfg.Storage.Set(key, 1, now.Add(cfg.Expiration))
			c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", cfg.Max))
			c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", cfg.Max-1))
			c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", now.Add(cfg.Expiration).Unix()))
			return c.Next()
		}

		if count >= cfg.Max {
			// Rate limit exceeded
			c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", cfg.Max))
			c.Set("X-RateLimit-Remaining", "0")
			c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", expiration.Unix()))

			if cfg.LimitReached != nil {
				return cfg.LimitReached(c)
			}

			return errors.NewAppError(errors.ValidationError, "Rate limit exceeded", 429)
		}

		// Increment counter
		cfg.Storage.Set(key, count+1, expiration)
		c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", cfg.Max))
		c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", cfg.Max-count-1))
		c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", expiration.Unix()))

		return c.Next()
	}
}
