package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"routrapp-api/internal/errors"
	"routrapp-api/internal/models"
)

// PermissionChecker interface allows for different permission checking strategies
type PermissionChecker interface {
	HasPermission(userRole, userID, organizationID uint, permission string) bool
}

// DefaultPermissionChecker uses the default role-based permission system
type DefaultPermissionChecker struct{}

// HasPermission checks if a user has a specific permission based on their role
func (dpc *DefaultPermissionChecker) HasPermission(userRole, userID, organizationID uint, permission string) bool {
	// Get default permissions for the role type
	var roleType models.RoleType
	switch userRole {
	case 1: // Assuming role ID 1 is owner - this could be improved with a role lookup
		roleType = models.RoleTypeOwner
	case 2: // Assuming role ID 2 is technician
		roleType = models.RoleTypeTechnician
	default:
		return false
	}

	defaultPerms := models.GetDefaultPermissions(roleType)
	
	// Check if permission matches any of the default permissions
	for _, perm := range defaultPerms {
		if permissionMatches(perm, permission) {
			return true
		}
	}

	return false
}

// permissionMatches checks if a stored permission matches the requested permission
// Supports wildcard permissions (e.g., "routes.*" matches "routes.read")
func permissionMatches(storedPerm, requestedPerm string) bool {
	// Exact match
	if storedPerm == requestedPerm {
		return true
	}

	// Wildcard match (e.g., "routes.*" matches "routes.read")
	if strings.HasSuffix(storedPerm, ".*") {
		prefix := strings.TrimSuffix(storedPerm, ".*")
		return strings.HasPrefix(requestedPerm, prefix+".")
	}

	// Global wildcard
	if storedPerm == "*" {
		return true
	}

	return false
}

// RequirePermission creates middleware that requires a specific permission
func RequirePermission(permission string) gin.HandlerFunc {
	return RequirePermissionWithChecker(permission, &DefaultPermissionChecker{})
}

// RequirePermissionWithChecker creates middleware that requires a specific permission using a custom checker
func RequirePermissionWithChecker(permission string, checker PermissionChecker) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ensure user is authenticated first
		userID, exists := GetUserID(c)
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

		organizationID, exists := GetOrganizationID(c)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusUnauthorized,
					"Organization context required",
					map[string]interface{}{
						"code": "ORGANIZATION_REQUIRED",
					},
				),
			})
			return
		}

		// Get user role for permission checking
		// Note: This is a simplified approach - in production you might want to get the actual role ID
		userRole, exists := GetUserRole(c)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusUnauthorized,
					"User role not found",
					map[string]interface{}{
						"code": "MISSING_USER_ROLE",
					},
				),
			})
			return
		}

		// Convert role string to role ID for checker
		var roleID uint
		switch userRole {
		case "owner":
			roleID = 1 // This should be improved with actual role lookup
		case "technician":
			roleID = 2
		default:
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusForbidden,
					"Unknown user role",
					map[string]interface{}{
						"code": "UNKNOWN_ROLE",
						"role": userRole,
					},
				),
			})
			return
		}

		// Check if user has the required permission
		if !checker.HasPermission(roleID, userID, organizationID, permission) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusForbidden,
					"Insufficient permissions. Required permission: "+permission,
					map[string]interface{}{
						"code":                "INSUFFICIENT_PERMISSIONS",
						"required_permission": permission,
						"user_role":          userRole,
					},
				),
			})
			return
		}

		c.Next()
	}
}

// RequireAnyPermission creates middleware that requires any one of the specified permissions
func RequireAnyPermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Always require authentication, even if no specific permissions are specified
		userID, exists := GetUserID(c)
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

		organizationID, exists := GetOrganizationID(c)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusUnauthorized,
					"Organization context required",
					map[string]interface{}{
						"code": "ORGANIZATION_REQUIRED",
					},
				),
			})
			return
		}

		userRole, exists := GetUserRole(c)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusUnauthorized,
					"User role not found",
					map[string]interface{}{
						"code": "MISSING_USER_ROLE",
					},
				),
			})
			return
		}

		// If no permissions specified, allow any authenticated user
		if len(permissions) == 0 {
			c.Next()
			return
		}

		// Convert role string to role ID for checker
		var roleID uint
		switch userRole {
		case "owner":
			roleID = 1
		case "technician":
			roleID = 2
		default:
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusForbidden,
					"Unknown user role",
					map[string]interface{}{
						"code": "UNKNOWN_ROLE",
						"role": userRole,
					},
				),
			})
			return
		}

		checker := &DefaultPermissionChecker{}

		// Check each permission until one is found
		for _, permission := range permissions {
			if checker.HasPermission(roleID, userID, organizationID, permission) {
				// User has this permission, continue
				c.Next()
				return
			}
		}

		// User doesn't have any of the required permissions
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": errors.NewAppErrorWithDetails(
				http.StatusForbidden,
				"Insufficient permissions. Required one of: "+strings.Join(permissions, ", "),
				map[string]interface{}{
					"code":                 "INSUFFICIENT_PERMISSIONS",
					"required_permissions": permissions,
					"user_role":           userRole,
				},
			),
		})
	}
}

// RequireAllPermissions creates middleware that requires all specified permissions
func RequireAllPermissions(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(permissions) == 0 {
			c.Next()
			return
		}

		userID, exists := GetUserID(c)
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

		organizationID, exists := GetOrganizationID(c)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusUnauthorized,
					"Organization context required",
					map[string]interface{}{
						"code": "ORGANIZATION_REQUIRED",
					},
				),
			})
			return
		}

		userRole, exists := GetUserRole(c)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusUnauthorized,
					"User role not found",
					map[string]interface{}{
						"code": "MISSING_USER_ROLE",
					},
				),
			})
			return
		}

		var roleID uint
		switch userRole {
		case "owner":
			roleID = 1
		case "technician":
			roleID = 2
		default:
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusForbidden,
					"Unknown user role",
					map[string]interface{}{
						"code": "UNKNOWN_ROLE",
						"role": userRole,
					},
				),
			})
			return
		}

		checker := &DefaultPermissionChecker{}

		// Check all permissions - all must be satisfied
		for _, permission := range permissions {
			if !checker.HasPermission(roleID, userID, organizationID, permission) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": errors.NewAppErrorWithDetails(
						http.StatusForbidden,
						"Insufficient permissions. Missing permission: "+permission,
						map[string]interface{}{
							"code":               "INSUFFICIENT_PERMISSIONS",
							"missing_permission": permission,
							"user_role":         userRole,
						},
					),
				})
				return
			}
		}

		c.Next()
	}
}

// Convenience middleware functions for common permission patterns

// RequireOrganizationAccess requires access to organization data
func RequireOrganizationAccess() gin.HandlerFunc {
	return RequirePermission("organizations.read")
}

// RequireUserManagement requires user management permissions
func RequireUserManagement() gin.HandlerFunc {
	return RequirePermission("users.manage")
}

// RequireTechnicianManagement requires technician management permissions
func RequireTechnicianManagement() gin.HandlerFunc {
	return RequirePermission("technicians.manage")
}

// RequireRouteAccess requires route access permissions
func RequireRouteAccess() gin.HandlerFunc {
	return RequireAnyPermission("routes.read", "routes.*")
}

// RequireRouteManagement requires route management permissions
func RequireRouteManagement() gin.HandlerFunc {
	return RequirePermission("routes.manage")
}

// RequireOwnerOnly requires owner role (backward compatibility)
func RequireOwnerOnly() gin.HandlerFunc {
	return RequireRole("owner")
}

// RequireResourceOwnership creates middleware that checks if user owns a specific resource
// This is useful for endpoints like /users/{id} where users should only access their own data
func RequireResourceOwnership(resourceParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := GetUserID(c)
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

		userRole, exists := GetUserRole(c)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusUnauthorized,
					"User role not found",
					map[string]interface{}{
						"code": "MISSING_USER_ROLE",
					},
				),
			})
			return
		}

		// Owners can access any resource
		if userRole == "owner" {
			c.Next()
			return
		}

		// For non-owners, check if they're accessing their own resource
		resourceID := c.Param(resourceParam)
		if resourceID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusBadRequest,
					"Resource ID required",
					map[string]interface{}{
						"code": "MISSING_RESOURCE_ID",
					},
				),
			})
			return
		}

		// Convert userID to string for comparison
		userIDStr := strconv.Itoa(int(userID))
		if resourceID != userIDStr {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": errors.NewAppErrorWithDetails(
					http.StatusForbidden,
					"Access denied. You can only access your own resources",
					map[string]interface{}{
						"code": "RESOURCE_ACCESS_DENIED",
					},
				),
			})
			return
		}

		c.Next()
	}
}

// HasPermission checks if the current user has a specific permission
// This is a helper function that can be used within handlers
func HasPermission(c *gin.Context, permission string) bool {
	userID, exists := GetUserID(c)
	if !exists {
		return false
	}

	organizationID, exists := GetOrganizationID(c)
	if !exists {
		return false
	}

	userRole, exists := GetUserRole(c)
	if !exists {
		return false
	}

	var roleID uint
	switch userRole {
	case "owner":
		roleID = 1
	case "technician":
		roleID = 2
	default:
		return false
	}

	checker := &DefaultPermissionChecker{}
	return checker.HasPermission(roleID, userID, organizationID, permission)
} 