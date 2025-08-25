package main

import (
	"github.com/test/myapp/app/jobs"
	"github.com/test/myapp/cmd"
	"github.com/test/myapp/config"
	"github.com/test/myapp/framework"
	"github.com/test/myapp/framework/database"
	"github.com/test/myapp/framework/queue"
	"github.com/test/myapp/routes"
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
		port = "8080"
	}
	
	log.Printf("🔥 %s starting on port %s", config.Denv("APP_NAME"), port)
	log.Printf("📚 API Documentation: http://localhost:%s/docs", port)
	
	app.Listen(":" + port)
}
