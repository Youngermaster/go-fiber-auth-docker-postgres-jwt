package handler

import (
	"app/database"
	"app/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// UserResponse represents the user data returned to clients (without sensitive fields)
type UserResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Names     string `json:"names"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// toUserResponse converts a User model to UserResponse
func toUserResponse(user *model.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Names:     user.Names,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// GetUser retrieves a user by ID
func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var user model.User

	if err := db.First(&user, id).Error; err != nil {
		return ErrorResponseJSON(c, fiber.StatusNotFound, "User not found", nil)
	}

	return SuccessResponse(c, "User found", toUserResponse(&user))
}

// CreateUser creates a new user (registration)
func CreateUser(c *fiber.Ctx) error {
	type CreateUserInput struct {
		Username string `json:"username" validate:"required,min=3,max=50"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=100"`
		Names    string `json:"names" validate:"max=255"`
	}

	input := new(CreateUserInput)
	if err := c.BodyParser(input); err != nil {
		return ErrorResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate input
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return ErrorResponseJSON(c, fiber.StatusBadRequest, "Validation failed", err.Error())
	}

	// Normalize input
	input.Email = NormalizeEmail(input.Email)
	input.Username = NormalizeUsername(input.Username)
	input.Names = SanitizeString(input.Names, 255)

	// Validate password strength
	if err := ValidatePasswordStrength(input.Password); err != nil {
		return ErrorResponseJSON(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	// Hash password
	hashedPassword, err := HashPassword(input.Password)
	if err != nil {
		return ErrorResponseJSON(c, fiber.StatusInternalServerError, "Failed to process password", nil)
	}

	// Create user
	user := model.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,
		Names:    input.Names,
	}

	db := database.DB
	if err := db.Create(&user).Error; err != nil {
		// Check for duplicate email/username
		return ErrorResponseJSON(c, fiber.StatusConflict, "User with this email or username already exists", nil)
	}

	return CreatedResponse(c, "User created successfully", toUserResponse(&user))
}

// UpdateUser updates user information
func UpdateUser(c *fiber.Ctx) error {
	type UpdateUserInput struct {
		Names string `json:"names" validate:"max=255"`
	}

	id := c.Params("id")

	// Verify ownership - user can only update their own profile
	if !ValidateTokenOwnership(c, id) {
		return ErrorResponseJSON(c, fiber.StatusForbidden, "You don't have permission to update this user", nil)
	}

	input := new(UpdateUserInput)
	if err := c.BodyParser(input); err != nil {
		return ErrorResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate input
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return ErrorResponseJSON(c, fiber.StatusBadRequest, "Validation failed", err.Error())
	}

	// Sanitize input
	input.Names = SanitizeString(input.Names, 255)

	db := database.DB
	var user model.User

	if err := db.First(&user, id).Error; err != nil {
		return ErrorResponseJSON(c, fiber.StatusNotFound, "User not found", nil)
	}

	// Update only allowed fields
	user.Names = input.Names

	if err := db.Save(&user).Error; err != nil {
		return ErrorResponseJSON(c, fiber.StatusInternalServerError, "Failed to update user", nil)
	}

	return SuccessResponse(c, "User updated successfully", toUserResponse(&user))
}

// DeleteUser deletes a user account
func DeleteUser(c *fiber.Ctx) error {
	type DeleteUserInput struct {
		Password string `json:"password" validate:"required"`
	}

	id := c.Params("id")

	// Verify ownership - user can only delete their own account
	if !ValidateTokenOwnership(c, id) {
		return ErrorResponseJSON(c, fiber.StatusForbidden, "You don't have permission to delete this user", nil)
	}

	input := new(DeleteUserInput)
	if err := c.BodyParser(input); err != nil {
		return ErrorResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate input
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return ErrorResponseJSON(c, fiber.StatusBadRequest, "Validation failed", err.Error())
	}

	db := database.DB
	var user model.User

	if err := db.First(&user, id).Error; err != nil {
		return ErrorResponseJSON(c, fiber.StatusNotFound, "User not found", nil)
	}

	// Verify password before deletion (security measure)
	if !CheckPasswordHash(input.Password, user.Password) {
		return ErrorResponseJSON(c, fiber.StatusUnauthorized, "Invalid password", nil)
	}

	// Soft delete (GORM's default with gorm.Model)
	if err := db.Delete(&user).Error; err != nil {
		return ErrorResponseJSON(c, fiber.StatusInternalServerError, "Failed to delete user", nil)
	}

	return SuccessResponse(c, "User deleted successfully", nil)
}
