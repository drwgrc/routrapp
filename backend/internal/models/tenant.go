package models

// Tenant represents a tenant in the multi-tenant system
type Tenant struct {
	Base
	Name           string `gorm:"type:varchar(100)" json:"name"`
	SubDomain      string `gorm:"type:varchar(100);uniqueIndex" json:"sub_domain"`
	ContactEmail   string `gorm:"type:varchar(100)" json:"contact_email"`
	ContactPhone   string `gorm:"type:varchar(20)" json:"contact_phone,omitempty"`
	LogoURL        string `gorm:"type:varchar(255)" json:"logo_url,omitempty"`
	PrimaryColor   string `gorm:"type:varchar(20)" json:"primary_color,omitempty"`
	SecondaryColor string `gorm:"type:varchar(20)" json:"secondary_color,omitempty"`
	Active         bool   `gorm:"default:true" json:"active"`
	PlanType       string `gorm:"type:varchar(20);default:'basic'" json:"plan_type"`
} 