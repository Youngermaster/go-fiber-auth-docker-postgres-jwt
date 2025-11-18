package handler

import (
	"app/database"
	"app/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ProductResponse represents the product data returned to clients
type ProductResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
	UserID      uint   `json:"user_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// toProductResponse converts a Product model to ProductResponse
func toProductResponse(product *model.Product) ProductResponse {
	return ProductResponse{
		ID:          product.ID,
		Title:       product.Title,
		Description: product.Description,
		Amount:      product.Amount,
		UserID:      product.UserID,
		CreatedAt:   product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// GetAllProducts retrieves all products with pagination
func GetAllProducts(c *fiber.Ctx) error {
	db := database.DB
	var products []model.Product

	// Get and validate pagination parameters
	page, limit := GetPaginationParams(c)
	offset := CalculateOffset(page, limit)

	// Get total count
	var total int64
	db.Model(&model.Product{}).Count(&total)

	// Get products with pagination
	if err := db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&products).Error; err != nil {
		return ErrorResponseJSON(c, fiber.StatusInternalServerError, "Failed to fetch products", nil)
	}

	// Convert to response format
	productResponses := make([]ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = toProductResponse(&product)
	}

	return SuccessResponse(c, "Products retrieved successfully", fiber.Map{
		"products": productResponses,
		"meta": PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: CalculateTotalPages(total, limit),
		},
	})
}

// GetProduct retrieves a single product by ID
func GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var product model.Product

	if err := db.First(&product, id).Error; err != nil {
		return ErrorResponseJSON(c, fiber.StatusNotFound, "Product not found", nil)
	}

	return SuccessResponse(c, "Product found", toProductResponse(&product))
}

// CreateProduct creates a new product
func CreateProduct(c *fiber.Ctx) error {
	type ProductInput struct {
		Title       string `json:"title" validate:"required,min=1,max=255"`
		Description string `json:"description" validate:"required,min=1,max=2000"`
		Amount      int    `json:"amount" validate:"required,min=0"`
	}

	input := new(ProductInput)
	if err := c.BodyParser(input); err != nil {
		return ErrorResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate input
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return ErrorResponseJSON(c, fiber.StatusBadRequest, "Validation failed", err.Error())
	}

	// Sanitize input
	input.Title = SanitizeString(input.Title, 255)
	input.Description = SanitizeString(input.Description, 2000)

	// Get user ID from JWT token
	userID, err := GetUserIDFromToken(c)
	if err != nil {
		return ErrorResponseJSON(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	product := model.Product{
		Title:       input.Title,
		Description: input.Description,
		Amount:      input.Amount,
		UserID:      userID,
	}

	db := database.DB
	if err := db.Create(&product).Error; err != nil {
		return ErrorResponseJSON(c, fiber.StatusInternalServerError, "Failed to create product", nil)
	}

	return CreatedResponse(c, "Product created successfully", toProductResponse(&product))
}

// UpdateProduct updates an existing product
func UpdateProduct(c *fiber.Ctx) error {
	type UpdateProductInput struct {
		Title       string `json:"title" validate:"omitempty,min=1,max=255"`
		Description string `json:"description" validate:"omitempty,min=1,max=2000"`
		Amount      int    `json:"amount" validate:"omitempty,min=0"`
	}

	id := c.Params("id")
	db := database.DB

	// Get user ID from token
	userID, err := GetUserIDFromToken(c)
	if err != nil {
		return ErrorResponseJSON(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	// Find product
	var product model.Product
	if err := db.First(&product, id).Error; err != nil {
		return ErrorResponseJSON(c, fiber.StatusNotFound, "Product not found", nil)
	}

	// Verify ownership
	if product.UserID != userID {
		return ErrorResponseJSON(c, fiber.StatusForbidden, "You don't have permission to update this product", nil)
	}

	// Parse input
	input := new(UpdateProductInput)
	if err := c.BodyParser(input); err != nil {
		return ErrorResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate input
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return ErrorResponseJSON(c, fiber.StatusBadRequest, "Validation failed", err.Error())
	}

	// Update only provided fields
	if input.Title != "" {
		product.Title = SanitizeString(input.Title, 255)
	}
	if input.Description != "" {
		product.Description = SanitizeString(input.Description, 2000)
	}
	if input.Amount >= 0 {
		product.Amount = input.Amount
	}

	if err := db.Save(&product).Error; err != nil {
		return ErrorResponseJSON(c, fiber.StatusInternalServerError, "Failed to update product", nil)
	}

	return SuccessResponse(c, "Product updated successfully", toProductResponse(&product))
}

// DeleteProduct deletes a product
func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	// Get user ID from token
	userID, err := GetUserIDFromToken(c)
	if err != nil {
		return ErrorResponseJSON(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	var product model.Product
	if err := db.First(&product, id).Error; err != nil {
		return ErrorResponseJSON(c, fiber.StatusNotFound, "Product not found", nil)
	}

	// Verify ownership
	if product.UserID != userID {
		return ErrorResponseJSON(c, fiber.StatusForbidden, "You don't have permission to delete this product", nil)
	}

	// Soft delete (GORM's default with gorm.Model)
	if err := db.Delete(&product).Error; err != nil {
		return ErrorResponseJSON(c, fiber.StatusInternalServerError, "Failed to delete product", nil)
	}

	return SuccessResponse(c, "Product deleted successfully", nil)
}
