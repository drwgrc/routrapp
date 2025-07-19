package middleware

import (
	"time"

	"routrapp-api/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoggerMiddleware logs HTTP requests with timing and response information
func LoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			// Use our structured logger instead of default gin logging
			logEntry := logger.GetLogger().WithFields(logrus.Fields{
				"method":      param.Method,
				"path":        param.Path,
				"status_code": param.StatusCode,
				"latency":     param.Latency.String(),
				"client_ip":   param.ClientIP,
				"user_agent":  param.Request.UserAgent(),
				"time":        param.TimeStamp.Format(time.RFC3339),
			})

			// Add error information if present
			if param.ErrorMessage != "" {
				logEntry = logEntry.WithField("error", param.ErrorMessage)
			}

			// Log at different levels based on status code
			switch {
			case param.StatusCode >= 500:
				logEntry.Error("HTTP Request")
			case param.StatusCode >= 400:
				logEntry.Warn("HTTP Request")
			default:
				logEntry.Info("HTTP Request")
			}

			// Return empty string since we're handling logging ourselves
			return ""
		},
		Output: logger.GetLogger().Out,
	})
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			// Generate a simple request ID (in production, consider using UUID)
			requestID = generateRequestID()
		}
		
		c.Header("X-Request-ID", requestID)
		c.Set("X-Request-ID", requestID)
		
		c.Next()
	}
}

// generateRequestID generates a simple request ID
// In production, consider using github.com/google/uuid for proper UUIDs
func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + generateRandomString(6)
}

// generateRandomString generates a random string of given length
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
} 