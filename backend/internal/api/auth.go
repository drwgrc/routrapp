package api

import (
	"net/http"
	"strings"
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

// Register handles POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req validation.UserRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithContext(c).Errorf("Invalid registration request: %v", err)
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

	logger.WithContext(c).Infof("Registration attempt for email: %s", req.Email)

	// Enhanced password validation
	if err := auth.ValidatePassword(req.Password); err != nil {
		logger.WithContext(c).Warnf("Registration failed: weak password for email %s", req.Email)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusBadRequest,
				"Password does not meet security requirements: "+err.Error(),
				map[string]interface{}{
					"code": "WEAK_PASSWORD",
				},
			),
		})
		return
	}

	// Check for common passwords
	if auth.IsCommonPassword(req.Password) {
		logger.WithContext(c).Warnf("Registration failed: common password used for email %s", req.Email)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusBadRequest,
				"Password is too common, please choose a more secure password",
				map[string]interface{}{
					"code": "COMMON_PASSWORD",
				},
			),
		})
		return
	}

	// Verify organization exists
	var organization models.Organization
	if err := h.db.Where("id = ? AND active = ?", req.TenantID, true).First(&organization).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.WithContext(c).Warnf("Registration failed: organization not found %d", req.TenantID)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusBadRequest,
					"Invalid organization",
					map[string]interface{}{
						"code": "INVALID_ORGANIZATION",
					},
				),
			})
			return
		}
		logger.WithContext(c).Errorf("Database error during registration: %v", err)
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

	// Check if user already exists
	var existingUser models.User
	err := h.db.Where("organization_id = ? AND email = ?", req.TenantID, req.Email).First(&existingUser).Error
	if err == nil {
		logger.WithContext(c).Warnf("Registration failed: email already exists %s", req.Email)
		c.JSON(http.StatusConflict, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusConflict,
				"Email address is already registered",
				map[string]interface{}{
					"code": "EMAIL_EXISTS",
				},
			),
		})
		return
	} else if err != gorm.ErrRecordNotFound {
		logger.WithContext(c).Errorf("Database error checking user existence: %v", err)
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

	// Find the appropriate role
	var role models.Role
	if err := h.db.Where("organization_id = ? AND name = ? AND active = ?", req.TenantID, req.Role.String(), true).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.WithContext(c).Warnf("Registration failed: role not found %s for organization %d", req.Role.String(), req.TenantID)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusBadRequest,
					"Invalid role for organization",
					map[string]interface{}{
						"code": "INVALID_ROLE",
					},
				),
			})
			return
		}
		logger.WithContext(c).Errorf("Database error finding role: %v", err)
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

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		logger.WithContext(c).Errorf("Failed to hash password during registration: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusInternalServerError,
				"Failed to process password",
				map[string]interface{}{
					"code": "PASSWORD_PROCESSING_ERROR",
				},
			),
		})
		return
	}

	// Create user
	user := models.User{
		Base: models.Base{
			OrganizationID: req.TenantID,
		},
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		RoleID:    role.ID,
		Active:    true,
	}

	if err := h.db.Create(&user).Error; err != nil {
		logger.WithContext(c).Errorf("Failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusInternalServerError,
				"Failed to create user account",
				map[string]interface{}{
					"code": "USER_CREATION_ERROR",
				},
			),
		})
		return
	}

	// Load the role for response
	if err := h.db.Preload("Role").First(&user, user.ID).Error; err != nil {
		logger.WithContext(c).Errorf("Failed to load user role after creation: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusInternalServerError,
				"User created but failed to load details",
				map[string]interface{}{
					"code": "USER_LOAD_ERROR",
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

	logger.WithContext(c).Infof("User %s registered successfully", req.Email)
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    userResponse,
		"message": "User registered successfully",
	})
}

// ChangePassword handles POST /api/v1/auth/change-password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req validation.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithContext(c).Errorf("Invalid change password request: %v", err)
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

	// Get user from context (set by auth middleware)
	userIDValue, exists := c.Get("user_id")
	if !exists {
		logger.WithContext(c).Error("User ID not found in context during password change")
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
		logger.WithContext(c).Error("Invalid user ID type in context during password change")
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

	logger.WithContext(c).Infof("Password change request for user ID: %d", userID)

	// Enhanced password validation for new password
	if err := auth.ValidatePassword(req.NewPassword); err != nil {
		logger.WithContext(c).Warnf("Password change failed: weak password for user %d", userID)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusBadRequest,
				"New password does not meet security requirements: "+err.Error(),
				map[string]interface{}{
					"code": "WEAK_PASSWORD",
				},
			),
		})
		return
	}

	// Check for common passwords
	if auth.IsCommonPassword(req.NewPassword) {
		logger.WithContext(c).Warnf("Password change failed: common password used for user %d", userID)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusBadRequest,
				"New password is too common, please choose a more secure password",
				map[string]interface{}{
					"code": "COMMON_PASSWORD",
				},
			),
		})
		return
	}

	// Check if new password is the same as current password
	if req.CurrentPassword == req.NewPassword {
		logger.WithContext(c).Warnf("Password change failed: same password provided for user %d", userID)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusBadRequest,
				"New password must be different from current password",
				map[string]interface{}{
					"code": "SAME_PASSWORD",
				},
			),
		})
		return
	}

	// Find user and verify current password
	var user models.User
	if err := h.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.WithContext(c).Warnf("Password change failed: user not found %d", userID)
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
		logger.WithContext(c).Errorf("Database error during password change: %v", err)
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
		logger.WithContext(c).Warnf("Password change failed: user %d is inactive", userID)
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

	// Verify current password
	if err := auth.VerifyPassword(req.CurrentPassword, user.Password); err != nil {
		logger.WithContext(c).Warnf("Password change failed: invalid current password for user %d", userID)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusUnauthorized,
				"Current password is incorrect",
				map[string]interface{}{
					"code": "INVALID_CURRENT_PASSWORD",
				},
			),
		})
		return
	}

	// Hash new password
	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		logger.WithContext(c).Errorf("Failed to hash new password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusInternalServerError,
				"Failed to process new password",
				map[string]interface{}{
					"code": "PASSWORD_PROCESSING_ERROR",
				},
			),
		})
		return
	}

	// Update password in database and clear refresh tokens for security
	if err := h.db.Model(&user).Updates(map[string]interface{}{
		"password":      hashedPassword,
		"refresh_token": "", // Force re-login for security
		"updated_at":    time.Now(),
	}).Error; err != nil {
		logger.WithContext(c).Errorf("Failed to update password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusInternalServerError,
				"Failed to update password",
				map[string]interface{}{
					"code": "PASSWORD_UPDATE_ERROR",
				},
			),
		})
		return
	}

	logger.WithContext(c).Infof("Password changed successfully for user %d", userID)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password changed successfully. Please log in again.",
	})
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

// Register handles POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req validation.RegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithContext(c).Errorf("Invalid registration request: %v", err)
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

	logger.WithContext(c).Infof("Registration attempt for email: %s, organization: %s", req.Email, req.OrganizationName)

	// Check if subdomain is already taken
	var existingOrg models.Organization
	if err := h.db.Where("sub_domain = ? AND deleted_at IS NULL", req.SubDomain).First(&existingOrg).Error; err == nil {
		logger.WithContext(c).Warnf("Registration failed: subdomain %s already exists", req.SubDomain)
		c.JSON(http.StatusConflict, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusConflict,
				"Subdomain is already taken",
				map[string]interface{}{
					"code": "SUBDOMAIN_TAKEN",
				},
			),
		})
		return
	} else if err != gorm.ErrRecordNotFound {
		logger.WithContext(c).Errorf("Database error checking subdomain: %v", err)
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

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		logger.WithContext(c).Errorf("Failed to hash password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusInternalServerError,
				"Failed to process registration",
				map[string]interface{}{
					"code": "PASSWORD_HASH_ERROR",
				},
			),
		})
		return
	}

	// Begin transaction for creating organization and user atomically
	tx := h.db.Begin()
	if tx.Error != nil {
		logger.WithContext(c).Errorf("Failed to begin transaction: %v", tx.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusInternalServerError,
				"Internal server error",
				map[string]interface{}{
					"code": "TRANSACTION_ERROR",
				},
			),
		})
		return
	}

	// Flag to track if panic recovery has been triggered
	var panicRecovered bool
	defer func() {
		if r := recover(); r != nil {
			panicRecovered = true
			tx.Rollback()
			logger.WithContext(c).Errorf("Panic in registration transaction: %v", r)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusInternalServerError,
					"Registration failed due to internal error",
					map[string]interface{}{
						"code": "INTERNAL_PANIC_ERROR",
					},
				),
			})
		}
	}()

	// Create organization
	org := models.Organization{
		Name:         req.OrganizationName,
		SubDomain:    req.SubDomain,
		ContactEmail: req.OrganizationEmail,
		Active:       true,
		PlanType:     "basic",
	}

	if err := tx.Create(&org).Error; err != nil {
		tx.Rollback()
		logger.WithContext(c).Errorf("Failed to create organization: %v", err)
		
		// Check if it's a duplicate subdomain error
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			c.JSON(http.StatusConflict, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusConflict,
					"Subdomain is already taken",
					map[string]interface{}{
						"code": "SUBDOMAIN_TAKEN",
					},
				),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusInternalServerError,
					"Failed to create organization",
					map[string]interface{}{
						"code": "ORGANIZATION_CREATION_ERROR",
					},
				),
			})
		}
		return
	}

	// Create default roles for the organization
	ownerRole := models.Role{
		Base: models.Base{
			OrganizationID: org.ID,
		},
		Name:        models.RoleTypeOwner,
		DisplayName: "Organization Owner",
		Description: "Full access to all organization resources and settings",
		Permissions: `["organizations.*", "users.*", "technicians.*", "routes.*", "roles.*"]`,
		Active:      true,
	}

	if err := tx.Create(&ownerRole).Error; err != nil {
		tx.Rollback()
		logger.WithContext(c).Errorf("Failed to create owner role: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusInternalServerError,
				"Failed to create organization roles",
				map[string]interface{}{
					"code": "ROLE_CREATION_ERROR",
				},
			),
		})
		return
	}

	technicianRole := models.Role{
		Base: models.Base{
			OrganizationID: org.ID,
		},
		Name:        models.RoleTypeTechnician,
		DisplayName: "Technician",
		Description: "Access to routes and personal information",
		Permissions: `["routes.read", "routes.update_status", "technicians.read_own", "technicians.update_own"]`,
		Active:      true,
	}

	if err := tx.Create(&technicianRole).Error; err != nil {
		tx.Rollback()
		logger.WithContext(c).Errorf("Failed to create technician role: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusInternalServerError,
				"Failed to create organization roles",
				map[string]interface{}{
					"code": "ROLE_CREATION_ERROR",
				},
			),
		})
		return
	}

	// Check if user email already exists in this organization (shouldn't happen but extra safety)
	var existingUser models.User
	if err := tx.Where("email = ? AND organization_id = ? AND deleted_at IS NULL", req.Email, org.ID).First(&existingUser).Error; err == nil {
		tx.Rollback()
		logger.WithContext(c).Warnf("Registration failed: email %s already exists in organization %d", req.Email, org.ID)
		c.JSON(http.StatusConflict, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusConflict,
				"Email is already registered",
				map[string]interface{}{
					"code": "EMAIL_ALREADY_EXISTS",
				},
			),
		})
		return
	} else if err != gorm.ErrRecordNotFound {
		tx.Rollback()
		logger.WithContext(c).Errorf("Database error checking email uniqueness: %v", err)
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

	// Create user as organization owner
	user := models.User{
		Base: models.Base{
			OrganizationID: org.ID,
		},
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		RoleID:    ownerRole.ID,
		Active:    true,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		logger.WithContext(c).Errorf("Failed to create user: %v", err)
		
		// Check if it's a duplicate email error
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			c.JSON(http.StatusConflict, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusConflict,
					"Email is already registered",
					map[string]interface{}{
						"code": "EMAIL_ALREADY_EXISTS",
					},
				),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusInternalServerError,
					"Failed to create user",
					map[string]interface{}{
						"code": "USER_CREATION_ERROR",
					},
				),
			})
		}
		return
	}

	// Load the role for the user
	if err := tx.Preload("Role").First(&user, user.ID).Error; err != nil {
		tx.Rollback()
		logger.WithContext(c).Errorf("Failed to load user role: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusInternalServerError,
				"Failed to complete registration",
				map[string]interface{}{
					"code": "USER_LOAD_ERROR",
				},
			),
		})
		return
	}

	// Generate tokens
	accessToken, err := h.jwtService.GenerateAccessToken(user.ID, user.OrganizationID, user.Email, user.Role.Name.String())
	if err != nil {
		tx.Rollback()
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
		tx.Rollback()
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
	
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		logger.WithContext(c).Errorf("Failed to update user login info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusInternalServerError,
				"Failed to complete registration",
				map[string]interface{}{
					"code": "USER_UPDATE_ERROR",
				},
			),
		})
		return
	}

	// Check if panic was recovered before attempting commit
	if panicRecovered {
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		logger.WithContext(c).Errorf("Failed to commit registration transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusInternalServerError,
				"Failed to complete registration",
				map[string]interface{}{
					"code": "TRANSACTION_COMMIT_ERROR",
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

	orgResponse := validation.OrganizationResponse{
		ID:             org.ID,
		Name:           org.Name,
		SubDomain:      org.SubDomain,
		ContactEmail:   org.ContactEmail,
		ContactPhone:   org.ContactPhone,
		LogoURL:        org.LogoURL,
		PrimaryColor:   org.PrimaryColor,
		SecondaryColor: org.SecondaryColor,
		Active:         org.Active,
		PlanType:       org.PlanType,
		CreatedAt:      org.CreatedAt,
		UpdatedAt:      org.UpdatedAt,
	}

	registrationResponse := validation.RegistrationResponse{
		User:         userResponse,
		Organization: orgResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    constants.JWT_ACCESS_TOKEN_EXPIRY,
	}

	logger.WithContext(c).Infof("User %s registered successfully for organization %s", req.Email, req.OrganizationName)
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    registrationResponse,
		"message": "Registration successful",
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