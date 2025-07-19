package postgres

import (
	"time"
)

// MigrationRecord represents a migration record in the database
type MigrationRecord struct {
	Version     int       `json:"version"`
	Description string    `json:"description"`
	AppliedAt   time.Time `json:"applied_at"`
}

// MigrationTracker interface defines methods for tracking migration state
// This will be implemented when database connection is added
type MigrationTracker interface {
	// GetAppliedMigrations returns all migrations that have been applied
	GetAppliedMigrations() ([]MigrationRecord, error)
	
	// IsMigrationApplied checks if a specific migration version has been applied
	IsMigrationApplied(version int) (bool, error)
	
	// MarkMigrationApplied marks a migration as applied
	MarkMigrationApplied(version int, description string) error
	
	// MarkMigrationRolledBack removes a migration record (for rollbacks)
	MarkMigrationRolledBack(version int) error
	
	// GetCurrentSchemaVersion returns the highest applied migration version
	GetCurrentSchemaVersion() (int, error)
	
	// InitializeTrackingTable creates the schema_migrations table if it doesn't exist
	InitializeTrackingTable() error
}

// MockMigrationTracker is a mock implementation for testing and development
// This can be used until database connection is implemented
type MockMigrationTracker struct {
	appliedMigrations map[int]MigrationRecord
}

// NewMockMigrationTracker creates a new mock migration tracker
func NewMockMigrationTracker() *MockMigrationTracker {
	return &MockMigrationTracker{
		appliedMigrations: make(map[int]MigrationRecord),
	}
}

// GetAppliedMigrations returns all applied migrations from the mock tracker
func (mt *MockMigrationTracker) GetAppliedMigrations() ([]MigrationRecord, error) {
	var records []MigrationRecord
	for _, record := range mt.appliedMigrations {
		records = append(records, record)
	}
	return records, nil
}

// IsMigrationApplied checks if a migration is applied in the mock tracker
func (mt *MockMigrationTracker) IsMigrationApplied(version int) (bool, error) {
	_, exists := mt.appliedMigrations[version]
	return exists, nil
}

// MarkMigrationApplied marks a migration as applied in the mock tracker
func (mt *MockMigrationTracker) MarkMigrationApplied(version int, description string) error {
	mt.appliedMigrations[version] = MigrationRecord{
		Version:     version,
		Description: description,
		AppliedAt:   time.Now(),
	}
	return nil
}

// MarkMigrationRolledBack removes a migration from the mock tracker
func (mt *MockMigrationTracker) MarkMigrationRolledBack(version int) error {
	delete(mt.appliedMigrations, version)
	return nil
}

// GetCurrentSchemaVersion returns the highest applied migration version
func (mt *MockMigrationTracker) GetCurrentSchemaVersion() (int, error) {
	maxVersion := 0
	for version := range mt.appliedMigrations {
		if version > maxVersion {
			maxVersion = version
		}
	}
	return maxVersion, nil
}

// InitializeTrackingTable is a no-op for the mock tracker
func (mt *MockMigrationTracker) InitializeTrackingTable() error {
	return nil
}

// Future implementation note:
// When database connection is added, create a PostgresMigrationTracker
// that implements the MigrationTracker interface with actual database operations:
//
// type PostgresMigrationTracker struct {
//     db *sql.DB
// }
//
// func NewPostgresMigrationTracker(db *sql.DB) *PostgresMigrationTracker {
//     return &PostgresMigrationTracker{db: db}
// }
//
// Then implement all the interface methods with actual SQL queries 