package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	// DefaultCost is the default bcrypt cost
	DefaultCost = 12
)

// HashPassword hashes a password using bcrypt with the default cost
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// VerifyPassword compares a password with its hash
func VerifyPassword(password, hashedPassword string) error {
	if password == "" {
		return errors.New("password cannot be empty")
	}
	if hashedPassword == "" {
		return errors.New("hashed password cannot be empty")
	}

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// IsValidPassword checks if a password meets the minimum requirements
func IsValidPassword(password string) bool {
	// Minimum 8 characters - validation is also done at the request level
	return len(password) >= 8
} 