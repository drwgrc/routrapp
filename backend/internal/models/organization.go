package models

import (
	"time"

	"gorm.io/gorm"
)

// Organization represents a tenant organization in the system
// Note: Organizations don't embed Base since they are the root multi-tenant entity
type Organization struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"type:varchar(100);not null" json:"name"`
	SubDomain      string         `gorm:"type:varchar(100);not null;uniqueIndex" json:"sub_domain"`
	ContactEmail   string         `gorm:"type:varchar(100);not null" json:"contact_email"`
	ContactPhone   string         `gorm:"type:varchar(20)" json:"contact_phone,omitempty"`
	LogoURL        string         `gorm:"type:varchar(255)" json:"logo_url,omitempty"`
	PrimaryColor   string         `gorm:"type:varchar(20)" json:"primary_color,omitempty"`
	SecondaryColor string         `gorm:"type:varchar(20)" json:"secondary_color,omitempty"`
	Active         bool           `gorm:"default:true" json:"active"`
	PlanType       string         `gorm:"type:varchar(20);default:'basic'" json:"plan_type"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Users       []User       `gorm:"foreignKey:OrganizationID" json:"users,omitempty"`
	Technicians []Technician `gorm:"foreignKey:OrganizationID" json:"technicians,omitempty"`
	Routes      []Route      `gorm:"foreignKey:OrganizationID" json:"routes,omitempty"`
}

// TableName returns the table name for Organization
func (Organization) TableName() string {
	return "organizations"
} 