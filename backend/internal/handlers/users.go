package handlers

import (
	"net/http"
	"strconv"

	"routrapp-api/internal/errors"
	"routrapp-api/internal/logger"
	"routrapp-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related requests
type UserHandler struct {
}

// NewUserHandler creates a new user handler
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// MockUser represents a mock user for testing
type MockUser struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Active    bool   `json:"active"`
}

// GetUser handles GET /api/v1/users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	
	// Check for empty ID (bad request scenario)
	if idParam == "" {
		middleware.HandleAppError(c, errors.BadRequest("User ID is required"))
		return
	}

	// Parse user ID
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		middleware.HandleAppError(c, errors.ValidationError("id", "must be a valid integer"))
		return
	}

	logger.WithContext(c).Infof("Fetching user with ID: %d", userID)

	// Handle different test scenarios based on ID
	switch {
	case userID == 123:
		// Success scenario
		user := MockUser{
			ID:     123,
			Name:   "John Doe",
			Email:  "john.doe@example.com",
			Active: true,
		}
		logger.WithContext(c).Info("Successfully retrieved user")
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    user,
		})
		
	case userID == 999:
		// Not found scenario
		middleware.HandleAppError(c, errors.NotFound("User not found"))
		return
		
	case userID == 500:
		// Internal server error scenario
		middleware.HandleAppError(c, errors.InternalServerError("Database connection failed"))
		return
		
	default:
		// Generic user response for other IDs
		user := MockUser{
			ID:     userID,
			Name:   "Test User",
			Email:  "test@example.com",
			Active: true,
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    user,
		})
	}
}

// GetUserWithEmptyID handles GET /api/v1/users/ (trailing slash)
func (h *UserHandler) GetUserWithEmptyID(c *gin.Context) {
	middleware.HandleAppError(c, errors.BadRequest("User ID cannot be empty"))
}

// TriggerPanic handles GET /api/v1/panic
func (h *UserHandler) TriggerPanic(c *gin.Context) {
	logger.WithContext(c).Warn("Panic endpoint triggered - this is intentional for testing")
	panic("This is a test panic to demonstrate recovery middleware")
}

// GetUsers handles GET /api/v1/users
func (h *UserHandler) GetUsers(c *gin.Context) {
	logger.WithContext(c).Info("Fetching all users")
	
	users := []MockUser{
		{ID: 1, Name: "Alice Smith", Email: "alice@example.com", Active: true},
		{ID: 2, Name: "Bob Johnson", Email: "bob@example.com", Active: true},
		{ID: 3, Name: "Charlie Brown", Email: "charlie@example.com", Active: false},
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    users,
		"count":   len(users),
	})
}

// CreateUser handles POST /api/v1/users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleAppError(c, errors.BadRequest("Invalid request data: " + err.Error()))
		return
	}
	
	logger.WithContext(c).Infof("Creating user: %s (%s)", req.Name, req.Email)
	
	// Simulate creating a user
	user := MockUser{
		ID:     999, // Mock ID
		Name:   req.Name,
		Email:  req.Email,
		Active: true,
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    user,
		"message": "User created successfully",
	})
} 