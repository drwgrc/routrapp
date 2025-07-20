package models

// TechnicianStatus represents the current status of a technician
type TechnicianStatus string

// Technician status constants
const (
	TechnicianStatusActive    TechnicianStatus = "active"
	TechnicianStatusInactive  TechnicianStatus = "inactive"
	TechnicianStatusOnRoute   TechnicianStatus = "on_route"
	TechnicianStatusOnBreak   TechnicianStatus = "on_break"
	TechnicianStatusOffDuty   TechnicianStatus = "off_duty"
)

// Technician represents a technician in the system
type Technician struct {
	Base
	UserID         uint             `gorm:"uniqueIndex" json:"user_id"`
	User           User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Status         TechnicianStatus `gorm:"type:varchar(20);default:'inactive'" json:"status"`
	PhoneNumber    string           `gorm:"type:varchar(20)" json:"phone_number"`
	CurrentLat     *float64         `json:"current_lat,omitempty"`
	CurrentLng     *float64         `json:"current_lng,omitempty"`
	LastLocationAt *int64           `json:"last_location_at,omitempty"`
	Notes          string           `gorm:"type:text" json:"notes,omitempty"`
} 