package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"routrapp-api/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// RecoveryMiddleware handles panics and recovers gracefully
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// Get stack trace
		stack := debug.Stack()
		
		// Log the panic with full context and stack trace
		logger.WithContext(c).WithFields(logrus.Fields{
			"panic":      fmt.Sprintf("%v", recovered),
			"stack":      string(stack),
			"request_id": c.GetString("X-Request-ID"),
		}).Error("Panic recovered")

		// Return a clean error response
		errorResponse := ErrorResponse{
			Error:   "internal_error",
			Message: "An unexpected error occurred",
			Code:    http.StatusInternalServerError,
		}

		c.JSON(http.StatusInternalServerError, errorResponse)
		c.Abort()
	})
}

// CustomRecoveryMiddleware allows for custom recovery handling
func CustomRecoveryMiddleware(handler gin.RecoveryFunc) gin.HandlerFunc {
	return gin.CustomRecovery(handler)
} 