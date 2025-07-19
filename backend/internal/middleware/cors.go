package middleware

import (
	"routrapp-api/internal/config"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware creates a CORS middleware with the configured frontend URL
func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set CORS headers
		c.Header("Access-Control-Allow-Origin", cfg.CORS.FrontendURL)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
} 