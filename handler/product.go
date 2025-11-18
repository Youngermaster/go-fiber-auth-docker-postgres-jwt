package handler

import (
	"app/database"
	"app/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// GetAllProducts query all products with pagination
func GetAllProducts(c *fiber.Ctx) error {
	db := database.DB
	var products []model.Product

	// Pagination parameters
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Get total count
	var total int64
	db.Model(&model.Product{}).Count(&total)

	// Get products with pagination
	db.Limit(limit).Offset(offset).Find(&products)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All products",
		"data":    products,
		"meta": fiber.Map{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GetProduct query product
func GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var product model.Product
	db.Find(&product, id)
	if product.Title == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No product found with ID", "data": nil})

	}
	return c.JSON(fiber.Map{"status": "success", "message": "Product found", "data": product})
}

// CreateProduct new product
func CreateProduct(c *fiber.Ctx) error {
	type ProductInput struct {
		Title       string `json:"title" validate:"required,min=1,max=255"`
		Description string `json:"description" validate:"required,min=1,max=2000"`
		Amount      int    `json:"amount" validate:"required,min=0"`
	}

	db := database.DB
	input := new(ProductInput)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body", "errors": err.Error()})
	}

	// Validate input
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Validation failed", "errors": err.Error()})
	}

	// Get user ID from JWT token
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	product := model.Product{
		Title:       input.Title,
		Description: input.Description,
		Amount:      input.Amount,
		UserID:      userID,
	}

	if err := db.Create(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to create product", "errors": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": "Created product", "data": product})
}

// DeleteProduct delete product
func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	// Get user ID from JWT token
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	var product model.Product
	db.First(&product, id)
	if product.Title == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "No product found with ID", "data": nil})
	}

	// Check if the user owns the product
	if product.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "You don't have permission to delete this product", "data": nil})
	}

	db.Delete(&product)
	return c.JSON(fiber.Map{"status": "success", "message": "Product successfully deleted", "data": nil})
}
