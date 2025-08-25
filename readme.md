# 🔥 Golara Framework

### The Laravel-inspired Go Framework

A complete Laravel-style MVC framework built on top of Fiber and GORM, providing familiar Laravel conventions with Go's performance benefits.

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/tests-passing-brightgreen.svg)](#testing)

## ✨ Features

- **🏗️ Laravel-style MVC Architecture** - Controllers, Models, Views pattern
- **🗄️ Eloquent-style ORM** - Fluent query builder with relationships
- **🔄 Multi-Database Support** - MySQL, PostgreSQL, SQLite
- **⚡ Redis Integration** - Caching and queue management
- **📋 Job Queue System** - Background job processing
- **✅ Request Validation** - Laravel-style validation rules
- **📁 File Storage System** - Laravel-style file operations
- **📡 Event System** - Observer pattern implementation
- **🏗️ Service Container** - Dependency injection
- **🌍 Environment Configuration** - Multi-environment support
- **📚 Auto API Documentation** - Swagger/OpenAPI generation
- **🧪 Testing Framework** - HTTP testing utilities
- **⚙️ CLI Scaffolding** - Laravel-style artisan commands
- **🛡️ Security Features** - JWT auth, rate limiting, CORS
- **🔧 Middleware System** - Request/response pipeline

## 🚀 Quick Start

### 1. Installation

```bash
# Clone the framework
git clone https://github.com/your-repo/golara
cd golara

# Initialize your project (dynamic setup)
go run main.go -subcommand init

# Install dependencies
go mod tidy

# Build the application
make build
```

### 2. Project Setup

The `init` command will prompt you for:
- **Module Name**: Your Go module (e.g., `github.com/username/myapp`)
- **App Name**: Your application name
- **Database Name**: Your database name
- **Author**: Your name (optional)

This automatically configures all imports and creates your custom `.env.yaml`.

### 3. Basic Usage

```go
package main

import (
    "yourmodule/framework"
    "github.com/gofiber/fiber/v2"
)

func main() {
    app := framework.New(framework.Config{
        AppName: "My API",
        Version: "1.0.0",
    })
    
    app.App.Get("/", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"message": "Hello Golara!"})
    })
    
    app.Listen(":9000")
}
```

## 🏗️ Laravel-style MVC Architecture

### Controllers (Resource Controllers)

Generate Laravel-style resource controllers:

```bash
./bin/golara -subcommand make:controller User
```

Generated controller includes all REST methods:

```go
// UserController handles user related requests
type UserController struct {
    DB *database.DatabaseManager
}

// Index - GET /users (list all)
func (ctrl *UserController) Index(c *fiber.Ctx) error { ... }

// Show - GET /users/:id (show one)
func (ctrl *UserController) Show(c *fiber.Ctx) error { ... }

// Store - POST /users (create new)
func (ctrl *UserController) Store(c *fiber.Ctx) error { ... }

// Update - PUT /users/:id (update existing)
func (ctrl *UserController) Update(c *fiber.Ctx) error { ... }

// Destroy - DELETE /users/:id (delete)
func (ctrl *UserController) Destroy(c *fiber.Ctx) error { ... }
```

### Models (Eloquent-style)

Generate models with relationships:

```bash
./bin/golara -subcommand make:model Product
```

```go
type Product struct {
    gorm.Model
    Name        string  `json:"name" gorm:"size:255;not null"`
    Price       float64 `json:"price"`
    CategoryID  uint    `json:"category_id"`
    Category    Category `json:"category" gorm:"foreignKey:CategoryID"`
}

// Eloquent-style queries
products := []Product{}
qb.Table("products").
   Where("price", ">", 100).
   With("Category").
   OrderBy("created_at", "DESC").
   Paginate(1, 15, &products)
```

### Migrations

Laravel-style database migrations:

```bash
./bin/golara -subcommand make:migration create_products_table
./bin/golara -subcommand migrate
./bin/golara -subcommand migrate:rollback
```

## 🛠️ CLI Commands (Laravel Artisan-style)

```bash
# Project setup
./bin/golara -subcommand init                    # Initialize new project

# Database operations
./bin/golara -subcommand migrate                 # Run migrations
./bin/golara -subcommand migrate:rollback        # Rollback migrations
./bin/golara -subcommand migrate:status          # Migration status

# Code generation (MVC scaffolding)
./bin/golara -subcommand make:controller User    # Generate controller
./bin/golara -subcommand make:model Product      # Generate model
./bin/golara -subcommand make:middleware Auth    # Generate middleware
./bin/golara -subcommand make:job SendEmail      # Generate job
./bin/golara -subcommand make:migration users    # Generate migration

# Development
make dev                                         # Hot reload server
make test                                        # Run all tests
make docs                                        # Generate documentation
```

## 🌟 Laravel-style Features

### Request Validation

```go
// Laravel-style validation
validator := validation.NewValidator()
validator.AddRule("email", validation.Required()).
         AddRule("email", validation.Email()).
         AddRule("password", validation.Min(8))

if !validator.Validate(c) {
    return c.Status(422).JSON(fiber.Map{
        "errors": validator.GetErrors(),
    })
}
```

### Job Queues

```go
// Background job processing (Laravel-style)
job := jobs.NewSendEmailJob(email, subject, body)
app.Queue.Dispatch(job)

// Job definition
type SendEmailJob struct {
    queue.BaseJob
    Email   string
    Subject string
    Body    string
}

func (j *SendEmailJob) Handle() error {
    // Process job
    return nil
}
```

### Events & Listeners

```go
// Event-driven architecture
event := events.NewUserRegisteredEvent(user)
app.Events.Dispatch(event)

// Event listener
app.Events.ListenFunc("user.registered", func(event events.Event) error {
    // Handle event
    return nil
})
```

### Middleware

```go
// JWT Authentication
app.App.Use("/api", middleware.JWT(middleware.JWTConfig{
    SecretKey: "your-secret",
}))

// Rate Limiting
app.App.Use(middleware.RateLimit(middleware.RateLimitConfig{
    Max:        100,
    Expiration: time.Hour,
}))
```

### Caching (Laravel-style)

```go
// Remember pattern
cache := app.Cache.Store("redis")
cache.Remember("users", 1*time.Hour, func() (interface{}, error) {
    return fetchUsers()
}, &users)
```

### File Storage (Laravel-style)

```go
// File operations
storage := app.Storage.Disk("local")
storage.Put("files/document.pdf", content)
storage.Get("files/document.pdf")
storage.Delete("files/document.pdf")

// File uploads
file, _ := c.FormFile("upload")
path, _ := storage.Store(file, "uploads")
url := storage.URL(path)
```

## 📁 Project Structure (Laravel-inspired)

```
yourapp/
├── app/                    # Application layer
│   ├── controllers/        # HTTP controllers (MVC)
│   ├── models/            # Database models (MVC)
│   ├── jobs/              # Background jobs
│   ├── middleware/        # Custom middleware
│   └── providers/         # Service providers
├── config/                # Configuration files
├── database/              # Database layer
│   └── migrations/        # Database migrations
├── framework/             # Framework core
│   ├── cache/            # Caching system
│   ├── database/         # ORM & Query Builder
│   ├── storage/          # File storage system
│   ├── queue/            # Job queue system
│   ├── events/           # Event system
│   ├── validation/       # Validation system
│   ├── middleware/       # Built-in middleware
│   └── cli/              # CLI tools
├── routes/               # Route definitions
├── storage/              # File storage
├── tests/                # Test suites
│   ├── unit/            # Unit tests
│   ├── integration/     # Integration tests
│   └── examples/        # Example tests
├── .env.yaml            # Environment configuration
├── main.go              # Application entry point
└── Makefile             # Build automation
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run specific test suites
make test-unit           # Unit tests
make test-integration    # Integration tests
make test-examples       # Example tests
```

### Test Coverage
- ✅ Framework Bootstrap & Configuration
- ✅ MVC Architecture (Controllers, Models)
- ✅ Database ORM & Query Builder
- ✅ Cache System (Redis & Memory)
- ✅ Job Queue System
- ✅ Event System & Listeners
- ✅ Request Validation
- ✅ HTTP Routing & Middleware
- ✅ Service Container & DI
- ✅ CLI Tools & Generators

## 🔧 Configuration

### Environment Configuration (.env.yaml)

```yaml
# Application
APP_NAME: "My Golara App"
APP_ENV: "development"
APP_PORT: "9000"

# Database
DB_CONNECTION: "mysql"
DB_HOST: "localhost"
DB_DATABASE: "myapp_db"
DB_USERNAME: "root"
DB_PASSWORD: ""

# Redis
REDIS_HOST: "localhost"
REDIS_PORT: "6379"

# Cache & Queue
CACHE_DRIVER: "redis"
QUEUE_CONNECTION: "redis"

# File Storage
STORAGE_DRIVER: "local"
STORAGE_URL: "/storage"
```

### Database Support
- **MySQL** - Primary database with full feature support
- **PostgreSQL** - Complete PostgreSQL integration
- **SQLite** - Lightweight database for development

### Cache & Queue Drivers
- **Redis** - Distributed caching and job processing
- **Memory** - In-memory caching and synchronous jobs

## 🚀 Production Deployment

### Docker Support

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/golara .
CMD ["./golara"]
```

### Environment Setup

```bash
# Production environment
export APP_ENV=production
export APP_DEBUG=false

# Database configuration
export DB_HOST=your-db-host
export DB_PASSWORD=your-secure-password

# Redis configuration
export REDIS_HOST=your-redis-host
export REDIS_PASSWORD=your-redis-password
```

## 📚 Documentation

- **[API Documentation](http://localhost:9000/docs)** - Auto-generated Swagger docs
- **[HTML Documentation](docs/html/index.html)** - Complete developer guide
- **[Examples](examples/)** - Working code examples

## 🎯 Why Choose Golara?

### Laravel Familiarity + Go Performance

- **Familiar Syntax**: Laravel developers feel at home
- **Go Performance**: Compiled binary, fast execution
- **Type Safety**: Go's strong typing prevents runtime errors
- **Concurrency**: Built-in goroutine support
- **Single Binary**: Easy deployment, no dependencies

### Production Ready

- **Security First**: JWT auth, rate limiting, input validation
- **Scalable**: Redis clustering, database pooling
- **Observable**: Structured logging, error tracking
- **Testable**: Comprehensive testing framework
- **Maintainable**: Clean architecture, dependency injection

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Write tests for your changes
4. Ensure all tests pass (`make test`)
5. Commit changes (`git commit -m 'Add amazing feature'`)
6. Push to branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by Laravel's elegant syntax and conventions
- Built on top of [Fiber](https://github.com/gofiber/fiber) and [GORM](https://github.com/go-gorm/gorm)
- Thanks to the Go and Laravel communities

## 🎯 Complete MVC Implementation

### **M - Models** ✅
```bash
./bin/golara -subcommand make:model Product
```
- Eloquent-style models with GORM
- Relationships and migrations
- Database query builder

### **V - Views** ✅  
```bash
./bin/golara -subcommand make:view products/index
```
- Laravel Blade-like template engine
- Layout support: `resources/views/layouts/app.html`
- Helper functions: `{{asset}}`, `{{route}}`, `{{csrf}}`
- Automatic layout generation

### **C - Controllers** ✅
```bash
# API Controller (JSON responses)
./bin/golara -subcommand make:controller Product

# Web Controller (HTML views)  
./bin/golara -subcommand make:controller Product --web
```

**API Controller Methods:**
- `Index()` - GET /products (JSON list)
- `Show()` - GET /products/:id (JSON single)
- `Store()` - POST /products (JSON create)
- `Update()` - PUT /products/:id (JSON update)
- `Destroy()` - DELETE /products/:id (JSON delete)

**Web Controller Methods:**
- `Index()` - GET /products (HTML list)
- `Show()` - GET /products/:id (HTML single)
- `Create()` - GET /products/create (HTML form)
- `Store()` - POST /products (form submit)
- `Edit()` - GET /products/:id/edit (HTML form)
- `Update()` - PUT /products/:id (form submit)
- `Destroy()` - DELETE /products/:id (redirect)

## 🏗️ Complete Project Structure

```
myapp/
├── app/
│   ├── controllers/           # MVC Controllers
│   │   ├── product_controller.go    # API controller
│   │   └── blog_controller.go       # Web controller  
│   ├── models/               # MVC Models
│   │   └── product.go
│   ├── jobs/                 # Background jobs
│   └── middleware/           # Custom middleware
├── resources/
│   └── views/                # MVC Views
│       ├── layouts/
│       │   └── app.html      # Master layout
│       ├── welcome.html      # Home page
│       └── products/
│           ├── index.html    # List view
│           ├── show.html     # Detail view
│           ├── create.html   # Create form
│           └── edit.html     # Edit form
├── public/                   # Static assets
│   ├── css/app.css
│   ├── js/app.js
│   └── images/
├── routes/                   # Route definitions
├── database/migrations/      # Database migrations
├── framework/                # Framework core
│   ├── view/                 # Template engine
│   ├── database/             # ORM & Query Builder
│   ├── cache/                # Caching system
│   └── queue/                # Job system
├── tests/                    # Test suites
│   ├── unit/                # Unit tests
│   ├── integration/         # Integration tests
│   └── examples/            # Example tests
├── examples/                 # Usage examples
├── .env.yaml                # Configuration
└── docker-compose.yml       # Docker setup
```

## 🎨 Template Engine Features

### **Laravel Blade-like Syntax**
```html
<!-- resources/views/products/index.html -->
{{define "content"}}
<h1>{{.title}}</h1>
{{range .products}}
    <div class="product">
        <h2>{{.Name}}</h2>
        <p>Price: ${{.Price}}</p>
        <a href="{{route "products.show"}} {{.ID}}">View</a>
    </div>
{{end}}
{{end}}
```

### **Web Controller with Views**
```go
// BlogController (Web)
func (ctrl *BlogController) Index(c *fiber.Ctx) error {
    var blogs []models.Blog
    // ... fetch data
    
    return view.View(c, "blogs/index", view.ViewData{
        "blogs": blogs,
        "title": "All Blogs",
    }, "layouts/app")
}
```

### **Helper Functions**
- `{{asset "css/app.css"}}` - Asset URLs
- `{{route "home"}}` - Named routes
- `{{csrf}}` - CSRF token
- `{{.title}}` - Data binding

## 🐳 Docker Production Setup

```bash
# Development
docker-compose up

# Production
docker-compose --profile production up
```

**Features:**
- Multi-stage build for security
- Health checks
- Volume persistence
- Environment variables

## 📊 Laravel vs Golara

| Feature | Laravel | Golara |
|---------|---------|--------|
| **Models** | Eloquent ORM | GORM + Query Builder |
| **Views** | Blade Templates | Go Templates + Helpers |
| **Controllers** | Resource Controllers | API + Web Controllers |
| **Routing** | Route::resource() | Fiber routing |
| **Migrations** | php artisan migrate | golara migrate |
| **Jobs** | Queue::dispatch() | app.Queue.Dispatch() |
| **Storage** | Storage::disk() | app.Storage.Disk() |
| **Performance** | PHP (interpreted) | Go (compiled) |

## 📄 License

MIT License - Free for commercial and personal use.

---

**Made with ❤️ for developers who love Laravel but need Go's performance**

**Complete MVC Framework Ready! 🔥**