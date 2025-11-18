package router

import (
	"app/handler"
	"app/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// SetupAuthRoutes configures all authentication-related routes
func SetupAuthRoutes(router fiber.Router) {
	auth := router.Group("/auth")

	// Rate limiter for authentication endpoints: 5 attempts per minute per IP
	authLimiter := limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  "error",
				"message": "Too many requests. Please try again later.",
			})
		},
	})

	// Public routes (no authentication required)
	auth.Post("/login", authLimiter, handler.Login)
	auth.Post("/refresh", authLimiter, handler.RefreshToken) // Rate limited to prevent abuse

	// Protected routes (require valid JWT access token)
	auth.Post("/logout", middleware.Protected(), handler.Logout)
	auth.Post("/logout-all", middleware.Protected(), handler.LogoutAll)
	auth.Get("/sessions", middleware.Protected(), handler.GetActiveSessions)

	// TODO: Add these routes when implementing additional auth features
	// auth.Post("/register", authLimiter, handler.Register)  // User registration
	// auth.Post("/forgot-password", authLimiter, handler.ForgotPassword)
	// auth.Post("/reset-password", authLimiter, handler.ResetPassword)
	// auth.Post("/verify-email", handler.VerifyEmail)
	// auth.Post("/resend-verification", authLimiter, handler.ResendVerification)
}
