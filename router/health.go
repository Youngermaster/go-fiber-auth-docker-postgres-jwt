package router

import (
	"app/database"
	"time"

	"github.com/gofiber/fiber/v2"
)

// SetupHealthRoutes configures health check endpoints for monitoring and orchestration
func SetupHealthRoutes(router fiber.Router) {
	health := router.Group("/health")

	// Basic health check - returns 200 if service is running
	health.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "healthy",
			"service":   "go-fiber-auth-api",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Detailed health check - checks database connectivity
	health.Get("/ready", func(c *fiber.Ctx) error {
		// Check database connection
		sqlDB, err := database.DB.DB()
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":    "unhealthy",
				"service":   "go-fiber-auth-api",
				"timestamp": time.Now().UTC().Format(time.RFC3339),
				"checks": fiber.Map{
					"database": "disconnected",
				},
			})
		}

		if err := sqlDB.Ping(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":    "unhealthy",
				"service":   "go-fiber-auth-api",
				"timestamp": time.Now().UTC().Format(time.RFC3339),
				"checks": fiber.Map{
					"database": "unreachable",
				},
			})
		}

		return c.JSON(fiber.Map{
			"status":    "healthy",
			"service":   "go-fiber-auth-api",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"checks": fiber.Map{
				"database": "connected",
			},
		})
	})

	// Liveness probe - for Kubernetes/Docker orchestration
	health.Get("/live", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "alive",
		})
	})
}
