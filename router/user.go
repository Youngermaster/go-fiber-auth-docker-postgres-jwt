package router

import (
	"app/handler"
	"app/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupUserRoutes configures all user-related routes
func SetupUserRoutes(router fiber.Router) {
	users := router.Group("/users")

	// Public routes
	users.Post("/", handler.CreateUser) // Registration

	// Protected routes - require authentication
	users.Get("/:id", middleware.Protected(), handler.GetUser)
	users.Patch("/:id", middleware.Protected(), handler.UpdateUser)
	users.Delete("/:id", middleware.Protected(), handler.DeleteUser)

	// TODO: Add additional user routes as needed
	// users.Get("/", middleware.Protected(), middleware.AdminOnly(), handler.GetAllUsers)
	// users.Get("/me", middleware.Protected(), handler.GetCurrentUser)
	// users.Patch("/me/password", middleware.Protected(), handler.ChangePassword)
	// users.Get("/me/sessions", middleware.Protected(), handler.GetUserSessions)
	// users.Delete("/me/sessions/:id", middleware.Protected(), handler.RevokeSession)
}
