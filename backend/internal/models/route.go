package models

import "time"

// RouteStatus represents the current status of a route
type RouteStatus string

// Route status constants
const (
	RouteStatusPending   RouteStatus = "pending"
	RouteStatusAssigned  RouteStatus = "assigned"
	RouteStatusStarted   RouteStatus = "started"
	RouteStatusCompleted RouteStatus = "completed"
	RouteStatusCancelled RouteStatus = "cancelled"
	RouteStatusPaused    RouteStatus = "paused"
)

// Route represents a route in the system
type Route struct {
	Base
	Name          string      `gorm:"type:varchar(100)" json:"name"`
	Description   string      `gorm:"type:text" json:"description,omitempty"`
	TechnicianID  *uint       `gorm:"index" json:"technician_id,omitempty"`
	Technician    *Technician `gorm:"foreignKey:TechnicianID" json:"technician,omitempty"`
	Status        RouteStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	ScheduledDate *time.Time  `json:"scheduled_date,omitempty"`
	StartedAt     *time.Time  `json:"started_at,omitempty"`
	CompletedAt   *time.Time  `json:"completed_at,omitempty"`
	CancelledAt   *time.Time  `json:"cancelled_at,omitempty"`
	Stops         []RouteStop `gorm:"foreignKey:RouteID" json:"stops,omitempty"`
	IsOptimized   bool        `gorm:"default:false" json:"is_optimized"`
	TotalDistance float64     `gorm:"type:decimal(10,2)" json:"total_distance"`
	TotalDuration int         `json:"total_duration"` // in seconds
	Notes         string      `gorm:"type:text" json:"notes,omitempty"`
}

// RouteStop represents a stop/waypoint in a route
type RouteStop struct {
	Base
	RouteID      uint       `gorm:"index" json:"route_id"`
	Name         string     `gorm:"type:varchar(100)" json:"name"`
	Address      string     `gorm:"type:varchar(255)" json:"address"`
	Lat          float64    `json:"lat"`
	Lng          float64    `json:"lng"`
	SequenceNum  int        `json:"sequence_num"`
	StopType     string     `gorm:"type:varchar(20)" json:"stop_type"`
	Duration     int        `json:"duration"` // estimated time at stop in minutes
	Notes        string     `gorm:"type:text" json:"notes,omitempty"`
	TimeWindow   *TimeWindow `gorm:"embedded" json:"time_window,omitempty"`
	IsCompleted  bool       `gorm:"default:false" json:"is_completed"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	PhotosCount  int        `gorm:"default:0" json:"photos_count"`
	NotesCount   int        `gorm:"default:0" json:"notes_count"`
}

// TimeWindow represents a time window for a route stop
type TimeWindow struct {
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
}

// RouteActivity represents activities performed during a route
type RouteActivity struct {
	Base
	RouteID      uint      `gorm:"index" json:"route_id"`
	RouteStopID  *uint     `gorm:"index" json:"route_stop_id,omitempty"`
	TechnicianID uint      `gorm:"index" json:"technician_id"`
	ActivityType string    `gorm:"type:varchar(50)" json:"activity_type"`
	Notes        string    `gorm:"type:text" json:"notes,omitempty"`
	Lat          *float64  `json:"lat,omitempty"`
	Lng          *float64  `json:"lng,omitempty"`
	PhotoURL     string    `json:"photo_url,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
} 