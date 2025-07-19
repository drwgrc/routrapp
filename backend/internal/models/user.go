package models

import "time"

// Role represents user's role in the system
type Role string

// User role constants
const (
	RoleOwner     Role = "owner"
	RoleTechnician Role = "technician"
)

// User represents a user in the system
type User struct {
	Base
	Email        string     `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Password     string     `gorm:"type:varchar(255)" json:"-"`
	FirstName    string     `gorm:"type:varchar(100)" json:"first_name"`
	LastName     string     `gorm:"type:varchar(100)" json:"last_name"`
	Role         Role       `gorm:"type:varchar(20);default:'technician'" json:"role"`
	TenantID     uint       `gorm:"index" json:"tenant_id"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	RefreshToken string     `gorm:"type:varchar(255)" json:"-"`
	Active       bool       `gorm:"default:true" json:"active"`
} 