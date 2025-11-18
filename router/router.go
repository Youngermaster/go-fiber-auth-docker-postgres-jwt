package router

import (
	"app/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// SetupRoutes initializes all application routes in a modular, scalable way
// This is the main router that delegates to domain-specific route files
func SetupRoutes(app *fiber.App) {
	// Global middleware - applied to all routes
	app.Use(recover.New())   // Recover from panics
	app.Use(requestid.New()) // Add unique request ID for tracking
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${method} ${path} ${latency} [${locals:requestid}]\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "UTC",
	}))

	// API version 1 - all routes are versioned for future compatibility
	api := app.Group("/api/v1")

	// Root endpoint - API status
	api.Get("/", handler.Hello)

	// Health check endpoints (no versioning needed, used by orchestrators)
	SetupHealthRoutes(app)

	// Domain-specific routes - each in its own file for modularity
	SetupAuthRoutes(api)    // Authentication & authorization routes
	SetupUserRoutes(api)    // User management routes
	SetupProductRoutes(api) // Product/resource routes

	// TODO: Add more domain routes as your application grows
	// SetupOrderRoutes(api)
	// SetupPaymentRoutes(api)
	// SetupNotificationRoutes(api)
	// SetupAdminRoutes(api)

	// 404 Handler - must be last
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Route not found",
			"path":    c.Path(),
		})
	})
}

