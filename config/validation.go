package config

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// RequiredEnvVars lists all required environment variables for the application
var RequiredEnvVars = []string{
	"DB_HOST",
	"DB_PORT",
	"DB_USER",
	"DB_PASSWORD",
	"DB_NAME",
	"SECRET",
	"ACCESS_TOKEN_SECRET",
	"REFRESH_TOKEN_SECRET",
}

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Message string
}

// ValidateConfig validates all required environment variables are set and valid
func ValidateConfig() []ValidationError {
	var errors []ValidationError

	// Check all required variables exist
	for _, key := range RequiredEnvVars {
		value := Config(key)
		if value == "" {
			errors = append(errors, ValidationError{
				Field:   key,
				Message: fmt.Sprintf("Required environment variable %s is not set", key),
			})
		}
	}

	// If basic validation failed, return early
	if len(errors) > 0 {
		return errors
	}

	// Validate DB_PORT is a valid number
	if _, err := strconv.Atoi(Config("DB_PORT")); err != nil {
		errors = append(errors, ValidationError{
			Field:   "DB_PORT",
			Message: "DB_PORT must be a valid number",
		})
	}

	// Validate JWT secrets
	if err := ValidateJWTSecret(Config("SECRET"), "SECRET"); err.Field != "" {
		errors = append(errors, err)
	}
	if err := ValidateJWTSecret(Config("ACCESS_TOKEN_SECRET"), "ACCESS_TOKEN_SECRET"); err.Field != "" {
		errors = append(errors, err)
	}
	if err := ValidateJWTSecret(Config("REFRESH_TOKEN_SECRET"), "REFRESH_TOKEN_SECRET"); err.Field != "" {
		errors = append(errors, err)
	}

	// Ensure secrets are different (security best practice)
	accessSecret := Config("ACCESS_TOKEN_SECRET")
	refreshSecret := Config("REFRESH_TOKEN_SECRET")
	if accessSecret == refreshSecret {
		errors = append(errors, ValidationError{
			Field:   "ACCESS_TOKEN_SECRET/REFRESH_TOKEN_SECRET",
			Message: "Access and refresh token secrets must be different for security",
		})
	}

	return errors
}

// ValidateJWTSecret validates a JWT secret meets security requirements
func ValidateJWTSecret(secret, fieldName string) ValidationError {
	// Minimum length check
	const minLength = 32
	if len(secret) < minLength {
		return ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("%s must be at least %d characters long (current: %d)", fieldName, minLength, len(secret)),
		}
	}

	// Check it's not a weak/default value
	weakSecrets := []string{
		"secret",
		"password",
		"changeme",
		"example",
		"test",
		"default",
		"admin",
		"12345",
	}

	secretLower := strings.ToLower(secret)
	for _, weak := range weakSecrets {
		if strings.Contains(secretLower, weak) {
			return ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("%s appears to contain weak/default values. Please use a cryptographically random secret", fieldName),
			}
		}
	}

	// Check for sufficient complexity (at least some variation)
	if !hasMinimumComplexity(secret) {
		return ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("%s lacks sufficient complexity. Generate using: openssl rand -base64 32", fieldName),
		}
	}

	return ValidationError{}
}

// hasMinimumComplexity checks if string has minimum character variation
func hasMinimumComplexity(s string) bool {
	// Check for at least 2 different types of characters
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(s)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(s)
	hasSpecial := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(s)

	complexityCount := 0
	if hasLetter {
		complexityCount++
	}
	if hasNumber {
		complexityCount++
	}
	if hasSpecial {
		complexityCount++
	}

	return complexityCount >= 2
}

// GenerateSecureSecret generates a cryptographically secure random secret
// This is a helper function for developers to generate secrets
func GenerateSecureSecret(length int) (string, error) {
	if length < 32 {
		length = 32
	}

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate secure secret: %w", err)
	}

	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// PrintValidationErrors prints validation errors in a readable format
func PrintValidationErrors(errors []ValidationError) {
	fmt.Println("\n❌ Configuration Validation Failed:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	for _, err := range errors {
		fmt.Printf("  • %s: %s\n", err.Field, err.Message)
	}
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("\n💡 To generate secure secrets, run:")
	fmt.Println("   openssl rand -base64 32")
	fmt.Println()
}
