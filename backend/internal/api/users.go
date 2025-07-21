package api

import (
	"net/http"
	"strconv"

	"routrapp-api/internal/errors"
	"routrapp-api/internal/logger"
	"routrapp-api/internal/middleware"
	"routrapp-api/internal/models"
	"routrapp-api/internal/validation"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserHandler handles user-related requests
type UserHandler struct {
	db *gorm.DB
}

// NewUserHandler creates a new user handler
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		db: db,
	}
}

// UpdateProfile handles PUT /api/v1/users/profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req validation.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithContext(c).Errorf("Invalid user update request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusBadRequest,
				"Invalid request data: "+err.Error(),
				map[string]interface{}{
					"code": "VALIDATION_ERROR",
				},
			),
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userIDValue, exists := c.Get("user_id")
	if !exists {
		logger.WithContext(c).Error("User ID not found in context during profile update")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusUnauthorized,
				"Authentication required",
				map[string]interface{}{
					"code": "AUTHENTICATION_REQUIRED",
				},
			),
		})
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		logger.WithContext(c).Error("Invalid user ID type in context during profile update")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusInternalServerError,
				"Internal server error",
				map[string]interface{}{
					"code": "INTERNAL_ERROR",
				},
			),
		})
		return
	}

	logger.WithContext(c).Infof("Profile update request for user ID: %d", userID)

	// Check if at least one field is provided for update
	if req.FirstName == nil && req.LastName == nil {
		logger.WithContext(c).Warnf("Profile update failed: no fields provided for user %d", userID)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusBadRequest,
				"At least one field must be provided for update",
				map[string]interface{}{
					"code": "NO_FIELDS_PROVIDED",
				},
			),
		})
		return
	}

	// Find user
	var user models.User
	if err := h.db.Preload("Role").Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.WithContext(c).Warnf("Profile update failed: user not found %d", userID)
			c.JSON(http.StatusNotFound, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusNotFound,
					"User not found",
					map[string]interface{}{
						"code": "USER_NOT_FOUND",
					},
				),
			})
			return
		}
		logger.WithContext(c).Errorf("Database error during profile update: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusInternalServerError,
				"Internal server error",
				map[string]interface{}{
					"code": "INTERNAL_ERROR",
				},
			),
		})
		return
	}

	// Check if user is active
	if !user.Active {
		logger.WithContext(c).Warnf("Profile update failed: user %d is inactive", userID)
		c.JSON(http.StatusForbidden, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusForbidden,
				"Account is disabled",
				map[string]interface{}{
					"code": "ACCOUNT_DISABLED",
				},
			),
		})
		return
	}

	// Prepare update data
	updateData := make(map[string]interface{})
	if req.FirstName != nil {
		updateData["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		updateData["last_name"] = *req.LastName
	}

	// Update user in database
	if err := h.db.Model(&user).Updates(updateData).Error; err != nil {
		logger.WithContext(c).Errorf("Failed to update user profile: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusInternalServerError,
				"Failed to update profile",
				map[string]interface{}{
					"code": "PROFILE_UPDATE_ERROR",
				},
			),
		})
		return
	}

	// Reload user to get updated data
	if err := h.db.Preload("Role").First(&user, user.ID).Error; err != nil {
		logger.WithContext(c).Errorf("Failed to reload user after profile update: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusInternalServerError,
				"Profile updated but failed to reload data",
				map[string]interface{}{
					"code": "USER_RELOAD_ERROR",
				},
			),
		})
		return
	}

	// Prepare response
	userResponse := validation.UserResponse{
		BaseResponse: validation.BaseResponse{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Active:    user.Active,
		Role:      user.Role.Name.String(),
	}

	logger.WithContext(c).Infof("Profile updated successfully for user %d", userID)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    userResponse,
		"message": "Profile updated successfully",
	})
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