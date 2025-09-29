package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// GenerateSecureToken generates a cryptographically secure random token
func GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate secure token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateVerificationToken generates a secure verification token (64 chars)
func GenerateVerificationToken() (string, error) {
	return GenerateSecureToken(32) // 32 bytes = 64 hex characters
}

// ValidateTokenFormat validates if token has correct format
func ValidateTokenFormat(token string) bool {
	if len(token) != 64 { // 32 bytes = 64 hex chars
		return false
	}

	// Check if all characters are valid hex
	for _, char := range token {
		if !((char >= '0' && char <= '9') ||
			(char >= 'a' && char <= 'f') ||
			(char >= 'A' && char <= 'F')) {
			return false
		}
	}
	return true
}