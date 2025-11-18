package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Response represents a standard API response structure
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents a standard error response structure
type ErrorResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}

// SuccessResponse sends a standardized success response
func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.JSON(Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// CreatedResponse sends a standardized 201 created response
func CreatedResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// ErrorResponseJSON sends a standardized error response
func ErrorResponseJSON(c *fiber.Ctx, statusCode int, message string, errors interface{}) error {
	return c.Status(statusCode).JSON(ErrorResponse{
		Status:  "error",
		Message: message,
		Errors:  errors,
	})
}

// GetUserIDFromToken extracts user ID from JWT token in context
func GetUserIDFromToken(c *fiber.Ctx) (uint, error) {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))
	return userID, nil
}

// ValidateTokenOwnership checks if the token's user ID matches the resource owner
func ValidateTokenOwnership(c *fiber.Ctx, resourceID string) bool {
	id, err := strconv.ParseUint(resourceID, 10, 32)
	if err != nil {
		return false
	}

	userID, err := GetUserIDFromToken(c)
	if err != nil {
		return false
	}

	return uint(id) == userID
}

// GetPaginationParams extracts and validates pagination parameters from query
func GetPaginationParams(c *fiber.Ctx) (page int, limit int) {
	page = c.QueryInt("page", 1)
	limit = c.QueryInt("limit", 10)

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return page, limit
}

// CalculateOffset calculates the database offset from page and limit
func CalculateOffset(page, limit int) int {
	return (page - 1) * limit
}

// CalculateTotalPages calculates total pages from total items and limit
func CalculateTotalPages(total int64, limit int) int64 {
	return (total + int64(limit) - 1) / int64(limit)
}
