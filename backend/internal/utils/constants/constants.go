package constants

import (
	"os"
	"time"
)

// Application constants
const (
	// API version
	API_VERSION = "v1"
	
	// Context keys
	TENANT_CONTEXT_KEY = "tenant_context"
	USER_CONTEXT_KEY   = "user_context"
	
	// JWT settings - default values
	DEFAULT_JWT_SECRET                = "dev-secret-key-change-in-production"
	JWT_ACCESS_TOKEN_EXPIRY          = 15 * 60                // 15 minutes in seconds
	JWT_REFRESH_TOKEN_EXPIRY         = 7 * 24 * 60 * 60       // 7 days in seconds
)

// JWT_SECRET returns the JWT secret from environment or default value
func JWT_SECRET() string {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		return secret
	}
	return DEFAULT_JWT_SECRET
}

const (
	// Server defaults
	DefaultPort         = "8080"
	DefaultReadTimeout  = 10 * time.Second
	DefaultWriteTimeout = 10 * time.Second
	
	// CORS defaults
	DefaultFrontendURL = "http://localhost:3000"

	// Database defaults
	DefaultDBHost        = "localhost"
	DefaultDBPort        = "5432"
	DefaultDBUser        = "postgres"
	DefaultDBPassword    = "postgres"
	DefaultDBName        = "routrapp"
	DefaultDBSSLMode     = "disable"
	DefaultDBMaxIdleConns = 10
	DefaultDBMaxOpenConns = 100
	DefaultDBConnMaxLife  = 30 // in seconds
) 