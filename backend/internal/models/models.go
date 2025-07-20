package models

// This file provides a central import point for all models
// Import this package to get access to all model types

// Re-export all model types for convenience
type (
	// Core entities
	OrganizationModel = Organization
	UserModel         = User
	RoleModel         = Role
	TechnicianModel   = Technician
	RouteModel        = Route
	RouteStopModel    = RouteStop
	
	// Enums and types
	RoleTypeEnum         = RoleType
	TechnicianStatusEnum = TechnicianStatus
	RouteStatusEnum      = RouteStatus
	
	// Embedded types
	TimeWindowType     = TimeWindow
	RouteActivityModel = RouteActivity
)

// Model constants for easy access
const (
	// Role types
	OwnerRole      = RoleTypeOwner
	TechnicianRole = RoleTypeTechnician
	
	// Technician statuses
	TechActiveStatus    = TechnicianStatusActive
	TechInactiveStatus  = TechnicianStatusInactive
	TechOnRouteStatus   = TechnicianStatusOnRoute
	TechOnBreakStatus   = TechnicianStatusOnBreak
	TechOffDutyStatus   = TechnicianStatusOffDuty
	
	// Route statuses
	RoutePendingStatus   = RouteStatusPending
	RouteAssignedStatus  = RouteStatusAssigned
	RouteStartedStatus   = RouteStatusStarted
	RouteCompletedStatus = RouteStatusCompleted
	RouteCancelledStatus = RouteStatusCancelled
	RoutePausedStatus    = RouteStatusPaused
)

// AllModels returns a slice of all model types for auto-migration
func AllModels() []interface{} {
	return []interface{}{
		&Organization{},
		&Role{},
		&User{},
		&Technician{},
		&Route{},
		&RouteStop{},
		&RouteActivity{},
	}
} 