package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
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
	
	var roleStr string
	switch v := value.(type) {
	case string:
		roleStr = v
	case []byte:
		roleStr = string(v)
	default:
		return fmt.Errorf("cannot scan %T into RoleType", value)
	}
	
	// Convert to RoleType and validate
	roleType := RoleType(roleStr)
	if !roleType.IsValid() {
		return fmt.Errorf("invalid role type: %s", roleStr)
	}
	
	*r = roleType
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
	// Parse the JSON permissions from the database
	var permissions []string
	if r.Permissions != "" {
		if err := json.Unmarshal([]byte(r.Permissions), &permissions); err != nil {
			// If JSON parsing fails, fall back to default permissions for the role type
			return r.hasDefaultPermission(permission)
		}
	} else {
		// If no permissions are stored, use default permissions for the role type
		return r.hasDefaultPermission(permission)
	}

	// Check if the permission is explicitly granted
	for _, perm := range permissions {
		if r.permissionMatches(perm, permission) {
			return true
		}
	}

	return false
}

// hasDefaultPermission checks if the permission is granted by default for the role type
func (r *Role) hasDefaultPermission(permission string) bool {
	defaultPerms := GetDefaultPermissions(r.Name)
	for _, perm := range defaultPerms {
		if r.permissionMatches(perm, permission) {
			return true
		}
	}
	return false
}

// permissionMatches checks if a stored permission matches the requested permission
// Supports wildcard permissions (e.g., "routes.*" matches "routes.read")
func (r *Role) permissionMatches(storedPerm, requestedPerm string) bool {
	// Exact match
	if storedPerm == requestedPerm {
		return true
	}

	// Wildcard match (e.g., "routes.*" matches "routes.read")
	if strings.HasSuffix(storedPerm, ".*") {
		prefix := strings.TrimSuffix(storedPerm, ".*")
		return strings.HasPrefix(requestedPerm, prefix+".")
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