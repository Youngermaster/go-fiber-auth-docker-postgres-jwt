package router

import (
	"app/handler"
	"app/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupProductRoutes configures all product-related routes
func SetupProductRoutes(router fiber.Router) {
	products := router.Group("/products")

	// Public routes
	products.Get("/", handler.GetAllProducts)
	products.Get("/:id", handler.GetProduct)

	// Protected routes - require authentication
	products.Post("/", middleware.Protected(), handler.CreateProduct)
	products.Patch("/:id", middleware.Protected(), handler.UpdateProduct)
	products.Delete("/:id", middleware.Protected(), handler.DeleteProduct)

	// TODO: Add additional product routes as needed
	// products.Get("/search", handler.SearchProducts)
	// products.Get("/categories/:category", handler.GetProductsByCategory)
	// products.Post("/:id/reviews", middleware.Protected(), handler.CreateProductReview)
}

// NOTE: You'll need to implement UpdateProduct handler
// func UpdateProduct(c *fiber.Ctx) error { ... }
