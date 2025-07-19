package validation

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Example handlers demonstrating how to use validation structs

// ExampleUserRegistrationHandler demonstrates user registration with validation
func ExampleUserRegistrationHandler(c *gin.Context) {
	var req UserRegistrationRequest
	
	// Bind and validate request
	if err := BindAndValidate(c, &req); err != nil {
		return // Error response already sent by BindAndValidate
	}
	
	// Additional custom validation if needed
	if req.Password == req.Email {
		SendErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Password cannot be the same as email", nil)
		return
	}
	
	// Process the validated request...
	// userService.Register(req)
	
	SendSuccessResponse(c, http.StatusCreated, "User registered successfully", map[string]interface{}{
		"user_id": 123, // This would come from your service
	})
}

// ExampleRouteCreationHandler demonstrates route creation with validation
func ExampleRouteCreationHandler(c *gin.Context) {
	var req RouteCreateRequest
	
	// Bind and validate request
	if err := BindAndValidate(c, &req); err != nil {
		return
	}
	
	// Additional validation for route stops
	if err := ValidateRouteStops(req.Stops); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}
	
	// Process the validated request...
	// route := routeService.Create(req)
	
	response := RouteResponse{
		BaseResponse: BaseResponse{
			ID: 1, // This would come from your service
		},
		Name:        req.Name,
		Description: req.Description,
		// ... other fields
	}
	
	SendSuccessResponse(c, http.StatusCreated, "Route created successfully", response)
}

// ExampleGetRoutesHandler demonstrates paginated listing with filters
func ExampleGetRoutesHandler(c *gin.Context) {
	var pagination PaginationRequest
	var filters RouteFilterRequest
	
	// Bind and validate pagination
	if err := BindQueryAndValidate(c, &pagination); err != nil {
		return
	}
	
	// Bind and validate filters
	if err := BindQueryAndValidate(c, &filters); err != nil {
		return
	}
	
	// Process the validated request...
	// routes, total := routeService.List(pagination, filters)
	
	// Example response data
	routes := []RouteResponse{
		{
			BaseResponse: BaseResponse{ID: 1},
			Name:         "Route 1",
			Status:       "pending",
		},
	}
	total := int64(1)
	
	paginationInfo := CalculatePagination(pagination.Page, pagination.PageSize, total)
	PaginatedSuccessResponse(c, routes, paginationInfo)
}

// ExampleLocationUpdateHandler demonstrates location update with validation
func ExampleLocationUpdateHandler(c *gin.Context) {
	// Get technician ID from URL parameter
	technicianIDStr := c.Param("id")
	technicianID, err := strconv.ParseUint(technicianIDStr, 10, 32)
	if err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid technician ID", nil)
		return
	}
	
	var req LocationUpdateRequest
	
	// Bind and validate request
	if err := BindAndValidate(c, &req); err != nil {
		return
	}
	
	// Process the validated request...
	// technicianService.UpdateLocation(uint(technicianID), req)
	
	SendSuccessResponse(c, http.StatusOK, "Location updated successfully", map[string]interface{}{
		"technician_id": technicianID,
		"lat":           req.Lat,
		"lng":           req.Lng,
		"updated_at":    "2024-01-01T12:00:00Z", // This would come from your service
	})
}

// ExampleUserUpdateHandler demonstrates partial updates with validation
func ExampleUserUpdateHandler(c *gin.Context) {
	// Get user ID from URL parameter
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid user ID", nil)
		return
	}
	
	var req UserUpdateRequest
	
	// Bind and validate request
	if err := BindAndValidate(c, &req); err != nil {
		return
	}
	
	// Check if at least one field is provided for update
	if req.FirstName == nil && req.LastName == nil && req.Role == nil && req.Active == nil {
		SendErrorResponse(c, http.StatusBadRequest, "NO_FIELDS", "At least one field must be provided for update", nil)
		return
	}
	
	// Process the validated request...
	// user := userService.Update(uint(userID), req)
	
	response := UserResponse{
		BaseResponse: BaseResponse{ID: uint(userID)},
		Email:        "user@example.com", // This would come from your service
		FirstName:    "John",
		LastName:     "Doe",
		// ... other fields
	}
	
	SendSuccessResponse(c, http.StatusOK, "User updated successfully", response)
}

// ExampleErrorHandling demonstrates various error scenarios
func ExampleErrorHandling(c *gin.Context) {
	var req UserRegistrationRequest
	
	// This will demonstrate validation errors
	if err := c.ShouldBindJSON(&req); err != nil {
		SendErrorResponse(c, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON format", map[string]interface{}{
			"details": err.Error(),
		})
		return
	}
	
	// Validate the struct
	if validationErrors := ValidateStruct(&req); len(validationErrors) > 0 {
		SendErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", map[string]interface{}{
			"validation_errors": validationErrors,
		})
		return
	}
	
	// Business logic validation
	if req.Email == "admin@example.com" {
		SendErrorResponse(c, http.StatusConflict, "EMAIL_TAKEN", "Email address is already in use", nil)
		return
	}
	
	// Simulate a server error
	SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong", map[string]interface{}{
		"request_id": "req_123456",
	})
} 