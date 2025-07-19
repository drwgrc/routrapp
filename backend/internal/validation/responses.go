package validation

import (
	"routrapp-api/internal/models"
	"time"
)

// BaseResponse contains common response fields
type BaseResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserResponse represents sanitized user data for API responses
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

// UserLoginResponse represents response after successful login
type UserLoginResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresAt    time.Time    `json:"expires_at"`
}

// TenantResponse represents tenant data for API responses
type TenantResponse struct {
	BaseResponse
	Name           string `json:"name"`
	SubDomain      string `json:"sub_domain"`
	ContactEmail   string `json:"contact_email"`
	ContactPhone   string `json:"contact_phone,omitempty"`
	LogoURL        string `json:"logo_url,omitempty"`
	PrimaryColor   string `json:"primary_color,omitempty"`
	SecondaryColor string `json:"secondary_color,omitempty"`
	Active         bool   `json:"active"`
	PlanType       string `json:"plan_type"`
}

// TechnicianResponse represents technician data for API responses
type TechnicianResponse struct {
	BaseResponse
	UserID         uint                     `json:"user_id"`
	User           *UserResponse            `json:"user,omitempty"`
	TenantID       uint                     `json:"tenant_id"`
	Status         models.TechnicianStatus  `json:"status"`
	PhoneNumber    string                   `json:"phone_number"`
	CurrentLat     *float64                 `json:"current_lat,omitempty"`
	CurrentLng     *float64                 `json:"current_lng,omitempty"`
	LastLocationAt *int64                   `json:"last_location_at,omitempty"`
	Notes          string                   `json:"notes,omitempty"`
}

// RouteStopResponse represents route stop data for API responses
type RouteStopResponse struct {
	BaseResponse
	RouteID      uint                `json:"route_id"`
	TenantID     uint                `json:"tenant_id"`
	Name         string              `json:"name"`
	Address      string              `json:"address"`
	Lat          float64             `json:"lat"`
	Lng          float64             `json:"lng"`
	SequenceNum  int                 `json:"sequence_num"`
	StopType     string              `json:"stop_type"`
	Duration     int                 `json:"duration"`
	Notes        string              `json:"notes,omitempty"`
	TimeWindow   *TimeWindowResponse `json:"time_window,omitempty"`
	IsCompleted  bool                `json:"is_completed"`
	CompletedAt  *time.Time          `json:"completed_at,omitempty"`
	PhotosCount  int                 `json:"photos_count"`
	NotesCount   int                 `json:"notes_count"`
}

// TimeWindowResponse represents time window data for API responses
type TimeWindowResponse struct {
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
}

// RouteResponse represents route data for API responses
type RouteResponse struct {
	BaseResponse
	Name          string                `json:"name"`
	Description   string                `json:"description,omitempty"`
	TenantID      uint                  `json:"tenant_id"`
	TechnicianID  *uint                 `json:"technician_id,omitempty"`
	Technician    *TechnicianResponse   `json:"technician,omitempty"`
	Status        models.RouteStatus    `json:"status"`
	ScheduledDate *time.Time            `json:"scheduled_date,omitempty"`
	StartedAt     *time.Time            `json:"started_at,omitempty"`
	CompletedAt   *time.Time            `json:"completed_at,omitempty"`
	CancelledAt   *time.Time            `json:"cancelled_at,omitempty"`
	Stops         []RouteStopResponse   `json:"stops,omitempty"`
	IsOptimized   bool                  `json:"is_optimized"`
	TotalDistance float64               `json:"total_distance"`
	TotalDuration int                   `json:"total_duration"`
	Notes         string                `json:"notes,omitempty"`
}

// RouteActivityResponse represents route activity data for API responses
type RouteActivityResponse struct {
	BaseResponse
	RouteID      uint                `json:"route_id"`
	RouteStopID  *uint               `json:"route_stop_id,omitempty"`
	TenantID     uint                `json:"tenant_id"`
	TechnicianID uint                `json:"technician_id"`
	Technician   *TechnicianResponse `json:"technician,omitempty"`
	ActivityType string              `json:"activity_type"`
	Notes        string              `json:"notes,omitempty"`
	Lat          *float64            `json:"lat,omitempty"`
	Lng          *float64            `json:"lng,omitempty"`
	PhotoURL     string              `json:"photo_url,omitempty"`
	Timestamp    time.Time           `json:"timestamp"`
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

// ValidationErrorResponse represents validation error details
type ValidationErrorResponse struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// PaginatedResponse represents paginated API response
type PaginatedResponse struct {
	Data       interface{}      `json:"data"`
	Pagination PaginationInfo   `json:"pagination"`
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

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp string                 `json:"timestamp"`
	Database  string                 `json:"database,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

// SuccessResponse represents a simple success response
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// StatsResponse represents various statistics
type StatsResponse struct {
	TotalRoutes       int64   `json:"total_routes"`
	CompletedRoutes   int64   `json:"completed_routes"`
	PendingRoutes     int64   `json:"pending_routes"`
	ActiveTechnicians int64   `json:"active_technicians"`
	TotalDistance     float64 `json:"total_distance"`
	AverageRouteTime  float64 `json:"average_route_time"`
} 