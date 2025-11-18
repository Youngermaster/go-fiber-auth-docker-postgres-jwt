package handler

import (
	"app/config"
	"app/database"
	"app/model"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func getUserByEmail(e string) (*model.User, error) {
	db := database.DB
	var user model.User
	if err := db.Where(&model.User{Email: e}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func getUserByUsername(u string) (*model.User, error) {
	db := database.DB
	var user model.User
	if err := db.Where(&model.User{Username: u}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Login authenticates a user and returns a JWT token
func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Identity string `json:"identity" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	input := new(LoginInput)
	if err := c.BodyParser(input); err != nil {
		return ErrorResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	identity := NormalizeEmail(input.Identity) // Normalize for consistency
	pass := input.Password

	var userModel *model.User
	var err error

	// Try to find user by email or username
	if ValidateEmail(identity) {
		userModel, err = getUserByEmail(identity)
	} else {
		userModel, err = getUserByUsername(NormalizeUsername(identity))
	}

	if err != nil {
		return ErrorResponseJSON(c, fiber.StatusInternalServerError, "Internal server error", nil)
	}

	// If user not found, still check password hash to prevent timing attacks
	if userModel == nil {
		CheckPasswordHash(pass, "")
		return ErrorResponseJSON(c, fiber.StatusUnauthorized, "Invalid credentials", nil)
	}

	// Verify password
	if !CheckPasswordHash(pass, userModel.Password) {
		return ErrorResponseJSON(c, fiber.StatusUnauthorized, "Invalid credentials", nil)
	}

	// Generate JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = userModel.Username
	claims["user_id"] = userModel.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // TODO: Reduce to 15-30 minutes and implement refresh tokens

	signedToken, err := token.SignedString([]byte(config.Config("SECRET")))
	if err != nil {
		return ErrorResponseJSON(c, fiber.StatusInternalServerError, "Failed to generate token", nil)
	}

	return SuccessResponse(c, "Login successful", fiber.Map{
		"token": signedToken,
		"user":  toUserResponse(userModel),
	})
}
