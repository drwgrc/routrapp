package models

import "time"

// User represents a user in the system
type User struct {
	Base
	Email        string     `gorm:"type:varchar(100)" json:"email"`
	Password     string     `gorm:"type:varchar(255)" json:"-"`
	FirstName    string     `gorm:"type:varchar(100)" json:"first_name"`
	LastName     string     `gorm:"type:varchar(100)" json:"last_name"`
	RoleID       uint       `gorm:"not null;index" json:"role_id"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	RefreshToken string     `gorm:"type:varchar(1000)" json:"-"`
	Active       bool       `gorm:"default:true" json:"active"`

	// Relationships
	Role        Role        `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Technician  *Technician `gorm:"foreignKey:UserID" json:"technician,omitempty"`
}

// TableName returns the table name for User
func (User) TableName() string {
	return "users"
}

// Indexes returns the database indexes for the User model
func (User) Indexes() []string {
	return []string{
		"CREATE UNIQUE INDEX IF NOT EXISTS idx_users_org_email ON users(organization_id, email) WHERE deleted_at IS NULL",
	}
}

// GetFullName returns the user's full name
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

// HasPermission checks if the user has a specific permission through their role
func (u *User) HasPermission(permission string) bool {
	if u.Role.Name == "" {
		return false // Role not loaded
	}
	return u.Role.HasPermission(permission)
}

// IsOwner checks if the user has owner role
func (u *User) IsOwner() bool {
	if u.Role.Name == "" {
		return false // Role not loaded
	}
	return u.Role.Name == RoleTypeOwner
}

// IsTechnician checks if the user has technician role
func (u *User) IsTechnician() bool {
	if u.Role.Name == "" {
		return false // Role not loaded
	}
	return u.Role.Name == RoleTypeTechnician
} 