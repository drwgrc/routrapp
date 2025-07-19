package errors

import (
	"fmt"
	"net/http"
)

// AppError represents a structured application error
type AppError struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a new AppError
func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// NewAppErrorWithDetails creates a new AppError with additional details
func NewAppErrorWithDetails(code int, message string, details map[string]interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Common error constructors
func BadRequest(message string) *AppError {
	return NewAppError(http.StatusBadRequest, message)
}

func BadRequestWithDetails(message string, details map[string]interface{}) *AppError {
	return NewAppErrorWithDetails(http.StatusBadRequest, message, details)
}

func Unauthorized(message string) *AppError {
	return NewAppError(http.StatusUnauthorized, message)
}

func Forbidden(message string) *AppError {
	return NewAppError(http.StatusForbidden, message)
}

func NotFound(message string) *AppError {
	return NewAppError(http.StatusNotFound, message)
}

func Conflict(message string) *AppError {
	return NewAppError(http.StatusConflict, message)
}

func InternalServerError(message string) *AppError {
	return NewAppError(http.StatusInternalServerError, message)
}

func ServiceUnavailable(message string) *AppError {
	return NewAppError(http.StatusServiceUnavailable, message)
}

// ValidationError creates a bad request error for validation failures
func ValidationError(field string, message string) *AppError {
	return NewAppErrorWithDetails(
		http.StatusBadRequest,
		"Validation failed",
		map[string]interface{}{
			"field": field,
			"error": message,
		},
	)
}

// DatabaseError creates an internal server error for database issues
func DatabaseError(operation string, err error) *AppError {
	return NewAppErrorWithDetails(
		http.StatusInternalServerError,
		"Database operation failed",
		map[string]interface{}{
			"operation": operation,
			"error":     err.Error(),
		},
	)
}

// ExternalServiceError creates a service unavailable error for external service failures
func ExternalServiceError(service string, err error) *AppError {
	return NewAppErrorWithDetails(
		http.StatusServiceUnavailable,
		fmt.Sprintf("External service %s is unavailable", service),
		map[string]interface{}{
			"service": service,
			"error":   err.Error(),
		},
	)
} 