package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"routrapp-api/internal/errors"
	"routrapp-api/internal/utils/auth"
	"routrapp-api/internal/utils/constants"
)

// TenantContext defines the organization context structure
type TenantContext struct {
	OrganizationID uint
	SubDomain      string
}

// TenantMiddleware extracts organization ID from JWT or subdomain and sets it in context
func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tenantCtx TenantContext

		// First try to get organization ID from JWT if user is authenticated
		if authHeader := c.GetHeader("Authorization"); authHeader != "" {
			jwtService := auth.DefaultJWTService()
			
			// Extract token from Bearer header
			if tokenString, err := auth.ExtractTokenFromHeader(authHeader); err == nil {
				// Validate the token using our JWT service
				if claims, err := jwtService.ValidateToken(tokenString); err == nil {
					tenantCtx.OrganizationID = claims.OrganizationID
				}
			}
		}

		// If no organization ID from JWT, try to get it from subdomain
		if tenantCtx.OrganizationID == 0 {
			host := c.Request.Host
			subdomain := extractSubdomain(host)
			
			if subdomain != "" {
				tenantCtx.SubDomain = subdomain
				// In a real app, you would look up the organization ID from the database
				// using the subdomain
				// For now, we'll set a placeholder error
				c.Set(constants.TENANT_CONTEXT_KEY, tenantCtx)
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": errors.NewAppErrorWithDetails(
						http.StatusBadRequest,
						"Organization not found for subdomain: "+subdomain,
						map[string]interface{}{
							"code": "TENANT_NOT_FOUND",
						},
					),
				})
				return
			}
		}

		// If we still don't have an organization ID, check if it's explicitly set in query params
		// (useful for public APIs or testing)
		if tenantCtx.OrganizationID == 0 && c.Query("organization_id") != "" {
			if orgID, err := strconv.ParseUint(c.Query("organization_id"), 10, 32); err == nil {
				tenantCtx.OrganizationID = uint(orgID)
			}
		}

		// Set the tenant context in Gin's context
		c.Set(constants.TENANT_CONTEXT_KEY, tenantCtx)

		// If we have a path that requires tenant context but we don't have it, abort
		if requiresTenantContext(c.FullPath()) && tenantCtx.OrganizationID == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusBadRequest,
					"Organization context is required for this endpoint",
					map[string]interface{}{
						"code": "TENANT_REQUIRED",
					},
				),
			})
			return
		}

		c.Next()
	}
}

// GetTenantContext retrieves the tenant context from Gin's context
func GetTenantContext(c *gin.Context) (TenantContext, bool) {
	if tenantCtx, exists := c.Get(constants.TENANT_CONTEXT_KEY); exists {
		if ctx, ok := tenantCtx.(TenantContext); ok {
			return ctx, true
		}
	}
	return TenantContext{}, false
}

// extractSubdomain extracts subdomain from host
func extractSubdomain(host string) string {
	// Remove port if present
	if colonIndex := strings.Index(host, ":"); colonIndex != -1 {
		host = host[:colonIndex]
	}

	parts := strings.Split(host, ".")
	
	// Handle localhost separately
	if host == "localhost" || len(parts) < 3 {
		return ""
	}
	
	// Return first part as subdomain
	return parts[0]
}

// requiresTenantContext determines if an endpoint requires tenant context
func requiresTenantContext(path string) bool {
	// Public endpoints that don't require tenant context
	publicPaths := []string{
		"/api/v1/health",
		"/api/v1/auth/login",
		"/api/v1/auth/register",
	}

	for _, publicPath := range publicPaths {
		if strings.HasPrefix(path, publicPath) {
			return false
		}
	}

	// All other API paths require tenant context
	return strings.HasPrefix(path, "/api/v1/")
}