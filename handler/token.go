package handler

import (
	"app/config"
	"app/database"
	"app/model"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gofiber/fiber/v2"
)

const (
	// AccessTokenDuration is the lifespan of an access token (15 minutes)
	// Short-lived for security - requires refresh token for extended sessions
	AccessTokenDuration = 15 * time.Minute

	// RefreshTokenDuration is the lifespan of a refresh token (7 days)
	// Can be extended up to 30 days based on security requirements
	RefreshTokenDuration = 7 * 24 * time.Hour

	// RefreshTokenLength is the length of the random refresh token string
	RefreshTokenLength = 64
)

// TokenPair represents both access and refresh tokens
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int64     `json:"expires_in"`     // Access token expiration in seconds
	ExpiresAt    time.Time `json:"-"`              // Internal use only
	TokenType    string    `json:"token_type"`     // Always "Bearer"
}

// GenerateTokenPair creates both access and refresh tokens for a user
func GenerateTokenPair(user *model.User, c *fiber.Ctx) (*TokenPair, error) {
	// Generate JWT access token
	accessToken, expiresAt, err := GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token (random string)
	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store refresh token in database
	session := &model.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    c.Get("User-Agent"),
		IPAddress:    c.IP(),
		ExpiresAt:    time.Now().Add(RefreshTokenDuration),
		LastUsedAt:   time.Now(),
		IsRevoked:    false,
	}

	db := database.DB
	if err := db.Create(session).Error; err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(AccessTokenDuration.Seconds()),
		ExpiresAt:    expiresAt,
		TokenType:    "Bearer",
	}, nil
}

// GenerateAccessToken creates a JWT access token for a user
func GenerateAccessToken(user *model.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(AccessTokenDuration)

	// Create token claims
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      expiresAt.Unix(),
		"iat":      time.Now().Unix(),
		"type":     "access",
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	secret := config.Config("ACCESS_TOKEN_SECRET")
	if secret == "" {
		secret = config.Config("SECRET") // Fallback for backward compatibility
	}

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expiresAt, nil
}

// GenerateRefreshToken creates a cryptographically secure random refresh token
func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, RefreshTokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// ValidateRefreshToken validates a refresh token and returns the associated session
func ValidateRefreshToken(refreshToken string) (*model.Session, error) {
	db := database.DB
	var session model.Session

	// Find session by refresh token
	if err := db.Where("refresh_token = ?", refreshToken).First(&session).Error; err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Check if session is valid
	if !session.IsValid() {
		return nil, fmt.Errorf("refresh token expired or revoked")
	}

	return &session, nil
}

// RevokeRefreshToken revokes a specific refresh token
func RevokeRefreshToken(refreshToken string) error {
	db := database.DB
	var session model.Session

	if err := db.Where("refresh_token = ?", refreshToken).First(&session).Error; err != nil {
		return fmt.Errorf("refresh token not found")
	}

	return session.Revoke(db)
}

// RevokeAllUserSessions revokes all sessions for a specific user (logout all devices)
func RevokeAllUserSessions(userID uint) error {
	db := database.DB
	return db.Model(&model.Session{}).
		Where("user_id = ? AND is_revoked = ?", userID, false).
		Update("is_revoked", true).Error
}

// CleanupExpiredSessions removes expired sessions from the database
// This should be called periodically (e.g., daily cron job)
func CleanupExpiredSessions() error {
	db := database.DB
	return db.Where("expires_at < ?", time.Now()).Delete(&model.Session{}).Error
}

// GetUserActiveSessions retrieves all active sessions for a user
func GetUserActiveSessions(userID uint) ([]model.Session, error) {
	db := database.DB
	var sessions []model.Session

	err := db.Where("user_id = ? AND is_revoked = ? AND expires_at > ?",
		userID, false, time.Now()).
		Order("last_used_at DESC").
		Find(&sessions).Error

	return sessions, err
}

// RotateRefreshToken creates a new refresh token and revokes the old one
// This implements token rotation for enhanced security
func RotateRefreshToken(oldRefreshToken string, c *fiber.Ctx) (*TokenPair, error) {
	// Validate old refresh token
	session, err := ValidateRefreshToken(oldRefreshToken)
	if err != nil {
		return nil, err
	}

	// Get user
	db := database.DB
	var user model.User
	if err := db.First(&user, session.UserID).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Revoke old token
	if err := session.Revoke(db); err != nil {
		return nil, fmt.Errorf("failed to revoke old token: %w", err)
	}

	// Generate new token pair
	tokenPair, err := GenerateTokenPair(&user, c)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new tokens: %w", err)
	}

	return tokenPair, nil
}
