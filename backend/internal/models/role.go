package models

import (
	"database/sql/driver"
	"fmt"
)

// RoleType represents the role type enumeration
type RoleType string

// Role type constants
const (
	RoleTypeOwner      RoleType = "owner"
	RoleTypeTechnician RoleType = "technician"
)

// String returns the string representation of RoleType
func (r RoleType) String() string {
	return string(r)
}

// Value implements the driver.Valuer interface for database storage
func (r RoleType) Value() (driver.Value, error) {
	return string(r), nil
}

// Scan implements the sql.Scanner interface for database retrieval
func (r *RoleType) Scan(value interface{}) error {
	if value == nil {
		*r = RoleTypeTechnician
		return nil
	}
	
	switch v := value.(type) {
	case string:
		*r = RoleType(v)
	case []byte:
		*r = RoleType(v)
	default:
		return fmt.Errorf("cannot scan %T into RoleType", value)
	}
	return nil
}

// IsValid checks if the role type is valid
func (r RoleType) IsValid() bool {
	switch r {
	case RoleTypeOwner, RoleTypeTechnician:
		return true
	default:
		return false
	}
}

// Role represents a user role with specific permissions
type Role struct {
	Base
	Name        RoleType `gorm:"type:varchar(20);not null;uniqueIndex:idx_role_org_name" json:"name"`
	DisplayName string   `gorm:"type:varchar(100);not null" json:"display_name"`
	Description string   `gorm:"type:text" json:"description,omitempty"`
	Permissions string   `gorm:"type:text" json:"permissions,omitempty"` // JSON array of permissions
	Active      bool     `gorm:"default:true" json:"active"`

	// Relationships
	Users []User `gorm:"foreignKey:RoleID" json:"users,omitempty"`
}

// TableName returns the table name for Role
func (Role) TableName() string {
	return "roles"
}

// HasPermission checks if the role has a specific permission
func (r *Role) HasPermission(permission string) bool {
	// This is a simplified implementation
	// In a real system, you'd parse the JSON permissions field
	switch r.Name {
	case RoleTypeOwner:
		return true // Owner has all permissions
	case RoleTypeTechnician:
		// Technicians have limited permissions
		techPermissions := []string{
			"routes.read",
			"routes.update_status",
			"technicians.read_own",
			"technicians.update_own",
		}
		for _, perm := range techPermissions {
			if perm == permission {
				return true
			}
		}
	}
	return false
}

// GetDefaultPermissions returns default permissions for a role type
func GetDefaultPermissions(roleType RoleType) []string {
	switch roleType {
	case RoleTypeOwner:
		return []string{
			"organizations.*",
			"users.*",
			"technicians.*",
			"routes.*",
			"roles.*",
		}
	case RoleTypeTechnician:
		return []string{
			"routes.read",
			"routes.update_status",
			"technicians.read_own",
			"technicians.update_own",
		}
	default:
		return []string{}
	}
} 