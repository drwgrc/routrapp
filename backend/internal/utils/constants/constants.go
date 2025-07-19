package constants

import "time"

const (
	// Server defaults
	DefaultPort         = "8080"
	DefaultReadTimeout  = 10 * time.Second
	DefaultWriteTimeout = 10 * time.Second
	
	// CORS defaults
	DefaultFrontendURL = "http://localhost:3000"
) 