package main

import (
	"app/config"
	"app/database"
	"app/router"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Validate configuration before starting the application
	// This ensures all required environment variables are set and valid
	if validationErrors := config.ValidateConfig(); len(validationErrors) > 0 {
		config.PrintValidationErrors(validationErrors)
		fmt.Println("❌ Application cannot start due to configuration errors.")
		fmt.Println("   Please fix the issues above and restart the application.")
		os.Exit(1)
	}

	fmt.Println("✅ Configuration validated successfully")

	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Go Fiber Auth API",
	})

	// CORS middleware with secure defaults
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:3001",
		AllowMethods:     "GET,POST,PATCH,DELETE",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
		MaxAge:           300,
	}))

	database.ConnectDB()

	router.SetupRoutes(app)
	log.Fatal(app.Listen(":3000"))
}
