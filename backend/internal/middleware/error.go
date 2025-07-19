package middleware

import (
	"net/http"

	"routrapp-api/internal/errors"
	"routrapp-api/internal/logger"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a structured error response
type ErrorResponse struct {
	Error   string                 `json:"error"`
	Message string                 `json:"message"`
	Code    int                    `json:"code"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// ErrorHandlerMiddleware handles errors and returns consistent JSON responses
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			// Log the error with context
			logger.WithContext(c).WithField("error", err.Error()).Error("Request error")

			// Handle different error types
			switch e := err.Err.(type) {
			case *errors.AppError:
				// Handle custom application errors
				response := ErrorResponse{
					Error:   "application_error",
					Message: e.Message,
					Code:    e.Code,
					Details: e.Details,
				}
				c.JSON(e.Code, response)
			default:
				// Handle generic errors
				response := ErrorResponse{
					Error:   "internal_error",
					Message: "An internal error occurred",
					Code:    http.StatusInternalServerError,
				}
				c.JSON(http.StatusInternalServerError, response)
			}

			// Abort to prevent further processing
			c.Abort()
		}
	}
}

// HandleError is a helper function for handlers to handle errors
func HandleError(c *gin.Context, err error) {
	if err != nil {
		c.Error(err)
		return
	}
}

// HandleAppError is a helper function for handlers to handle AppError specifically
func HandleAppError(c *gin.Context, appErr *errors.AppError) {
	if appErr != nil {
		c.Error(appErr)
		return
	}
} 