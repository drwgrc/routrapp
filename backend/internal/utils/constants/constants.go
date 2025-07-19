package constants

import "time"

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