package framework

import (
	"context"
	"github.com/test/myapp/config"
	"github.com/test/myapp/framework/cache"
	"github.com/test/myapp/framework/container"
	"github.com/test/myapp/framework/database"
	"github.com/test/myapp/framework/docs"
	"github.com/test/myapp/framework/errors"
	"github.com/test/myapp/framework/events"
	"github.com/test/myapp/framework/middleware"
	"github.com/test/myapp/framework/queue"
	"github.com/test/myapp/framework/validation"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/redis/go-redis/v9"
)

// Golara represents the main framework instance
type Golara struct {
	App        *fiber.App
	Container  *container.Container
	DB         *database.DatabaseManager
	Cache      *cache.CacheManager
	Queue      *queue.QueueManager
	Events     *events.EventDispatcher
	Middleware *middleware.MiddlewareRegistry
	Docs       *docs.DocGenerator
	Validator  *validation.Validator
}

// New creates a new Golara framework instance
func New(config ...Config) *Golara {
	cfg := defaultConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	app := fiber.New(fiber.Config{
		BodyLimit:    cfg.BodyLimit,
		ErrorHandler: errors.ErrorHandler,
	})

	// Initialize all components
	containerInstance := container.NewContainer()
	dbManager := database.NewDatabaseManager()
	cacheManager := cache.NewCacheManager()
	queueManager := queue.NewQueueManager()
	eventsDispatcher := events.NewEventDispatcher()
	middlewareRegistry := middleware.NewMiddlewareRegistry()
	docGenerator := docs.NewDocGenerator(cfg.AppName, cfg.Version)
	validator := validation.NewValidator()

	golara := &Golara{
		App:        app,
		Container:  containerInstance,
		DB:         dbManager,
		Cache:      cacheManager,
		Queue:      queueManager,
		Events:     eventsDispatcher,
		Middleware: middlewareRegistry,
		Docs:       docGenerator,
		Validator:  validator,
	}

	golara.setupServices()
	golara.setupDefaultMiddleware()
	golara.setupDocumentation()

	return golara
}

// Config holds framework configuration
type Config struct {
	AppName     string
	Version     string
	BodyLimit   int
	Debug       bool
	RedisAddr   string
	RedisPass   string
	RedisDB     int
	Environment string
}

func defaultConfig() Config {
	return Config{
		AppName:     "Golara Application",
		Version:     "1.0.0",
		BodyLimit:   20 * 1024 * 1024, // 20MB
		Debug:       true,
		RedisAddr:   "localhost:6379",
		RedisPass:   "",
		RedisDB:     0,
		Environment: "development",
	}
}

func (g *Golara) setupServices() {
	// Setup memory cache (Redis optional)
	memoryCache := cache.NewMemoryCache("golara")
	g.Cache.AddStore("memory", memoryCache)

	// Setup default memory queue
	memoryQueue := &queue.MemoryQueue{}
	g.Queue.AddQueue("default", memoryQueue)

	// Try Redis setup (optional)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err == nil {
		redisCache := cache.NewRedisCache("localhost:6379", "", 0, "golara")
		g.Cache.AddStore("redis", redisCache)

		// Setup Redis queue (replace default)
		redisQueue := queue.NewRedisQueue(redisClient, "default")
		g.Queue.AddQueue("default", redisQueue)
	}

	// Register services in container
	g.Container.Instance("app", g.App)
	g.Container.Instance("db", g.DB)
	g.Container.Instance("cache", g.Cache)
	g.Container.Instance("queue", g.Queue)
	g.Container.Instance("events", g.Events)
	g.Container.Instance("validator", g.Validator)
}

func (g *Golara) setupDefaultMiddleware() {
	// Request ID middleware
	g.App.Use(middleware.RequestID())

	// Logger middleware
	g.App.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	// Recovery middleware
	g.App.Use(recover.New())

	// CORS middleware
	corsOrigins := config.Denv("CORS_ALLOWED_DOMAINS")
	if corsOrigins == "" {
		corsOrigins = "http://localhost:3000,http://localhost:8080"
	}
	g.App.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins,
		AllowHeaders:     "Content-Type, Authorization, X-API-Key, X-Request-ID",
		AllowCredentials: false,
	}))
}

func (g *Golara) setupDocumentation() {
	g.Docs.ServeSwaggerUI(g.App)
}

// RegisterRoute registers a route and adds it to documentation
func (g *Golara) RegisterRoute(method, path, summary, description string, handler fiber.Handler, tags ...string) {
	switch method {
	case "GET":
		g.App.Get(path, handler)
	case "POST":
		g.App.Post(path, handler)
	case "PUT":
		g.App.Put(path, handler)
	case "DELETE":
		g.App.Delete(path, handler)
	}

	g.Docs.AddRoute(method, path, summary, description, tags)
}

// Group creates a route group with middleware
func (g *Golara) Group(prefix string, middleware ...fiber.Handler) fiber.Router {
	return g.App.Group(prefix, middleware...)
}

// ConnectDatabase connects to database
func (g *Golara) ConnectDatabase(name string, config database.DatabaseConfig) error {
	return g.DB.Connect(name, config)
}

// StartQueue starts queue workers
func (g *Golara) StartQueue(queueName string, concurrency int) {
	g.Queue.StartWorker(queueName, concurrency)
}

// Listen starts the server
func (g *Golara) Listen(addr string) error {
	log.Printf("🚀 Golara server starting on %s", addr)
	log.Printf("📚 API Documentation available at http://localhost%s/docs", addr)
	log.Printf("🗄️  Database connections: %v", g.DB)
	log.Printf("💾 Cache stores available: redis, memory")
	log.Printf("⚡ Queue workers ready")
	return g.App.Listen(addr)
}

// Shutdown gracefully shuts down the server
func (g *Golara) Shutdown() error {
	log.Println("🛑 Shutting down Golara server...")
	return g.App.Shutdown()
}
