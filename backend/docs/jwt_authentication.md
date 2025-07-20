# JWT Authentication System

This document describes the JWT (JSON Web Token) authentication system implemented for the RoutrApp API.

## Overview

The JWT authentication system provides secure token-based authentication for API endpoints. It supports both access tokens (short-lived) and refresh tokens (long-lived) with proper validation and multi-tenant organization context.

## Architecture

### Components

1. **JWT Service** (`internal/utils/auth/jwt.go`)

   - Token generation (access and refresh tokens)
   - Token validation and parsing
   - Claims extraction and verification

2. **Auth Middleware** (`internal/middleware/auth.go`)

   - Request authentication
   - User context injection
   - Role-based access control

3. **Constants** (`internal/utils/constants/constants.go`)
   - JWT configuration (secret key, expiry times)
   - Environment variable support

## Token Structure

### JWT Claims

```json
{
  "user_id": 123,
  "organization_id": 456,
  "email": "user@example.com",
  "role": "owner",
  "token_type": "access",
  "iss": "routrapp-api",
  "aud": ["routrapp-frontend"],
  "sub": "123",
  "iat": 1640995200,
  "exp": 1640996100
}
```

### Token Types

- **Access Token**: Short-lived (15 minutes) for API access
- **Refresh Token**: Long-lived (7 days) for obtaining new access tokens

## Configuration

### Environment Variables

```bash
# JWT Secret Key (required in production)
JWT_SECRET=your-super-secret-jwt-key-here
```

### Default Values

- Access token expiry: 15 minutes
- Refresh token expiry: 7 days
- Default secret: "dev-secret-key-change-in-production" (development only)

## Usage

### Generating Tokens

```go
import "routrapp-api/internal/utils/auth"

jwtService := auth.DefaultJWTService()

// Generate access token
accessToken, err := jwtService.GenerateAccessToken(
    userID, organizationID, email, role)

// Generate refresh token
refreshToken, err := jwtService.GenerateRefreshToken(
    userID, organizationID, email, role)
```

### Validating Tokens

```go
claims, err := jwtService.ValidateToken(tokenString)
if err != nil {
    // Handle invalid token
}

// Access token data
userID := claims.UserID
organizationID := claims.OrganizationID
email := claims.Email
role := claims.Role
```

### Middleware Usage

#### Required Authentication

```go
// Protect routes that require authentication
protected := router.Group("/api/v1/protected")
protected.Use(middleware.AuthMiddleware())
{
    protected.GET("/profile", profileHandler)
    protected.POST("/routes", createRouteHandler)
}
```

#### Optional Authentication

```go
// Routes where authentication is optional
public := router.Group("/api/v1/public")
public.Use(middleware.OptionalAuthMiddleware())
{
    public.GET("/health", healthHandler)
}
```

#### Role-Based Access Control

```go
// Owner-only endpoints
ownerOnly := router.Group("/api/v1/admin")
ownerOnly.Use(middleware.AuthMiddleware())
ownerOnly.Use(middleware.RequireOwner())
{
    ownerOnly.POST("/users", createUserHandler)
    ownerOnly.DELETE("/users/:id", deleteUserHandler)
}

// Custom role requirements
customRole := router.Group("/api/v1/custom")
customRole.Use(middleware.AuthMiddleware())
customRole.Use(middleware.RequireRole("technician"))
{
    customRole.GET("/routes/assigned", getAssignedRoutesHandler)
}
```

### Accessing User Context in Handlers

```go
func profileHandler(c *gin.Context) {
    // Get individual values
    userID, exists := middleware.GetUserID(c)
    if !exists {
        c.JSON(401, gin.H{"error": "User not authenticated"})
        return
    }

    email, _ := middleware.GetUserEmail(c)
    role, _ := middleware.GetUserRole(c)
    organizationID, _ := middleware.GetOrganizationID(c)

    // Or get full user context
    userContext, exists := middleware.GetUserContext(c)

    c.JSON(200, gin.H{
        "user_id":         userID,
        "email":           email,
        "role":            role,
        "organization_id": organizationID,
        "context":         userContext,
    })
}
```

## Security Features

### Token Validation

- **Signature verification**: Uses HMAC-SHA256 with secret key
- **Expiration checking**: Automatic token expiry validation
- **Token type validation**: Ensures access tokens are used for API access
- **Signing method validation**: Prevents algorithm confusion attacks

### Header Extraction

- **Bearer token format**: Requires "Bearer " prefix
- **Input validation**: Validates authorization header format
- **Error handling**: Proper error messages for debugging

### Multi-Tenant Support

- **Organization isolation**: Each token contains organization_id
- **Tenant context**: Automatic tenant context injection
- **Cross-tenant protection**: Prevents access to other organizations' data

## Error Codes

| Code                       | Description                                        | HTTP Status |
| -------------------------- | -------------------------------------------------- | ----------- |
| `MISSING_AUTH_HEADER`      | Authorization header not provided                  | 401         |
| `INVALID_AUTH_HEADER`      | Invalid authorization header format                | 401         |
| `INVALID_TOKEN`            | Token validation failed                            | 401         |
| `INVALID_TOKEN_TYPE`       | Wrong token type (e.g., refresh instead of access) | 401         |
| `INSUFFICIENT_PERMISSIONS` | User role doesn't have required permissions        | 403         |
| `AUTHENTICATION_REQUIRED`  | Endpoint requires authentication                   | 401         |

## Best Practices

### Security

1. **Use environment variables** for JWT secret in production
2. **Rotate JWT secrets** regularly in production
3. **Use HTTPS** to prevent token interception
4. **Implement token blacklisting** for logout functionality
5. **Monitor token usage** for suspicious activity

### Performance

1. **Cache JWT service instances** to avoid recreation
2. **Use short-lived access tokens** with refresh token rotation
3. **Implement rate limiting** on authentication endpoints
4. **Consider JWT blacklisting** for immediate revocation

### Development

1. **Use longer tokens in development** for easier debugging
2. **Log authentication failures** for troubleshooting
3. **Validate token structure** in tests
4. **Mock JWT service** in unit tests

## Integration with Frontend

### Login Response

```json
{
  "user": {
    "id": 123,
    "email": "user@example.com",
    "role": "owner"
  },
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 900
}
```

### API Request Headers

```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Token Refresh

```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

## Testing

### Unit Tests

```go
func TestJWTGeneration(t *testing.T) {
    jwtService := auth.NewJWTService("test-secret")

    token, err := jwtService.GenerateAccessToken(1, 1, "test@example.com", "owner")
    assert.NoError(t, err)
    assert.NotEmpty(t, token)

    claims, err := jwtService.ValidateToken(token)
    assert.NoError(t, err)
    assert.Equal(t, uint(1), claims.UserID)
    assert.True(t, claims.IsAccessToken())
}
```

### Integration Tests

```go
func TestAuthMiddleware(t *testing.T) {
    router := gin.New()
    router.Use(middleware.AuthMiddleware())
    router.GET("/protected", func(c *gin.Context) {
        userID, _ := middleware.GetUserID(c)
        c.JSON(200, gin.H{"user_id": userID})
    })

    // Test with valid token
    req := httptest.NewRequest("GET", "/protected", nil)
    req.Header.Set("Authorization", "Bearer "+validToken)

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)
}
```

## Migration from Existing Systems

If migrating from an existing authentication system:

1. **Dual authentication**: Support both old and new systems temporarily
2. **Token migration**: Provide endpoint to exchange old tokens for JWT
3. **User migration**: Update user records with JWT-compatible fields
4. **Frontend updates**: Update frontend to use JWT format
5. **Gradual rollout**: Phase out old system after successful migration

## Troubleshooting

### Common Issues

1. **"Invalid token"**: Check JWT secret configuration
2. **"Token expired"**: Verify system time sync
3. **"Wrong token type"**: Ensure access tokens are used for API calls
4. **"Missing auth header"**: Verify frontend sends Authorization header

### Debug Mode

Enable debug logging to troubleshoot authentication issues:

```go
// In development environment
logger.SetLevel(logger.DebugLevel)
```

This will log detailed information about token validation and authentication flows.
