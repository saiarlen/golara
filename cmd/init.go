package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type ProjectConfig struct {
	ModuleName string
	AppName    string
	DBName     string
	Port       string
	Author     string
}

func InitProject() error {
	fmt.Println("🔥 Welcome to Golara Framework Setup")
	fmt.Println("=====================================")

	config := &ProjectConfig{}
	scanner := bufio.NewScanner(os.Stdin)

	// Get module name
	fmt.Print("Enter your Go module name (e.g., github.com/username/myapp): ")
	scanner.Scan()
	config.ModuleName = strings.TrimSpace(scanner.Text())
	if config.ModuleName == "" {
		config.ModuleName = "myapp"
	}

	// Get app name
	fmt.Print("Enter your application name: ")
	scanner.Scan()
	config.AppName = strings.TrimSpace(scanner.Text())
	if config.AppName == "" {
		config.AppName = "My Golara App"
	}

	// Get database name
	fmt.Print("Enter your database name: ")
	scanner.Scan()
	config.DBName = strings.TrimSpace(scanner.Text())
	if config.DBName == "" {
		config.DBName = "golara_db"
	}

	// Get port
	fmt.Print("Enter application port (default: 9000): ")
	scanner.Scan()
	config.Port = strings.TrimSpace(scanner.Text())
	if config.Port == "" {
		config.Port = "9000"
	}

	// Get author
	fmt.Print("Enter author name (optional): ")
	scanner.Scan()
	config.Author = strings.TrimSpace(scanner.Text())

	return setupProject(config)
}

func setupProject(config *ProjectConfig) error {
	fmt.Println("\n🚀 Setting up your Golara project...")

	// Update go.mod
	if err := updateGoMod(config); err != nil {
		return fmt.Errorf("failed to update go.mod: %w", err)
	}

	// Update imports in all files
	if err := updateImports(config); err != nil {
		return fmt.Errorf("failed to update imports: %w", err)
	}

	// Create .env file
	if err := createEnvFile(config); err != nil {
		return fmt.Errorf("failed to create .env file: %w", err)
	}

	// Create main.go
	if err := createMainFile(config); err != nil {
		return fmt.Errorf("failed to create main.go: %w", err)
	}

	// Create air.toml
	if err := createAirConfig(config); err != nil {
		return fmt.Errorf("failed to create air.toml: %w", err)
	}

	fmt.Println("✅ Project setup completed!")
	fmt.Printf("📁 Module: %s\n", config.ModuleName)
	fmt.Printf("🏷️  App: %s\n", config.AppName)
	fmt.Printf("🗄️  Database: %s\n", config.DBName)
	fmt.Printf("🚪 Port: %s\n", config.Port)
	fmt.Println("\n🎯 Next steps:")
	fmt.Println("1. Run: go mod tidy")
	fmt.Println("2. Update .env.yaml with your database credentials")
	fmt.Println("3. Run: make migrate")
	fmt.Println("4. Run: make dev (hot reload) or go run main.go")

	return nil
}

func updateGoMod(config *ProjectConfig) error {
	content := fmt.Sprintf(`module %s

go 1.21

require (
	github.com/gofiber/fiber/v2 v2.52.4
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/redis/go-redis/v9 v9.0.5
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/driver/mysql v1.5.1
	gorm.io/driver/postgres v1.5.2
	gorm.io/driver/sqlite v1.5.2
	gorm.io/gorm v1.25.10
)
`, config.ModuleName)

	return os.WriteFile("go.mod", []byte(content), 0644)
}

func updateImports(config *ProjectConfig) error {
	return filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(path, ".go") {
			return err
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Replace old module names with new module name
		newContent := strings.ReplaceAll(string(content), "ekycapp/", config.ModuleName+"/")
		newContent = strings.ReplaceAll(newContent, "github.com/test/myapp/", config.ModuleName+"/")
		newContent = strings.ReplaceAll(newContent, "github.com/test/myapp/", config.ModuleName+"/")

		return os.WriteFile(path, []byte(newContent), info.Mode())
	})
}

func createEnvFile(config *ProjectConfig) error {
	envTemplate := `# {{.AppName}} Configuration

# Application
APP_NAME: "{{.AppName}}"
APP_ENV: "development"
APP_DEBUG: true
APP_URL: "http://localhost:{{.Port}}"
APP_PORT: "{{.Port}}"

# Database
DB_CONNECTION: "mysql"
DB_HOST: "localhost"
DB_PORT: "3306"
DB_DATABASE: "{{.DBName}}"
DB_USERNAME: "root"
DB_PASSWORD: ""

# Redis
REDIS_HOST: "localhost"
REDIS_PORT: "6379"
REDIS_PASSWORD: ""
REDIS_DB: 0

# Cache & Queue
CACHE_DRIVER: "redis"
QUEUE_CONNECTION: "redis"

# CORS
CORS_ALLOWED_DOMAINS: "http://localhost:3000,http://localhost:8080,http://localhost:{{.Port}}"

# JWT
JWT_SECRET: "your-jwt-secret-key-change-this-in-production"
`

	tmpl, err := template.New("env").Parse(envTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(".env.yaml")
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, config)
}

func createMainFile(config *ProjectConfig) error {
	mainTemplate := `package main

import (
	"{{.ModuleName}}/app/jobs"
	"{{.ModuleName}}/cmd"
	"{{.ModuleName}}/config"
	"{{.ModuleName}}/framework"
	"{{.ModuleName}}/framework/database"
	"{{.ModuleName}}/framework/queue"
	"{{.ModuleName}}/routes"
	"log"
	"os"
)

func main() {
	// Initialize configuration
	if err := config.InitDenv(); err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}

	// Handle CLI commands
	subCmd := "-subcommand"
	if len(os.Args) > 1 && os.Args[1] == subCmd {
		// Connect to database only for migration commands
		if len(os.Args) > 2 && (os.Args[2] == "migrate" || os.Args[2] == "migrate:rollback" || os.Args[2] == "migrate:status") {
			if err := config.ConnectDB(); err != nil {
				log.Fatalf("Failed to connect to database: %v", err)
			}
		}
		if err := cmd.RunCommands(); err != nil {
			log.Fatal(err)
		}
		return
	}

	// Connect to database for server mode
	if err := config.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create Golara application
	app := framework.New(framework.Config{
		AppName:     config.Denv("APP_NAME"),
		Version:     "1.0.0",
		Environment: config.Denv("APP_ENV"),
		RedisAddr:   config.Denv("REDIS_HOST") + ":" + config.Denv("REDIS_PORT"),
		RedisPass:   config.Denv("REDIS_PASSWORD"),
		RedisDB:     0,
	})

	// Connect to database
	err := app.ConnectDatabase("mysql", database.DatabaseConfig{
		Driver:   config.Denv("DB_CONNECTION"),
		Host:     config.Denv("DB_HOST"),
		Port:     config.Denv("DB_PORT"),
		Database: config.Denv("DB_DATABASE"),
		Username: config.Denv("DB_USERNAME"),
		Password: config.Denv("DB_PASSWORD"),
		Charset:  "utf8mb4",
	})
	if err != nil {
		log.Printf("Warning: Database connection failed: %v", err)
	}

	// Register job handlers
	app.Queue.RegisterJob("send_email", func() queue.Job {
		return &jobs.SendEmailJob{}
	})

	// Start queue workers
	app.StartQueue("default", 3)

	// Register routes
	routes.RegisterRoutes(app)

	// Start server
	port := config.Denv("APP_PORT")
	if port == "" {
		port = "{{.Port}}"
	}
	
	log.Printf("🔥 %s starting on port %s", config.Denv("APP_NAME"), port)
	log.Printf("📚 API Documentation: http://localhost:%s/docs", port)
	
	app.Listen(":" + port)
}
`

	tmpl, err := template.New("main").Parse(mainTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create("main.go")
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, config)
}

func createAirConfig(config *ProjectConfig) error {
	airTemplate := `root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main main.go"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata", "storage", "node_modules", "bin"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html", "yaml"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
`

	return os.WriteFile(".air.toml", []byte(airTemplate), 0644)
}