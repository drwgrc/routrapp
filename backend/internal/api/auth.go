package api

import (
	"net/http"
	"time"

	"routrapp-api/internal/errors"
	"routrapp-api/internal/logger"
	"routrapp-api/internal/models"
	"routrapp-api/internal/utils/auth"
	"routrapp-api/internal/utils/constants"
	"routrapp-api/internal/validation"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	db         *gorm.DB
	jwtService *auth.JWTService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		db:         db,
		jwtService: auth.DefaultJWTService(),
	}
}

// Login handles POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req validation.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithContext(c).Errorf("Invalid login request: %v", err)
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

	logger.WithContext(c).Infof("Login attempt for email: %s", req.Email)

	// Find user by email
	var user models.User
	if err := h.db.Preload("Role").Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.WithContext(c).Warnf("Login failed: user not found for email %s", req.Email)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusUnauthorized,
					"Invalid credentials",
					map[string]interface{}{
						"code": "INVALID_CREDENTIALS",
					},
				),
			})
			return
		}
		logger.WithContext(c).Errorf("Database error during login: %v", err)
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
		logger.WithContext(c).Warnf("Login failed: user %s is inactive", req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusUnauthorized,
				"Account is disabled",
				map[string]interface{}{
					"code": "ACCOUNT_DISABLED",
				},
			),
		})
		return
	}

	// Verify password
	if err := auth.VerifyPassword(req.Password, user.Password); err != nil {
		logger.WithContext(c).Warnf("Login failed: invalid password for email %s", req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusUnauthorized,
				"Invalid credentials",
				map[string]interface{}{
					"code": "INVALID_CREDENTIALS",
				},
			),
		})
		return
	}

	// Check if user has a valid role
	if user.Role.Name == "" {
		logger.WithContext(c).Errorf("User %s has no associated role", req.Email)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusInternalServerError,
				"User role configuration error",
				map[string]interface{}{
					"code": "ROLE_CONFIGURATION_ERROR",
				},
			),
		})
		return
	}

	// Generate tokens
	accessToken, err := h.jwtService.GenerateAccessToken(user.ID, user.OrganizationID, user.Email, user.Role.Name.String())
	if err != nil {
		logger.WithContext(c).Errorf("Failed to generate access token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusInternalServerError,
				"Failed to generate access token",
				map[string]interface{}{
					"code": "TOKEN_GENERATION_ERROR",
				},
			),
		})
		return
	}

	refreshToken, err := h.jwtService.GenerateRefreshToken(user.ID, user.OrganizationID, user.Email, user.Role.Name.String())
	if err != nil {
		logger.WithContext(c).Errorf("Failed to generate refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusInternalServerError,
				"Failed to generate refresh token",
				map[string]interface{}{
					"code": "TOKEN_GENERATION_ERROR",
				},
			),
		})
		return
	}

	// Update user's refresh token and last login time
	now := time.Now()
	user.RefreshToken = refreshToken
	user.LastLoginAt = &now
	
	if err := h.db.Save(&user).Error; err != nil {
		logger.WithContext(c).Errorf("Failed to update user login info: %v", err)
		// Don't fail the login for this, just log the error
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

	loginResponse := validation.LoginResponse{
		User:         userResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    constants.JWT_ACCESS_TOKEN_EXPIRY,
	}

	logger.WithContext(c).Infof("User %s logged in successfully", req.Email)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    loginResponse,
		"message": "Login successful",
	})
}

// Logout handles POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		logger.WithContext(c).Error("User ID not found in context during logout")
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

	logger.WithContext(c).Infof("Logout request for user ID: %v", userID)

	// Clear refresh token from database
	if err := h.db.Model(&models.User{}).Where("id = ?", userID).Update("refresh_token", "").Error; err != nil {
		logger.WithContext(c).Errorf("Failed to clear refresh token during logout: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusInternalServerError,
				"Failed to logout",
				map[string]interface{}{
					"code": "LOGOUT_ERROR",
				},
			),
		})
		return
	}

	logger.WithContext(c).Infof("User %v logged out successfully", userID)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout successful",
	})
}

// RefreshToken handles POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req validation.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithContext(c).Errorf("Invalid refresh token request: %v", err)
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

	logger.WithContext(c).Info("Token refresh request received")

	// Validate refresh token
	claims, err := h.jwtService.ValidateToken(req.RefreshToken)
	if err != nil {
		logger.WithContext(c).Warnf("Invalid refresh token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusUnauthorized,
				"Invalid refresh token",
				map[string]interface{}{
					"code": "INVALID_REFRESH_TOKEN",
				},
			),
		})
		return
	}

	// Ensure this is a refresh token
	if !claims.IsRefreshToken() {
		logger.WithContext(c).Warn("Token is not a refresh token")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusUnauthorized,
				"Invalid token type",
				map[string]interface{}{
					"code": "INVALID_TOKEN_TYPE",
				},
			),
		})
		return
	}

	// Find user first to check if they are active
	var user models.User
	if err := h.db.Preload("Role").Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.WithContext(c).Warnf("User not found for refresh token: %d", claims.UserID)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusUnauthorized,
					"Invalid refresh token",
					map[string]interface{}{
						"code": "INVALID_REFRESH_TOKEN",
					},
				),
			})
			return
		}
		logger.WithContext(c).Errorf("Database error during token refresh: %v", err)
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

	// Check if user is still active
	if !user.Active {
		logger.WithContext(c).Warnf("Token refresh failed: user %d is inactive", user.ID)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusUnauthorized,
				"Account is disabled",
				map[string]interface{}{
					"code": "ACCOUNT_DISABLED",
				},
			),
		})
		return
	}

	// Verify refresh token matches stored one
	if user.RefreshToken != req.RefreshToken {
		logger.WithContext(c).Warnf("Refresh token mismatch for user %d", claims.UserID)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusUnauthorized,
				"Invalid refresh token",
				map[string]interface{}{
					"code": "INVALID_REFRESH_TOKEN",
				},
			),
		})
		return
	}

	// Check if user has a valid role
	if user.Role.Name == "" {
		logger.WithContext(c).Errorf("User %d has no associated role during token refresh", user.ID)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusInternalServerError,
				"User role configuration error",
				map[string]interface{}{
					"code": "ROLE_CONFIGURATION_ERROR",
				},
			),
		})
		return
	}

	// Generate new access token
	accessToken, err := h.jwtService.GenerateAccessToken(user.ID, user.OrganizationID, user.Email, user.Role.Name.String())
	if err != nil {
		logger.WithContext(c).Errorf("Failed to generate new access token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusInternalServerError,
				"Failed to generate access token",
				map[string]interface{}{
					"code": "TOKEN_GENERATION_ERROR",
				},
			),
		})
		return
	}

	// Prepare response
	tokenResponse := validation.TokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   constants.JWT_ACCESS_TOKEN_EXPIRY,
	}

	logger.WithContext(c).Infof("Token refreshed successfully for user %d", user.ID)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tokenResponse,
		"message": "Token refreshed successfully",
	})
} 