package handler

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	// MinPasswordLength is the minimum allowed password length
	MinPasswordLength = 8
	// MaxPasswordLength is the maximum allowed password length
	MaxPasswordLength = 100
	// BcryptCost is the cost factor for bcrypt hashing
	BcryptCost = 14
)

var (
	// ErrWeakPassword is returned when password doesn't meet requirements
	ErrWeakPassword = errors.New("password must be at least 8 characters long")
	// ErrPasswordHashFailed is returned when password hashing fails
	ErrPasswordHashFailed = errors.New("failed to hash password")
	// ErrInvalidPassword is returned when password verification fails
	ErrInvalidPassword = errors.New("invalid password")
)

// HashPassword generates a bcrypt hash from a password
func HashPassword(password string) (string, error) {
	if len(password) < MinPasswordLength {
		return "", ErrWeakPassword
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", ErrPasswordHashFailed
	}

	return string(bytes), nil
}

// CheckPasswordHash compares a password with its hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidatePasswordStrength checks if password meets security requirements
// TODO: Implement more sophisticated password strength validation
// - Check for uppercase, lowercase, numbers, special characters
// - Check against common password lists
// - Implement password entropy checking
func ValidatePasswordStrength(password string) error {
	if len(password) < MinPasswordLength {
		return ErrWeakPassword
	}
	if len(password) > MaxPasswordLength {
		return errors.New("password is too long")
	}

	// TODO: Add more checks here
	// hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	// hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	// hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	// hasSpecial := regexp.MustCompile(`[!@#$%^&*]`).MatchString(password)

	return nil
}
