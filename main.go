package main

import (
	"ekycapp/app/middlewares"
	"ekycapp/config"
	"ekycapp/routes"
	"ekycapp/utils"
	"os"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {

	config.InitDenv()
	config.ConnectDB()
	subCmd := "-subcommand"
	if len(os.Args) > 1 && os.Args[1] == subCmd {
		err := Commands()
		if err == nil {
			return
		}
		return
	}
	utils.StorageInit()

	app := fiber.New(fiber.Config{
		BodyLimit: 20 * 1024 * 1024, // Set limit to 10 MB
	})

	// Configure the logger to write to the log file
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
		Output:     utils.LogFile("request.log"),
		TimeFormat: "2006-01-02 15:04:05",
	}))

	if config.Env("APP_ENV") == "production" {
		f := utils.LogFile("error.log")
		defer f.Close()
		log.SetOutput(f) //If production then log on file otherwise in terminal
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.Denv("CORS_ALLOWED_DOMAINS"),
		AllowHeaders:     "Content-Type, Authorization, Entity",
		AllowCredentials: true,
	}))
	app.Use(recover.New())

	app.Use(middlewares.ErrorHandler)
	app.Use(middlewares.PanicHandler)
	//app.Use(middlewares.LicenseValidator)

	routes.WebRoutes(app)
	routes.ApiRoutes(app)

	app.Listen(":9000")
}
