package handler

import (
	"net/mail"
	"strings"
)

// ValidateEmail checks if an email address is valid
func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// NormalizeEmail normalizes email address (lowercase and trim)
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// NormalizeUsername normalizes username (lowercase and trim)
func NormalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

// SanitizeString removes leading/trailing whitespace and limits length
func SanitizeString(input string, maxLength int) string {
	trimmed := strings.TrimSpace(input)
	if len(trimmed) > maxLength {
		return trimmed[:maxLength]
	}
	return trimmed
}
