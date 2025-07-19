package validation

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ValidatorInstance is the global validator instance
var ValidatorInstance *validator.Validate

// init initializes the validator with custom validators
func init() {
	ValidatorInstance = validator.New()
	
	// Register custom validators
	ValidatorInstance.RegisterValidation("latitude", validateLatitude)
	ValidatorInstance.RegisterValidation("longitude", validateLongitude)
	
	// Use JSON tag names in validation errors
	ValidatorInstance.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// validateLatitude validates that a latitude value is within valid range (-90 to 90)
func validateLatitude(fl validator.FieldLevel) bool {
	lat := fl.Field().Float()
	return lat >= -90 && lat <= 90
}

// validateLongitude validates that a longitude value is within valid range (-180 to 180)
func validateLongitude(fl validator.FieldLevel) bool {
	lng := fl.Field().Float()
	return lng >= -180 && lng <= 180
}

// ValidateStruct validates a struct and returns formatted validation errors
func ValidateStruct(s interface{}) []ValidationErrorResponse {
	var errors []ValidationErrorResponse
	
	err := ValidatorInstance.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, ValidationErrorResponse{
				Field:   err.Field(),
				Tag:     err.Tag(),
				Value:   fmt.Sprintf("%v", err.Value()),
				Message: getValidationMessage(err),
			})
		}
	}
	
	return errors
}

// getValidationMessage returns a human-readable validation error message
func getValidationMessage(err validator.FieldError) string {
	field := err.Field()
	tag := err.Tag()
	param := err.Param()
	
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		if err.Kind() == reflect.String {
			return fmt.Sprintf("%s must be at least %s characters long", field, param)
		}
		return fmt.Sprintf("%s must be at least %s", field, param)
	case "max":
		if err.Kind() == reflect.String {
			return fmt.Sprintf("%s must be at most %s characters long", field, param)
		}
		return fmt.Sprintf("%s must be at most %s", field, param)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, param)
	case "alphanum":
		return fmt.Sprintf("%s must contain only alphanumeric characters", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "hexcolor":
		return fmt.Sprintf("%s must be a valid hex color", field)
	case "latitude":
		return fmt.Sprintf("%s must be a valid latitude (-90 to 90)", field)
	case "longitude":
		return fmt.Sprintf("%s must be a valid longitude (-180 to 180)", field)
	case "dive":
		return fmt.Sprintf("%s contains invalid items", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// ValidationMiddleware returns a Gin middleware that handles validation errors
func ValidationMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Next()
		
		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			// Handle validation errors specifically
			if validationErrs, ok := err.Err.(validator.ValidationErrors); ok {
				var errors []ValidationErrorResponse
				for _, verr := range validationErrs {
					errors = append(errors, ValidationErrorResponse{
						Field:   verr.Field(),
						Tag:     verr.Tag(),
						Value:   fmt.Sprintf("%v", verr.Value()),
						Message: getValidationMessage(verr),
					})
				}
				
				SendErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", map[string]interface{}{
					"validation_errors": errors,
				})
				return
			}
			
			// Handle other errors
			SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		}
	})
}

// BindAndValidate binds request data and validates it
func BindAndValidate(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return err
	}
	
	if validationErrors := ValidateStruct(obj); len(validationErrors) > 0 {
		SendErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", map[string]interface{}{
			"validation_errors": validationErrors,
		})
		return fmt.Errorf("validation failed")
	}
	
	return nil
}

// BindQueryAndValidate binds query parameters and validates them
func BindQueryAndValidate(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		return err
	}
	
	if validationErrors := ValidateStruct(obj); len(validationErrors) > 0 {
		SendErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Query validation failed", map[string]interface{}{
			"validation_errors": validationErrors,
		})
		return fmt.Errorf("validation failed")
	}
	
	return nil
}

// SendSuccessResponse sends a successful API response
func SendSuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SendErrorResponse sends an error API response
func SendErrorResponse(c *gin.Context, statusCode int, code, message string, details map[string]interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// PaginatedSuccessResponse sends a paginated successful API response
func PaginatedSuccessResponse(c *gin.Context, data interface{}, pagination PaginationInfo) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: PaginatedResponse{
			Data:       data,
			Pagination: pagination,
		},
	})
}

// CalculatePagination calculates pagination info
func CalculatePagination(page, pageSize int, total int64) PaginationInfo {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	
	return PaginationInfo{
		Page:        page,
		PageSize:    pageSize,
		Total:       total,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}
}

// ValidateTimeWindow validates that end time is after start time
func ValidateTimeWindow(timeWindow *TimeWindowRequest) error {
	if timeWindow == nil {
		return nil
	}
	
	if timeWindow.StartTime != nil && timeWindow.EndTime != nil {
		if timeWindow.EndTime.Before(*timeWindow.StartTime) {
			return fmt.Errorf("end time must be after start time")
		}
	}
	
	return nil
}

// ValidateRouteStops validates route stops for conflicts and consistency
func ValidateRouteStops(stops []RouteStopCreateRequest) error {
	if len(stops) == 0 {
		return nil
	}
	
	// Check for duplicate sequence numbers
	sequenceMap := make(map[int]bool)
	for _, stop := range stops {
		if sequenceMap[stop.SequenceNum] {
			return fmt.Errorf("duplicate sequence number: %d", stop.SequenceNum)
		}
		sequenceMap[stop.SequenceNum] = true
		
		// Validate time window if present
		if err := ValidateTimeWindow(stop.TimeWindow); err != nil {
			return fmt.Errorf("invalid time window for stop %s: %v", stop.Name, err)
		}
	}
	
	return nil
}

// ValidateUpdateRouteStops validates route stops for updates
func ValidateUpdateRouteStops(stops []RouteStopUpdateRequest) error {
	if len(stops) == 0 {
		return nil
	}
	
	// Check for duplicate IDs
	idMap := make(map[uint]bool)
	sequenceMap := make(map[int]bool)
	
	for _, stop := range stops {
		if idMap[stop.ID] {
			return fmt.Errorf("duplicate stop ID: %d", stop.ID)
		}
		idMap[stop.ID] = true
		
		// Check for duplicate sequence numbers if provided
		if stop.SequenceNum != nil {
			if sequenceMap[*stop.SequenceNum] {
				return fmt.Errorf("duplicate sequence number: %d", *stop.SequenceNum)
			}
			sequenceMap[*stop.SequenceNum] = true
		}
		
		// Validate time window if present
		if err := ValidateTimeWindow(stop.TimeWindow); err != nil {
			return fmt.Errorf("invalid time window for stop ID %d: %v", stop.ID, err)
		}
	}
	
	return nil
} 