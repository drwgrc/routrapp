package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"routrapp-api/internal/errors"
	"routrapp-api/internal/utils/auth"
	"routrapp-api/internal/utils/constants"
)

// AuthMiddleware validates JWT tokens and sets user context using default JWT service
func AuthMiddleware() gin.HandlerFunc {
	jwtService := auth.DefaultJWTService()
	
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusUnauthorized,
					"Authorization header is required",
					map[string]interface{}{
						"code": "MISSING_AUTH_HEADER",
					},
				),
			})
			return
		}

		// Extract the token from Bearer prefix
		tokenString, err := auth.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusUnauthorized,
					"Invalid authorization header format: "+err.Error(),
					map[string]interface{}{
						"code": "INVALID_AUTH_HEADER",
					},
				),
			})
			return
		}

		// Validate the token
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusUnauthorized,
					"Invalid or expired token: "+err.Error(),
					map[string]interface{}{
						"code": "INVALID_TOKEN",
					},
				),
			})
			return
		}

		// Ensure this is an access token (not a refresh token)
		if !claims.IsAccessToken() {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusUnauthorized,
					"Access token required",
					map[string]interface{}{
						"code": "INVALID_TOKEN_TYPE",
					},
				),
			})
			return
		}

		// Set user context in Gin context
		c.Set(constants.USER_CONTEXT_KEY, claims.GetUserContext())
		
		// Also set individual claims for easy access
		c.Set("user_id", claims.UserID)
		c.Set("organization_id", claims.OrganizationID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// AuthMiddlewareWithJWT validates JWT tokens and sets user context using provided JWT service
func AuthMiddlewareWithJWT(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusUnauthorized,
					"Authorization header is required",
					map[string]interface{}{
						"code": "MISSING_AUTH_HEADER",
					},
				),
			})
			return
		}

		// Extract the token from Bearer prefix
		tokenString, err := auth.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusUnauthorized,
					"Invalid authorization header format: "+err.Error(),
					map[string]interface{}{
						"code": "INVALID_AUTH_HEADER",
					},
				),
			})
			return
		}

		// Validate the token
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusUnauthorized,
					"Invalid or expired token: "+err.Error(),
					map[string]interface{}{
						"code": "INVALID_TOKEN",
					},
				),
			})
			return
		}

		// Ensure this is an access token (not a refresh token)
		if !claims.IsAccessToken() {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusUnauthorized,
					"Access token required",
					map[string]interface{}{
						"code": "INVALID_TOKEN_TYPE",
					},
				),
			})
			return
		}

		// Set user context in Gin context
		c.Set(constants.USER_CONTEXT_KEY, claims.GetUserContext())
		
		// Also set individual claims for easy access
		c.Set("user_id", claims.UserID)
		c.Set("organization_id", claims.OrganizationID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// OptionalAuthMiddleware validates JWT tokens if present but doesn't require them
func OptionalAuthMiddleware() gin.HandlerFunc {
	jwtService := auth.DefaultJWTService()
	
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No auth header, continue without setting user context
			c.Next()
			return
		}

		// Extract the token from Bearer prefix
		tokenString, err := auth.ExtractTokenFromHeader(authHeader)
		if err != nil {
			// Invalid format, continue without setting user context
			c.Next()
			return
		}

		// Validate the token
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			// Invalid token, continue without setting user context
			c.Next()
			return
		}

		// Ensure this is an access token (not a refresh token)
		if !claims.IsAccessToken() {
			// Wrong token type, continue without setting user context
			c.Next()
			return
		}

		// Set user context in Gin context
		c.Set(constants.USER_CONTEXT_KEY, claims.GetUserContext())
		
		// Also set individual claims for easy access
		c.Set("user_id", claims.UserID)
		c.Set("organization_id", claims.OrganizationID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// GetUserContext retrieves the user context from Gin's context
func GetUserContext(c *gin.Context) (map[string]interface{}, bool) {
	if userCtx, exists := c.Get(constants.USER_CONTEXT_KEY); exists {
		if ctx, ok := userCtx.(map[string]interface{}); ok {
			return ctx, true
		}
	}
	return nil, false
}

// GetUserID retrieves the user ID from Gin's context
func GetUserID(c *gin.Context) (uint, bool) {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(uint); ok {
			return id, true
		}
	}
	return 0, false
}

// GetOrganizationID retrieves the organization ID from Gin's context
func GetOrganizationID(c *gin.Context) (uint, bool) {
	if orgID, exists := c.Get("organization_id"); exists {
		if id, ok := orgID.(uint); ok {
			return id, true
		}
	}
	return 0, false
}

// GetUserEmail retrieves the user email from Gin's context
func GetUserEmail(c *gin.Context) (string, bool) {
	if email, exists := c.Get("user_email"); exists {
		if e, ok := email.(string); ok {
			return e, true
		}
	}
	return "", false
}

// GetUserRole retrieves the user role from Gin's context
func GetUserRole(c *gin.Context) (string, bool) {
	if role, exists := c.Get("user_role"); exists {
		if r, ok := role.(string); ok {
			return r, true
		}
	}
	return "", false
}

// RequireRole creates middleware that requires a specific role
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := GetUserRole(c)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusUnauthorized,
					"User role not found in context",
					map[string]interface{}{
						"code": "MISSING_USER_ROLE",
					},
				),
			})
			return
		}

		if userRole != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusForbidden,
					"Insufficient permissions. Required role: "+requiredRole,
					map[string]interface{}{
						"code":          "INSUFFICIENT_PERMISSIONS",
						"required_role": requiredRole,
						"user_role":     userRole,
					},
				),
			})
			return
		}

		c.Next()
	}
}

// RequireOwner creates middleware that requires owner role
func RequireOwner() gin.HandlerFunc {
	return RequireRole("owner")
}

// RequireAuthentication checks if user is authenticated (any valid role)
func RequireAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := GetUserContext(c)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
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

		c.Next()
	}
} 