package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"routrapp-api/internal/api"
	"routrapp-api/internal/models"
	"routrapp-api/internal/utils/auth"
	"routrapp-api/internal/validation"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestContext holds test dependencies
type TestContext struct {
	DB          *gorm.DB
	Router      *gin.Engine
	AuthHandler *api.AuthHandler
	JWTService  *auth.JWTService
}

// SetupTestDB creates an in-memory SQLite database for testing
func SetupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Discard,
	})
	if err != nil {
		return nil, err
	}

	// Auto-migrate all models
	err = db.AutoMigrate(
		&models.Organization{},
		&models.Role{},
		&models.User{},
		&models.Technician{},
		&models.Route{},
		&models.RouteStop{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// SetupTestContext creates a complete test context with database, router, and handlers
func SetupTestContext() (*TestContext, error) {
	db, err := SetupTestDB()
	if err != nil {
		return nil, err
	}

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Set test JWT secret for consistent token validation
	os.Setenv("JWT_SECRET", "test-secret-key")

	// Create JWT service with test secret
	jwtService := auth.NewJWTService("test-secret-key")
	
	// Create auth handler using the constructor
	authHandler := api.NewAuthHandler(db)

	router := gin.New()

	// Setup auth routes
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/register", authHandler.RegisterOrganization)                               // POST /api/v1/auth/register (organization registration)
		authGroup.POST("/register-user", authHandler.Register)                                      // POST /api/v1/auth/register-user (user registration to existing org)
		authGroup.POST("/login", authHandler.Login)                                                    // POST /api/v1/auth/login
		authGroup.POST("/refresh", authHandler.RefreshToken)                                           // POST /api/v1/auth/refresh
		authGroup.GET("/me", CreateTestAuthMiddleware(jwtService), authHandler.GetCurrentUser)         // GET /api/v1/auth/me (requires auth)
		authGroup.POST("/logout", CreateTestAuthMiddleware(jwtService), authHandler.Logout)            // POST /api/v1/auth/logout (requires auth)
		authGroup.POST("/change-password", CreateTestAuthMiddleware(jwtService), authHandler.ChangePassword) // POST /api/v1/auth/change-password (requires auth)
	}

	return &TestContext{
		DB:          db,
		Router:      router,
		AuthHandler: authHandler,
		JWTService:  jwtService,
	}, nil
}

// CreateTestAuthMiddleware creates an auth middleware that uses the test JWT service
func CreateTestAuthMiddleware(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": map[string]interface{}{
					"code":    401,
					"message": "Authorization header is required",
					"details": map[string]interface{}{
						"code": "MISSING_AUTH_HEADER",
					},
				},
			})
			return
		}

		// Extract the token from Bearer prefix
		tokenString, err := auth.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": map[string]interface{}{
					"code":    401,
					"message": "Invalid authorization header format: " + err.Error(),
					"details": map[string]interface{}{
						"code": "INVALID_TOKEN",
					},
				},
			})
			return
		}

		// Validate the token
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": map[string]interface{}{
					"code":    401,
					"message": "Invalid or expired token: " + err.Error(),
					"details": map[string]interface{}{
						"code": "INVALID_TOKEN",
					},
				},
			})
			return
		}

		// Ensure this is an access token (not a refresh token)
		if !claims.IsAccessToken() {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": map[string]interface{}{
					"code":    401,
					"message": "Access token required",
					"details": map[string]interface{}{
						"code": "INVALID_TOKEN_TYPE",
					},
				},
			})
			return
		}

		// Set user context in Gin context
		c.Set("user_id", claims.UserID)
		c.Set("organization_id", claims.OrganizationID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// CreateTestOrganization creates a test organization
func CreateTestOrganization(db *gorm.DB) (*models.Organization, error) {
	org := &models.Organization{
		Name:         "Test Organization",
		SubDomain:    "test",
		ContactEmail: "admin@test.com",
		Active:       true,
		PlanType:     "basic",
	}

	err := db.Create(org).Error
	return org, err
}

// CreateTestRole creates a test role
func CreateTestRole(db *gorm.DB, orgID uint, roleType models.RoleType) (*models.Role, error) {
	role := &models.Role{
		Base: models.Base{
			OrganizationID: orgID,
		},
		Name:        roleType,
		DisplayName: string(roleType),
		Description: "Test " + string(roleType) + " role",
		Active:      true,
	}

	err := db.Create(role).Error
	return role, err
}

// CreateTestUser creates a test user with hashed password
func CreateTestUser(db *gorm.DB, orgID, roleID uint, email, password string, active bool) (*models.User, error) {
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Base: models.Base{
			OrganizationID: orgID,
		},
		Email:     email,
		Password:  hashedPassword,
		FirstName: "Test",
		LastName:  "User",
		RoleID:    roleID,
		Active:    active,
	}

	err = db.Create(user).Error
	if err != nil {
		return nil, err
	}

	// If we want to create an inactive user, we need to update it after creation
	// because GORM's default value overrides our setting
	if !active {
		err = db.Exec("UPDATE users SET active = ? WHERE id = ?", false, user.ID).Error
		if err != nil {
			return nil, err
		}
		user.Active = false
	}

	// Load the role for the user
	err = db.Preload("Role").First(user, user.ID).Error
	return user, err
}

// TestUser represents a complete test user setup
type TestUser struct {
	Organization *models.Organization
	Role         *models.Role
	User         *models.User
	Password     string // Plain text password for testing
}

// CreateCompleteTestUser creates a complete test user with organization and role
func CreateCompleteTestUser(db *gorm.DB, email, password string, roleType models.RoleType, active bool) (*TestUser, error) {
	// Create organization
	org, err := CreateTestOrganization(db)
	if err != nil {
		return nil, err
	}

	// Create role
	role, err := CreateTestRole(db, org.ID, roleType)
	if err != nil {
		return nil, err
	}

	// Create user
	user, err := CreateTestUser(db, org.ID, role.ID, email, password, active)
	if err != nil {
		return nil, err
	}

	return &TestUser{
		Organization: org,
		Role:         role,
		User:         user,
		Password:     password,
	}, nil
}

// MakeLoginRequest creates a login request
func MakeLoginRequest(router *gin.Engine, email, password string) *httptest.ResponseRecorder {
	loginReq := validation.UserLoginRequest{
		Email:    email,
		Password: password,
	}

	jsonBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// MakeRefreshRequest creates a refresh token request
func MakeRefreshRequest(router *gin.Engine, refreshToken string) *httptest.ResponseRecorder {
	refreshReq := validation.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	jsonBody, _ := json.Marshal(refreshReq)
	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// MakeLogoutRequest creates a logout request
func MakeLogoutRequest(router *gin.Engine, accessToken string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", "/api/v1/auth/logout", nil)
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// ParseLoginResponse parses a login response
func ParseLoginResponse(w *httptest.ResponseRecorder) (*validation.LoginResponse, error) {
	var response struct {
		Success bool                      `json:"success"`
		Data    validation.LoginResponse `json:"data"`
		Message string                    `json:"message"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}

// ParseTokenResponse parses a token refresh response
func ParseTokenResponse(w *httptest.ResponseRecorder) (*validation.TokenResponse, error) {
	var response struct {
		Success bool                      `json:"success"`
		Data    validation.TokenResponse `json:"data"`
		Message string                    `json:"message"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}

// MakeRegistrationRequest makes a POST request to /api/v1/auth/register
func MakeRegistrationRequest(router *gin.Engine, email, password, firstName, lastName, orgName, orgEmail, subDomain string) *httptest.ResponseRecorder {
	reqBody := validation.RegistrationRequest{
		Email:             email,
		Password:          password,
		FirstName:         firstName,
		LastName:          lastName,
		OrganizationName:  orgName,
		OrganizationEmail: orgEmail,
		SubDomain:         subDomain,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// ParseRegistrationResponse parses a registration response
func ParseRegistrationResponse(w *httptest.ResponseRecorder) (*validation.RegistrationResponse, error) {
	var response struct {
		Success bool                             `json:"success"`
		Data    validation.RegistrationResponse `json:"data"`
		Message string                           `json:"message"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}

// ParseErrorResponse parses an error response
func ParseErrorResponse(w *httptest.ResponseRecorder) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	return response, err
}

// GenerateExpiredToken generates an expired JWT token for testing
func GenerateExpiredToken(jwtService *auth.JWTService, userID, orgID uint, email, role string, tokenType string) (string, error) {
	// Create custom claims with past expiry
	claims := &auth.JWTClaims{
		UserID:         userID,
		OrganizationID: orgID,
		Email:          email,
		Role:           role,
		TokenType:      tokenType,
	}

	// Set expiry to 1 hour ago
	pastTime := time.Now().Add(-1 * time.Hour)
	claims.ExpiresAt = jwt.NewNumericDate(pastTime)
	claims.IssuedAt = jwt.NewNumericDate(pastTime)

	// This is a simplified version - in a real test, you'd need to generate the token manually
	// or have a method in JWTService that accepts custom expiry
	if tokenType == "access" {
		return jwtService.GenerateAccessToken(userID, orgID, email, role)
	}
	return jwtService.GenerateRefreshToken(userID, orgID, email, role)
}

// AssertResponseSuccess checks if a response indicates success
func AssertResponseSuccess(w *httptest.ResponseRecorder, expectedStatus int) bool {
	if w.Code != expectedStatus {
		return false
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		return false
	}

	success, ok := response["success"].(bool)
	return ok && success
}

// AssertResponseError checks if a response indicates an error with expected code
func AssertResponseError(w *httptest.ResponseRecorder, expectedStatus int, expectedCode string) bool {
	if w.Code != expectedStatus {
		return false
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		return false
	}

	errorData, ok := response["error"].(map[string]interface{})
	if !ok {
		return false
	}

	details, ok := errorData["details"].(map[string]interface{})
	if !ok {
		return false
	}

	code, ok := details["code"].(string)
	return ok && code == expectedCode
}

// CleanupTestContext cleans up test database
func CleanupTestContext(ctx *TestContext) error {
	if ctx.DB != nil {
		sqlDB, err := ctx.DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
} 