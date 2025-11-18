package router

import (
	"app/handler"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// SetupAuthRoutes configures all authentication-related routes
func SetupAuthRoutes(router fiber.Router) {
	auth := router.Group("/auth")

	// Rate limiter: 5 login attempts per minute per IP
	loginLimiter := limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  "error",
				"message": "Too many login attempts. Please try again later.",
			})
		},
	})

	// Public routes
	auth.Post("/login", loginLimiter, handler.Login)
	// TODO: Add these routes as you implement them
	// auth.Post("/register", registerLimiter, handler.Register)
	// auth.Post("/refresh", handler.RefreshToken)
	// auth.Post("/logout", middleware.Protected(), handler.Logout)
	// auth.Post("/forgot-password", handler.ForgotPassword)
	// auth.Post("/reset-password", handler.ResetPassword)
	// auth.Post("/verify-email", handler.VerifyEmail)
}
