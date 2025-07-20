# Middleware Package

This package contains all HTTP middleware for the RouteApp API, providing authentication, authorization, tenant isolation, logging, and error handling.

## Available Middleware

### Authentication Middleware (`auth.go`)

Handles JWT token validation and user context setup.

- `AuthMiddleware()` - Validates JWT tokens and sets user context (required)
- `OptionalAuthMiddleware()` - Validates JWT tokens if present (optional)
- `RequireAuthentication()` - Ensures user is authenticated
- `RequireRole(role)` - Requires specific role (owner/technician)
- `RequireOwner()` - Requires owner role

### RBAC Middleware (`rbac.go`)

Provides fine-grained permission-based access control.

#### Core Functions

- `RequirePermission(permission)` - Requires specific permission
- `RequireAnyPermission(permissions...)` - Requires any of the specified permissions
- `RequireAllPermissions(permissions...)` - Requires all specified permissions
- `RequireResourceOwnership(resourceParam)` - Resource ownership validation

#### Convenience Functions

- `RequireOrganizationAccess()` - Access to organization data
- `RequireUserManagement()` - User management permissions
- `RequireTechnicianManagement()` - Technician management permissions
- `RequireRouteAccess()` - Route reading permissions
- `RequireRouteManagement()` - Route management permissions

### Tenant Middleware (`tenant.go`)

Handles multi-tenant organization context.

- `TenantMiddleware()` - Extracts organization context from JWT or subdomain

### Other Middleware

- `CORSMiddleware(cfg)` - Cross-origin resource sharing
- `LoggerMiddleware()` - HTTP request logging
- `ErrorHandlerMiddleware()` - Centralized error handling
- `RecoveryMiddleware()` - Panic recovery

## RBAC Usage Examples

### Basic Permission Checking

```go
package routes

import (
    "github.com/gin-gonic/gin"
    "routrapp-api/internal/middleware"
)

func SetupRoutes(r *gin.Engine) {
    api := r.Group("/api/v1")

    // Apply authentication to all routes
    api.Use(middleware.AuthMiddleware())

    // Routes that require specific permissions
    routes := api.Group("/routes")
    {
        // Anyone with route access can read routes
        routes.GET("/", middleware.RequireRouteAccess(), listRoutes)

        // Only users with route management can create/update/delete
        routes.POST("/", middleware.RequirePermission("routes.create"), createRoute)
        routes.PUT("/:id", middleware.RequirePermission("routes.update"), updateRoute)
        routes.DELETE("/:id", middleware.RequirePermission("routes.delete"), deleteRoute)
    }

    // User management routes (owners only)
    users := api.Group("/users")
    users.Use(middleware.RequireUserManagement())
    {
        users.GET("/", listUsers)
        users.POST("/", createUser)
        users.PUT("/:id", updateUser)
        users.DELETE("/:id", deleteUser)
    }
}
```

### Resource Ownership Patterns

```go
// Users can only access their own profile unless they're owners
api.GET("/users/:id",
    middleware.RequireResourceOwnership("id"),
    getUserProfile)

// Technicians can only update their own data
api.PUT("/technicians/:id",
    middleware.RequireAnyPermission("technicians.manage", "technicians.update_own"),
    middleware.RequireResourceOwnership("id"),
    updateTechnician)
```

### Complex Permission Logic

```go
// Route assignment - requires either route management OR being the assigned technician
api.PUT("/routes/:id/assign",
    middleware.RequireAnyPermission(
        "routes.manage",
        "routes.assign_self",
    ),
    assignRoute)

// Route completion - requires multiple permissions
api.POST("/routes/:id/complete",
    middleware.RequireAllPermissions(
        "routes.update_status",
        "routes.complete",
    ),
    completeRoute)
```

### Custom Permission Checking in Handlers

```go
func updateRouteHandler(c *gin.Context) {
    routeID := c.Param("id")

    // Check if user can manage all routes OR owns this specific route
    canManageAll := middleware.HasPermission(c, "routes.manage")
    if !canManageAll {
        // Check if user is the assigned technician for this route
        route, err := getRouteByID(routeID)
        if err != nil {
            c.JSON(404, gin.H{"error": "Route not found"})
            return
        }

        userID, _ := middleware.GetUserID(c)
        if route.TechnicianID == nil || *route.TechnicianID != userID {
            c.JSON(403, gin.H{"error": "Access denied"})
            return
        }
    }

    // Proceed with route update
    // ...
}
```

## Permission System

### Permission Format

Permissions follow a hierarchical dot notation:

- `organizations.*` - All organization permissions
- `organizations.read` - Read organization data
- `organizations.update` - Update organization data
- `users.*` - All user permissions
- `users.create` - Create users
- `users.read` - Read user data
- `routes.manage` - Full route management
- `routes.read` - Read routes
- `routes.update_status` - Update route status

### Wildcard Matching

- `routes.*` matches `routes.read`, `routes.create`, `routes.update`, etc.
- `*` matches any permission (super admin)

### Default Role Permissions

#### Owner Role

```json
["organizations.*", "users.*", "technicians.*", "routes.*", "roles.*"]
```

#### Technician Role

```json
[
  "routes.read",
  "routes.update_status",
  "technicians.read_own",
  "technicians.update_own"
]
```

## Middleware Chain Order

The recommended middleware order for protected routes:

```go
r.Use(middleware.LoggerMiddleware())
r.Use(middleware.RecoveryMiddleware())
r.Use(middleware.CORSMiddleware(cfg))
r.Use(middleware.ErrorHandlerMiddleware())

api := r.Group("/api/v1")
api.Use(middleware.TenantMiddleware())
api.Use(middleware.AuthMiddleware())

// Now add specific permission middleware to route groups
protected := api.Group("/protected")
protected.Use(middleware.RequirePermission("some.permission"))
```

## Error Responses

All middleware returns consistent error responses:

```json
{
  "error": {
    "code": "INSUFFICIENT_PERMISSIONS",
    "message": "Insufficient permissions. Required permission: routes.create",
    "details": {
      "required_permission": "routes.create",
      "user_role": "technician"
    }
  }
}
```

Common error codes:

- `AUTHENTICATION_REQUIRED` - User not authenticated
- `INSUFFICIENT_PERMISSIONS` - User lacks required permission
- `ORGANIZATION_REQUIRED` - Organization context missing
- `MISSING_USER_ROLE` - User role not found in context
- `RESOURCE_ACCESS_DENIED` - User cannot access specific resource
- `UNKNOWN_ROLE` - Invalid role type

## Testing RBAC

### Unit Testing Middleware

```go
func TestRequirePermission(t *testing.T) {
    router := gin.New()
    router.Use(middleware.AuthMiddleware())
    router.GET("/test", middleware.RequirePermission("routes.read"), func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "success"})
    })

    // Test with owner (should have permission)
    req := httptest.NewRequest("GET", "/test", nil)
    req.Header.Set("Authorization", "Bearer " + generateOwnerToken())
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    assert.Equal(t, 200, w.Code)

    // Test with technician (should have permission for routes.read)
    req = httptest.NewRequest("GET", "/test", nil)
    req.Header.Set("Authorization", "Bearer " + generateTechnicianToken())
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)
    assert.Equal(t, 200, w.Code)
}
```

### Integration Testing

```go
func TestRoutePermissions(t *testing.T) {
    // Test that owners can create routes
    token := loginAsOwner()
    resp := createRoute(token, routeData)
    assert.Equal(t, 201, resp.StatusCode)

    // Test that technicians cannot create routes
    token = loginAsTechnician()
    resp = createRoute(token, routeData)
    assert.Equal(t, 403, resp.StatusCode)

    // Test that technicians can read routes
    resp = getRoutes(token)
    assert.Equal(t, 200, resp.StatusCode)
}
```

## Custom Permission Checkers

You can implement custom permission checking logic:

```go
type DatabasePermissionChecker struct {
    db *gorm.DB
}

func (dpc *DatabasePermissionChecker) HasPermission(userRole, userID, organizationID uint, permission string) bool {
    var role models.Role
    if err := dpc.db.First(&role, userRole).Error; err != nil {
        return false
    }

    return role.HasPermission(permission)
}

// Use custom checker
router.GET("/custom",
    middleware.RequirePermissionWithChecker("custom.permission", &DatabasePermissionChecker{db: db}),
    handler)
```

## Performance Considerations

- Permission checks are done in-memory using role defaults
- For high-traffic endpoints, consider caching permission results
- Use `HasPermission()` helper in handlers for conditional logic
- Avoid deeply nested permission chains

## Security Best Practices

1. **Principle of Least Privilege**: Grant minimal permissions needed
2. **Default Deny**: Require explicit permissions for all actions
3. **Validate Context**: Always check authentication before authorization
4. **Resource Ownership**: Use `RequireResourceOwnership()` for user-specific data
5. **Audit Permissions**: Log permission checks for security auditing
6. **Regular Review**: Periodically review and update permission assignments
