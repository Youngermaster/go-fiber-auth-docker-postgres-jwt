package handler

import (
	"app/database"
	"app/model"
	"errors"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
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

// Login authenticates a user and returns access + refresh tokens
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

	// Generate access and refresh tokens
	tokenPair, err := GenerateTokenPair(userModel, c)
	if err != nil {
		return ErrorResponseJSON(c, fiber.StatusInternalServerError, "Failed to generate tokens", nil)
	}

	return SuccessResponse(c, "Login successful", fiber.Map{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
		"token_type":    tokenPair.TokenType,
		"expires_in":    tokenPair.ExpiresIn,
		"user":          toUserResponse(userModel),
	})
}

// RefreshToken exchanges a refresh token for a new access token
func RefreshToken(c *fiber.Ctx) error {
	type RefreshInput struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	input := new(RefreshInput)
	if err := c.BodyParser(input); err != nil {
		return ErrorResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Rotate the refresh token (revoke old, generate new pair)
	tokenPair, err := RotateRefreshToken(input.RefreshToken, c)
	if err != nil {
		return ErrorResponseJSON(c, fiber.StatusUnauthorized, "Invalid or expired refresh token", nil)
	}

	return SuccessResponse(c, "Token refreshed successfully", fiber.Map{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
		"token_type":    tokenPair.TokenType,
		"expires_in":    tokenPair.ExpiresIn,
	})
}

// Logout revokes the current refresh token
func Logout(c *fiber.Ctx) error {
	type LogoutInput struct{
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	input := new(LogoutInput)
	if err := c.BodyParser(input); err != nil {
		return ErrorResponseJSON(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Revoke the refresh token
	if err := RevokeRefreshToken(input.RefreshToken); err != nil {
		// Still return success even if token not found (idempotent)
		return SuccessResponse(c, "Logged out successfully", nil)
	}

	return SuccessResponse(c, "Logged out successfully", nil)
}

// LogoutAll revokes all refresh tokens for the current user (logout from all devices)
func LogoutAll(c *fiber.Ctx) error {
	// Get user ID from JWT token
	userID, err := GetUserIDFromToken(c)
	if err != nil {
		return ErrorResponseJSON(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	// Revoke all sessions
	if err := RevokeAllUserSessions(userID); err != nil {
		return ErrorResponseJSON(c, fiber.StatusInternalServerError, "Failed to logout from all devices", nil)
	}

	return SuccessResponse(c, "Logged out from all devices successfully", nil)
}

// GetActiveSessions returns all active sessions for the current user
func GetActiveSessions(c *fiber.Ctx) error {
	// Get user ID from JWT token
	userID, err := GetUserIDFromToken(c)
	if err != nil {
		return ErrorResponseJSON(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	// Get active sessions
	sessions, err := GetUserActiveSessions(userID)
	if err != nil {
		return ErrorResponseJSON(c, fiber.StatusInternalServerError, "Failed to fetch sessions", nil)
	}

	// Format sessions for response (hide refresh tokens)
	type SessionResponse struct {
		ID         uint   `json:"id"`
		UserAgent  string `json:"user_agent"`
		IPAddress  string `json:"ip_address"`
		LastUsedAt string `json:"last_used_at"`
		ExpiresAt  string `json:"expires_at"`
	}

	sessionResponses := make([]SessionResponse, len(sessions))
	for i, session := range sessions {
		sessionResponses[i] = SessionResponse{
			ID:         session.ID,
			UserAgent:  session.UserAgent,
			IPAddress:  session.IPAddress,
			LastUsedAt: session.LastUsedAt.Format("2006-01-02T15:04:05Z07:00"),
			ExpiresAt:  session.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return SuccessResponse(c, "Active sessions retrieved successfully", fiber.Map{
		"sessions": sessionResponses,
		"count":    len(sessionResponses),
	})
}
