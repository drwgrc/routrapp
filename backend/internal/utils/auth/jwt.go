package auth

import (
	"errors"
	"fmt"
	"time"

	"routrapp-api/internal/utils/constants"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the JWT claims structure
type JWTClaims struct {
	UserID         uint   `json:"user_id"`
	OrganizationID uint   `json:"organization_id"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	TokenType      string `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// JWTService handles JWT token operations
type JWTService struct {
	secretKey []byte
}

// NewJWTService creates a new JWT service instance
func NewJWTService(secretKey string) *JWTService {
	return &JWTService{
		secretKey: []byte(secretKey),
	}
}

// GenerateAccessToken generates a new access token for the user
func (j *JWTService) GenerateAccessToken(userID, organizationID uint, email, role string) (string, error) {
	claims := JWTClaims{
		UserID:         userID,
		OrganizationID: organizationID,
		Email:          email,
		Role:           role,
		TokenType:      "access",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(constants.JWT_ACCESS_TOKEN_EXPIRY) * time.Second)),
			Subject:   fmt.Sprintf("%d", userID),
			Issuer:    "routrapp-api",
			Audience:  []string{"routrapp-frontend"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// GenerateRefreshToken generates a new refresh token for the user
func (j *JWTService) GenerateRefreshToken(userID, organizationID uint, email, role string) (string, error) {
	claims := JWTClaims{
		UserID:         userID,
		OrganizationID: organizationID,
		Email:          email,
		Role:           role,
		TokenType:      "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(constants.JWT_REFRESH_TOKEN_EXPIRY) * time.Second)),
			Subject:   fmt.Sprintf("%d", userID),
			Issuer:    "routrapp-api",
			Audience:  []string{"routrapp-frontend"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// ValidateToken validates and parses a JWT token
func (j *JWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Verify token is not expired
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

// ExtractTokenFromHeader extracts the JWT token from the Authorization header
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is empty")
	}

	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", errors.New("authorization header must start with 'Bearer '")
	}

	token := authHeader[7:]
	if token == "" {
		return "", errors.New("token is empty")
	}

	return token, nil
}

// IsAccessToken checks if the token is an access token
func (c *JWTClaims) IsAccessToken() bool {
	return c.TokenType == "access"
}

// IsRefreshToken checks if the token is a refresh token
func (c *JWTClaims) IsRefreshToken() bool {
	return c.TokenType == "refresh"
}

// GetUserContext returns a simplified user context from JWT claims
func (c *JWTClaims) GetUserContext() map[string]interface{} {
	return map[string]interface{}{
		"user_id":         c.UserID,
		"organization_id": c.OrganizationID,
		"email":           c.Email,
		"role":            c.Role,
	}
}

// DefaultJWTService returns a JWT service instance with the default secret key
func DefaultJWTService() *JWTService {
	return NewJWTService(constants.JWT_SECRET())
} 