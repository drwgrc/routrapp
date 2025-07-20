# Validation Package

This package provides comprehensive request/response validation for the RoutrApp Go backend using the `github.com/go-playground/validator/v10` library.

## Overview

The validation package includes:

- **Request structs** with validation tags for all API endpoints
- **Response structs** with proper JSON serialization
- **Helper functions** for binding, validating, and sending responses
- **Custom validators** for domain-specific validation (lat/lng, time windows, etc.)
- **Middleware** for handling validation errors consistently

## Files Structure

```
validation/
├── requests.go     # Request validation structs
├── responses.go    # Response structs with JSON tags
├── helpers.go      # Validation helpers and response functions
├── examples.go     # Example handler implementations
└── README.md       # This documentation
```

## Quick Start

### 1. Basic Handler with Validation

```go
func CreateUserHandler(c *gin.Context) {
    var req validation.UserRegistrationRequest

    // Bind and validate in one step
    if err := validation.BindAndValidate(c, &req); err != nil {
        return // Error response already sent
    }

    // Process validated request
    user := userService.Create(req)

    // Send standardized response
    validation.SendSuccessResponse(c, http.StatusCreated, "User created", user)
}
```

### 2. Paginated List with Filters

```go
func ListRoutesHandler(c *gin.Context) {
    var pagination validation.PaginationRequest
    var filters validation.RouteFilterRequest

    // Validate query parameters
    if err := validation.BindQueryAndValidate(c, &pagination); err != nil {
        return
    }
    if err := validation.BindQueryAndValidate(c, &filters); err != nil {
        return
    }

    routes, total := routeService.List(pagination, filters)
    paginationInfo := validation.CalculatePagination(pagination.Page, pagination.PageSize, total)

    validation.PaginatedSuccessResponse(c, routes, paginationInfo)
}
```

## Request Validation Structs

### User Operations

```go
// User registration
type UserRegistrationRequest struct {
    Email     string      `json:"email" binding:"required,email,max=100"`
    Password  string      `json:"password" binding:"required,min=8,max=255"`
    FirstName string      `json:"first_name" binding:"required,min=1,max=100"`
    LastName  string      `json:"last_name" binding:"required,min=1,max=100"`
    Role      models.Role `json:"role" binding:"required,oneof=owner technician"`
    TenantID  uint        `json:"tenant_id" binding:"required,min=1"`
}

// User login
type UserLoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

// User updates (all fields optional)
type UserUpdateRequest struct {
    FirstName *string      `json:"first_name,omitempty" binding:"omitempty,min=1,max=100"`
    LastName  *string      `json:"last_name,omitempty" binding:"omitempty,min=1,max=100"`
    Role      *models.Role `json:"role,omitempty" binding:"omitempty,oneof=owner technician"`
    Active    *bool        `json:"active,omitempty"`
}
```

### Route Operations

```go
// Route creation
type RouteCreateRequest struct {
    Name          string                    `json:"name" binding:"required,min=1,max=100"`
    Description   string                    `json:"description,omitempty" binding:"omitempty,max=1000"`
    TechnicianID  *uint                     `json:"technician_id,omitempty" binding:"omitempty,min=1"`
    ScheduledDate *time.Time                `json:"scheduled_date,omitempty"`
    Notes         string                    `json:"notes,omitempty" binding:"omitempty,max=1000"`
    Stops         []RouteStopCreateRequest  `json:"stops,omitempty" binding:"omitempty,dive"`
}

// Route stop creation
type RouteStopCreateRequest struct {
    Name        string               `json:"name" binding:"required,min=1,max=100"`
    Address     string               `json:"address" binding:"required,min=1,max=255"`
    Lat         float64              `json:"lat" binding:"required,latitude"`
    Lng         float64              `json:"lng" binding:"required,longitude"`
    SequenceNum int                  `json:"sequence_num" binding:"required,min=1"`
    StopType    string               `json:"stop_type" binding:"required,oneof=pickup delivery service maintenance"`
    Duration    int                  `json:"duration" binding:"required,min=1,max=1440"`
    Notes       string               `json:"notes,omitempty" binding:"omitempty,max=1000"`
    TimeWindow  *TimeWindowRequest   `json:"time_window,omitempty"`
}
```

### Location & Activity

```go
// Location updates
type LocationUpdateRequest struct {
    Lat float64 `json:"lat" binding:"required,latitude"`
    Lng float64 `json:"lng" binding:"required,longitude"`
}

// Route activities
type RouteActivityCreateRequest struct {
    RouteStopID  *uint   `json:"route_stop_id,omitempty" binding:"omitempty,min=1"`
    ActivityType string  `json:"activity_type" binding:"required,oneof=start stop complete pause resume note photo"`
    Notes        string  `json:"notes,omitempty" binding:"omitempty,max=1000"`
    Lat          *float64 `json:"lat,omitempty" binding:"omitempty,latitude"`
    Lng          *float64 `json:"lng,omitempty" binding:"omitempty,longitude"`
    PhotoURL     string  `json:"photo_url,omitempty" binding:"omitempty,url,max=255"`
}
```

## Response Structs

All response structs include proper JSON tags and omitempty for optional fields:

```go
// Standard API response wrapper
type APIResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Error   *ErrorInfo  `json:"error,omitempty"`
}

// User response (passwords excluded)
type UserResponse struct {
    BaseResponse
    Email       string      `json:"email"`
    FirstName   string      `json:"first_name"`
    LastName    string      `json:"last_name"`
    Role        models.Role `json:"role"`
    TenantID    uint        `json:"tenant_id"`
    LastLoginAt *time.Time  `json:"last_login_at,omitempty"`
    Active      bool        `json:"active"`
}

// Paginated responses
type PaginatedResponse struct {
    Data       interface{}      `json:"data"`
    Pagination PaginationInfo   `json:"pagination"`
}
```

## Validation Tags Reference

### Basic Validation

- `required` - Field is required
- `omitempty` - Skip validation if field is empty
- `min=N` - Minimum value/length
- `max=N` - Maximum value/length
- `email` - Valid email format
- `url` - Valid URL format
- `alphanum` - Only alphanumeric characters
- `hexcolor` - Valid hex color code

### Custom Validators

- `latitude` - Valid latitude (-90 to 90)
- `longitude` - Valid longitude (-180 to 180)
- `oneof=val1 val2` - Value must be one of the specified options
- `dive` - Validate each element in slice/array

### Examples

```go
Email     string  `binding:"required,email,max=100"`        // Required valid email, max 100 chars
Password  string  `binding:"required,min=8,max=255"`        // Required, 8-255 chars
Role      string  `binding:"required,oneof=owner tech"`     // Must be "owner" or "tech"
Lat       float64 `binding:"required,latitude"`             // Required valid latitude
Optional  *string `binding:"omitempty,min=1,max=50"`        // Optional, but if provided must be 1-50 chars
```

## Helper Functions

### Binding and Validation

```go
// Bind JSON and validate struct
func BindAndValidate(c *gin.Context, obj interface{}) error

// Bind query params and validate
func BindQueryAndValidate(c *gin.Context, obj interface{}) error

// Validate struct manually
func ValidateStruct(s interface{}) []ValidationErrorResponse
```

### Response Functions

```go
// Success responses
func SendSuccessResponse(c *gin.Context, statusCode int, message string, data interface{})
func PaginatedSuccessResponse(c *gin.Context, data interface{}, pagination PaginationInfo)

// Error responses
func SendErrorResponse(c *gin.Context, statusCode int, code, message string, details map[string]interface{})

// Pagination helpers
func CalculatePagination(page, pageSize int, total int64) PaginationInfo
```

### Domain-Specific Validation

```go
// Validate time windows
func ValidateTimeWindow(timeWindow *TimeWindowRequest) error

// Validate route stops for conflicts
func ValidateRouteStops(stops []RouteStopCreateRequest) error
func ValidateUpdateRouteStops(stops []RouteStopUpdateRequest) error
```

## Error Responses

### Validation Errors

When validation fails, the API returns a structured error response:

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": {
      "validation_errors": [
        {
          "field": "email",
          "tag": "required",
          "value": "",
          "message": "email is required"
        },
        {
          "field": "password",
          "tag": "min",
          "value": "123",
          "message": "password must be at least 8 characters long"
        }
      ]
    }
  }
}
```

### Business Logic Errors

```json
{
  "success": false,
  "error": {
    "code": "EMAIL_TAKEN",
    "message": "Email address is already in use",
    "details": null
  }
}
```

## Integration with Gin

### Add Validation Middleware

```go
// In your router setup
router.Use(validation.ValidationMiddleware())
```

### Example Route Handler

```go
func setupRoutes(router *gin.Engine) {
    api := router.Group("/api/v1")

    // User routes
    api.POST("/users", CreateUserHandler)
    api.PUT("/users/:id", UpdateUserHandler)
    api.GET("/users", ListUsersHandler)

    // Route routes
    api.POST("/routes", CreateRouteHandler)
    api.GET("/routes", ListRoutesHandler)
    api.PUT("/routes/:id", UpdateRouteHandler)

    // Location updates
    api.PUT("/technicians/:id/location", UpdateLocationHandler)
}
```

## Testing Validation

```go
func TestUserValidation(t *testing.T) {
    // Test valid request
    validReq := validation.UserRegistrationRequest{
        Email:     "test@example.com",
        Password:  "password123",
        FirstName: "John",
        LastName:  "Doe",
        Role:      "technician",
        TenantID:  1,
    }

    errors := validation.ValidateStruct(&validReq)
    assert.Empty(t, errors)

    // Test invalid request
    invalidReq := validation.UserRegistrationRequest{
        Email:    "invalid-email",
        Password: "123", // too short
        // missing required fields
    }

    errors = validation.ValidateStruct(&invalidReq)
    assert.NotEmpty(t, errors)
    assert.Len(t, errors, 4) // email, password, firstName, lastName
}
```

## Best Practices

1. **Always use BindAndValidate()** for request handling - it handles both binding and validation in one step

2. **Use pointer fields for optional updates** - allows distinguishing between "not provided" and "set to zero value"

3. **Add custom validation for business rules** - after struct validation, add domain-specific checks

4. **Provide meaningful error messages** - the helper functions generate user-friendly error messages

5. **Use consistent response format** - always use the provided response helper functions

6. **Validate nested structs** - use `dive` tag for slice/array validation

7. **Test validation thoroughly** - write unit tests for both valid and invalid cases

## Common Patterns

### Partial Updates

```go
type UserUpdateRequest struct {
    FirstName *string `json:"first_name,omitempty" binding:"omitempty,min=1"`
    LastName  *string `json:"last_name,omitempty" binding:"omitempty,min=1"`
}

// Check if any field was provided
if req.FirstName == nil && req.LastName == nil {
    validation.SendErrorResponse(c, 400, "NO_FIELDS", "At least one field required", nil)
    return
}
```

### ID Parameter Validation

```go
idStr := c.Param("id")
id, err := strconv.ParseUint(idStr, 10, 32)
if err != nil {
    validation.SendErrorResponse(c, 400, "INVALID_ID", "Invalid ID format", nil)
    return
}
```

### Complex Business Validation

```go
// After struct validation
if req.StartDate.After(req.EndDate) {
    validation.SendErrorResponse(c, 400, "INVALID_DATE_RANGE", "Start date must be before end date", nil)
    return
}
```
