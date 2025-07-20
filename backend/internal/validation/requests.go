package validation

import (
	"routrapp-api/internal/models"
	"time"
)

// UserRegistrationRequest represents request for user registration
type UserRegistrationRequest struct {
	Email     string           `json:"email" binding:"required,email,max=100"`
	Password  string           `json:"password" binding:"required,min=8,max=255"`
	FirstName string           `json:"first_name" binding:"required,min=1,max=100"`
	LastName  string           `json:"last_name" binding:"required,min=1,max=100"`
	Role      models.RoleType  `json:"role" binding:"required,oneof=owner technician"`
	TenantID  uint             `json:"tenant_id" binding:"required,min=1"`
}

// UserLoginRequest represents request for user login
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserUpdateRequest represents request for updating user profile
type UserUpdateRequest struct {
	FirstName *string           `json:"first_name,omitempty" binding:"omitempty,min=1,max=100"`
	LastName  *string           `json:"last_name,omitempty" binding:"omitempty,min=1,max=100"`
	Role      *models.RoleType  `json:"role,omitempty" binding:"omitempty,oneof=owner technician"`
	Active    *bool             `json:"active,omitempty"`
}

// ChangePasswordRequest represents request for changing user password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8,max=255"`
}

// TenantCreateRequest represents request for creating a new tenant
type TenantCreateRequest struct {
	Name           string `json:"name" binding:"required,min=1,max=100"`
	SubDomain      string `json:"sub_domain" binding:"required,min=1,max=100,alphanum"`
	ContactEmail   string `json:"contact_email" binding:"required,email,max=100"`
	ContactPhone   string `json:"contact_phone,omitempty" binding:"omitempty,max=20"`
	LogoURL        string `json:"logo_url,omitempty" binding:"omitempty,url,max=255"`
	PrimaryColor   string `json:"primary_color,omitempty" binding:"omitempty,hexcolor,max=20"`
	SecondaryColor string `json:"secondary_color,omitempty" binding:"omitempty,hexcolor,max=20"`
	PlanType       string `json:"plan_type,omitempty" binding:"omitempty,oneof=basic premium enterprise"`
}

// TenantUpdateRequest represents request for updating tenant
type TenantUpdateRequest struct {
	Name           *string `json:"name,omitempty" binding:"omitempty,min=1,max=100"`
	ContactEmail   *string `json:"contact_email,omitempty" binding:"omitempty,email,max=100"`
	ContactPhone   *string `json:"contact_phone,omitempty" binding:"omitempty,max=20"`
	LogoURL        *string `json:"logo_url,omitempty" binding:"omitempty,url,max=255"`
	PrimaryColor   *string `json:"primary_color,omitempty" binding:"omitempty,hexcolor,max=20"`
	SecondaryColor *string `json:"secondary_color,omitempty" binding:"omitempty,hexcolor,max=20"`
	Active         *bool   `json:"active,omitempty"`
	PlanType       *string `json:"plan_type,omitempty" binding:"omitempty,oneof=basic premium enterprise"`
}

// RouteCreateRequest represents request for creating a new route
type RouteCreateRequest struct {
	Name          string                     `json:"name" binding:"required,min=1,max=100"`
	Description   string                     `json:"description,omitempty" binding:"omitempty,max=1000"`
	TechnicianID  *uint                      `json:"technician_id,omitempty" binding:"omitempty,min=1"`
	ScheduledDate *time.Time                 `json:"scheduled_date,omitempty"`
	Notes         string                     `json:"notes,omitempty" binding:"omitempty,max=1000"`
	Stops         []RouteStopCreateRequest   `json:"stops,omitempty" binding:"omitempty,dive"`
}

// RouteUpdateRequest represents request for updating a route
type RouteUpdateRequest struct {
	Name          *string                    `json:"name,omitempty" binding:"omitempty,min=1,max=100"`
	Description   *string                    `json:"description,omitempty" binding:"omitempty,max=1000"`
	TechnicianID  *uint                      `json:"technician_id,omitempty" binding:"omitempty,min=1"`
	Status        *models.RouteStatus        `json:"status,omitempty" binding:"omitempty,oneof=pending assigned started completed cancelled paused"`
	ScheduledDate *time.Time                 `json:"scheduled_date,omitempty"`
	Notes         *string                    `json:"notes,omitempty" binding:"omitempty,max=1000"`
	Stops         []RouteStopUpdateRequest   `json:"stops,omitempty" binding:"omitempty,dive"`
}

// RouteStopCreateRequest represents request for creating a route stop
type RouteStopCreateRequest struct {
	Name        string               `json:"name" binding:"required,min=1,max=100"`
	Address     string               `json:"address" binding:"required,min=1,max=255"`
	Lat         float64              `json:"lat" binding:"required,latitude"`
	Lng         float64              `json:"lng" binding:"required,longitude"`
	SequenceNum int                  `json:"sequence_num" binding:"required,min=1"`
	StopType    string               `json:"stop_type" binding:"required,oneof=pickup delivery service maintenance"`
	Duration    int                  `json:"duration" binding:"required,min=1,max=1440"` // max 24 hours in minutes
	Notes       string               `json:"notes,omitempty" binding:"omitempty,max=1000"`
	TimeWindow  *TimeWindowRequest   `json:"time_window,omitempty"`
}

// RouteStopUpdateRequest represents request for updating a route stop
type RouteStopUpdateRequest struct {
	ID          uint                 `json:"id" binding:"required,min=1"`
	Name        *string              `json:"name,omitempty" binding:"omitempty,min=1,max=100"`
	Address     *string              `json:"address,omitempty" binding:"omitempty,min=1,max=255"`
	Lat         *float64             `json:"lat,omitempty" binding:"omitempty,latitude"`
	Lng         *float64             `json:"lng,omitempty" binding:"omitempty,longitude"`
	SequenceNum *int                 `json:"sequence_num,omitempty" binding:"omitempty,min=1"`
	StopType    *string              `json:"stop_type,omitempty" binding:"omitempty,oneof=pickup delivery service maintenance"`
	Duration    *int                 `json:"duration,omitempty" binding:"omitempty,min=1,max=1440"`
	Notes       *string              `json:"notes,omitempty" binding:"omitempty,max=1000"`
	TimeWindow  *TimeWindowRequest   `json:"time_window,omitempty"`
	IsCompleted *bool                `json:"is_completed,omitempty"`
}

// TimeWindowRequest represents time window validation
type TimeWindowRequest struct {
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
}

// TechnicianCreateRequest represents request for creating a technician
type TechnicianCreateRequest struct {
	UserID      uint   `json:"user_id" binding:"required,min=1"`
	PhoneNumber string `json:"phone_number" binding:"required,min=10,max=20"`
	Notes       string `json:"notes,omitempty" binding:"omitempty,max=1000"`
}

// TechnicianUpdateRequest represents request for updating a technician
type TechnicianUpdateRequest struct {
	Status      *models.TechnicianStatus `json:"status,omitempty" binding:"omitempty,oneof=active inactive on_route on_break off_duty"`
	PhoneNumber *string                  `json:"phone_number,omitempty" binding:"omitempty,min=10,max=20"`
	Notes       *string                  `json:"notes,omitempty" binding:"omitempty,max=1000"`
}

// LocationUpdateRequest represents request for updating technician location
type LocationUpdateRequest struct {
	Lat float64 `json:"lat" binding:"required,latitude"`
	Lng float64 `json:"lng" binding:"required,longitude"`
}

// RouteActivityCreateRequest represents request for creating route activity
type RouteActivityCreateRequest struct {
	RouteStopID  *uint   `json:"route_stop_id,omitempty" binding:"omitempty,min=1"`
	ActivityType string  `json:"activity_type" binding:"required,oneof=start stop complete pause resume note photo"`
	Notes        string  `json:"notes,omitempty" binding:"omitempty,max=1000"`
	Lat          *float64 `json:"lat,omitempty" binding:"omitempty,latitude"`
	Lng          *float64 `json:"lng,omitempty" binding:"omitempty,longitude"`
	PhotoURL     string  `json:"photo_url,omitempty" binding:"omitempty,url,max=255"`
}

// PaginationRequest represents common pagination parameters
type PaginationRequest struct {
	Page     int `form:"page,default=1" binding:"min=1"`
	PageSize int `form:"page_size,default=20" binding:"min=1,max=100"`
}

// FilterRequest represents common filtering parameters
type FilterRequest struct {
	Search   string `form:"search,omitempty" binding:"omitempty,max=100"`
	SortBy   string `form:"sort_by,omitempty" binding:"omitempty,oneof=id name created_at updated_at"`
	SortDesc bool   `form:"sort_desc,omitempty"`
}

// RouteFilterRequest represents route-specific filtering
type RouteFilterRequest struct {
	FilterRequest
	Status       []models.RouteStatus `form:"status,omitempty" binding:"omitempty,dive,oneof=pending assigned started completed cancelled paused"`
	TechnicianID *uint                `form:"technician_id,omitempty" binding:"omitempty,min=1"`
	DateFrom     *time.Time           `form:"date_from,omitempty"`
	DateTo       *time.Time           `form:"date_to,omitempty"`
}

// TechnicianFilterRequest represents technician-specific filtering  
type TechnicianFilterRequest struct {
	FilterRequest
	Status []models.TechnicianStatus `form:"status,omitempty" binding:"omitempty,dive,oneof=active inactive on_route on_break off_duty"`
	Active *bool                     `form:"active,omitempty"`
}

// RefreshTokenRequest represents request for refreshing access token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RegistrationRequest represents request for user registration with organization creation
type RegistrationRequest struct {
	// User details
	Email     string `json:"email" binding:"required,email,max=100"`
	Password  string `json:"password" binding:"required,min=8,max=255"`
	FirstName string `json:"first_name" binding:"required,min=1,max=100"`
	LastName  string `json:"last_name" binding:"required,min=1,max=100"`
	
	// Organization details
	OrganizationName  string `json:"organization_name" binding:"required,min=1,max=100"`
	OrganizationEmail string `json:"organization_email" binding:"required,email,max=100"`
	SubDomain         string `json:"sub_domain" binding:"required,min=1,max=100,alphanum"`
}
