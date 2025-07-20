package validation

import (
	"routrapp-api/internal/models"
	"time"
)

// StandardResponse represents a standard API response structure
type StandardResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
	Message    string      `json:"message,omitempty"`
}

// MetaData represents pagination metadata
type MetaData struct {
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
	TotalPages  int `json:"total_pages"`
	TotalCount  int `json:"total_count"`
}

// UserResponse represents a user in API responses
type UserResponse struct {
	BaseResponse
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Active    bool      `json:"active"`
	Role      string    `json:"role,omitempty"`
}

// TenantResponse represents a tenant in API responses
type TenantResponse struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name"`
	SubDomain      string    `json:"sub_domain"`
	ContactEmail   string    `json:"contact_email"`
	ContactPhone   string    `json:"contact_phone,omitempty"`
	LogoURL        string    `json:"logo_url,omitempty"`
	PrimaryColor   string    `json:"primary_color,omitempty"`
	SecondaryColor string    `json:"secondary_color,omitempty"`
	Active         bool      `json:"active"`
	PlanType       string    `json:"plan_type"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// RouteResponse represents a route in API responses
type RouteResponse struct {
	BaseResponse
	Name          string               `json:"name"`
	Description   string               `json:"description,omitempty"`
	Status        models.RouteStatus   `json:"status"`
	TechnicianID  *uint                `json:"technician_id,omitempty"`
	Technician    *TechnicianResponse  `json:"technician,omitempty"`
	ScheduledDate *time.Time           `json:"scheduled_date,omitempty"`
	StartedAt     *time.Time           `json:"started_at,omitempty"`
	CompletedAt   *time.Time           `json:"completed_at,omitempty"`
	Notes         string               `json:"notes,omitempty"`
	Stops         []RouteStopResponse  `json:"stops,omitempty"`
}

// RouteStopResponse represents a route stop in API responses
type RouteStopResponse struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Address     string     `json:"address"`
	Lat         float64    `json:"lat"`
	Lng         float64    `json:"lng"`
	SequenceNum int        `json:"sequence_num"`
	StopType    string     `json:"stop_type"`
	Duration    int        `json:"duration"`
	Notes       string     `json:"notes,omitempty"`
	TimeWindow  string     `json:"time_window,omitempty"`
	IsCompleted bool       `json:"is_completed"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TechnicianResponse represents a technician in API responses
type TechnicianResponse struct {
	ID          uint                      `json:"id"`
	User        UserResponse              `json:"user"`
	Status      models.TechnicianStatus   `json:"status"`
	PhoneNumber string                    `json:"phone_number"`
	Notes       string                    `json:"notes,omitempty"`
	LastLat     *float64                  `json:"last_lat,omitempty"`
	LastLng     *float64                  `json:"last_lng,omitempty"`
	LastSeen    *time.Time                `json:"last_seen,omitempty"`
	CreatedAt   time.Time                 `json:"created_at"`
	UpdatedAt   time.Time                 `json:"updated_at"`
}

// LoginResponse represents the response for successful login
type LoginResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    int          `json:"expires_in"`
}

// TokenResponse represents the response for token refresh
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// ValidationErrorResponse represents validation error details
type ValidationErrorResponse struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// PaginationInfo represents pagination metadata
type PaginationInfo struct {
	Page         int   `json:"page"`
	PageSize     int   `json:"page_size"`
	Total        int64 `json:"total"`
	TotalPages   int   `json:"total_pages"`
	HasNext      bool  `json:"has_next"`
	HasPrevious  bool  `json:"has_previous"`
}

// BaseResponse contains common response fields (for backwards compatibility)
type BaseResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// APIResponse represents a standard API response wrapper
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo represents detailed error information
type ErrorInfo struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}