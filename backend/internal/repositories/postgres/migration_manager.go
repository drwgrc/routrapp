package postgres

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Migration represents a single database migration
type Migration struct {
	Version   int
	Name      string
	UpSQL     string
	DownSQL   string
	Timestamp time.Time
}

// MigrationManager handles database migrations
type MigrationManager struct {
	migrationsPath string
	// db will be added when database connection is implemented
	// db *sql.DB
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(migrationsPath string) *MigrationManager {
	return &MigrationManager{
		migrationsPath: migrationsPath,
	}
}

// LoadMigrations loads all migration files from the migrations directory
func (mm *MigrationManager) LoadMigrations() ([]Migration, error) {
	var migrations []Migration
	
	entries, err := os.ReadDir(mm.migrationsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}
	
	migrationMap := make(map[int]*Migration)
	
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		filename := entry.Name()
		if !strings.HasSuffix(filename, ".sql") {
			continue
		}
		
		version, name, direction, err := parseMigrationFilename(filename)
		if err != nil {
			continue // Skip invalid filenames
		}
		
		content, err := os.ReadFile(filepath.Join(mm.migrationsPath, filename))
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}
		
		if migrationMap[version] == nil {
			migrationMap[version] = &Migration{
				Version: version,
				Name:    name,
			}
		}
		
		if direction == "up" {
			migrationMap[version].UpSQL = string(content)
		} else if direction == "down" {
			migrationMap[version].DownSQL = string(content)
		}
	}
	
	// Convert map to slice and sort by version
	for _, migration := range migrationMap {
		migrations = append(migrations, *migration)
	}
	
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})
	
	return migrations, nil
}

// CreateMigration creates a new migration file pair (up and down)
func (mm *MigrationManager) CreateMigration(name string) (string, string, error) {
	// Get next version number
	migrations, err := mm.LoadMigrations()
	if err != nil {
		return "", "", fmt.Errorf("failed to load existing migrations: %w", err)
	}
	
	nextVersion := 1
	if len(migrations) > 0 {
		nextVersion = migrations[len(migrations)-1].Version + 1
	}
	
	// Generate filenames
	versionStr := fmt.Sprintf("%03d", nextVersion)
	upFilename := fmt.Sprintf("%s_%s.up.sql", versionStr, strings.ReplaceAll(name, " ", "_"))
	downFilename := fmt.Sprintf("%s_%s.down.sql", versionStr, strings.ReplaceAll(name, " ", "_"))
	
	upPath := filepath.Join(mm.migrationsPath, upFilename)
	downPath := filepath.Join(mm.migrationsPath, downFilename)
	
	// Create up migration file
	upContent := fmt.Sprintf(`-- Migration: %s
-- Version: %d
-- Created: %s
-- Direction: UP

-- Add your SQL statements here to apply this migration
-- Example:
-- CREATE TABLE example_table (
--     id SERIAL PRIMARY KEY,
--     name VARCHAR(255) NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );
`, name, nextVersion, time.Now().Format("2006-01-02 15:04:05"))
	
	if err := os.WriteFile(upPath, []byte(upContent), 0644); err != nil {
		return "", "", fmt.Errorf("failed to create up migration file: %w", err)
	}
	
	// Create down migration file
	downContent := fmt.Sprintf(`-- Migration: %s
-- Version: %d
-- Created: %s
-- Direction: DOWN

-- Add your SQL statements here to rollback this migration
-- Example:
-- DROP TABLE IF EXISTS example_table;
`, name, nextVersion, time.Now().Format("2006-01-02 15:04:05"))
	
	if err := os.WriteFile(downPath, []byte(downContent), 0644); err != nil {
		return "", "", fmt.Errorf("failed to create down migration file: %w", err)
	}
	
	return upPath, downPath, nil
}

// ValidateMigrations checks if all migrations have both up and down files
func (mm *MigrationManager) ValidateMigrations() error {
	migrations, err := mm.LoadMigrations()
	if err != nil {
		return err
	}
	
	for _, migration := range migrations {
		if migration.UpSQL == "" {
			return fmt.Errorf("migration %d (%s) is missing up SQL", migration.Version, migration.Name)
		}
		if migration.DownSQL == "" {
			return fmt.Errorf("migration %d (%s) is missing down SQL", migration.Version, migration.Name)
		}
	}
	
	return nil
}

// GetMigrationStatus returns the current state of migrations
// This will be enhanced when database connection is added
func (mm *MigrationManager) GetMigrationStatus() ([]Migration, error) {
	migrations, err := mm.LoadMigrations()
	if err != nil {
		return nil, err
	}
	
	// TODO: When database connection is implemented, check which migrations 
	// have been applied by querying the migrations table
	// For now, just return all available migrations
	
	return migrations, nil
}

// parseMigrationFilename parses migration filename and extracts version, name, and direction
// Expected format: 001_migration_name.up.sql or 001_migration_name.down.sql
func parseMigrationFilename(filename string) (version int, name string, direction string, err error) {
	// Remove .sql extension
	basename := strings.TrimSuffix(filename, ".sql")
	
	// Split by last dot to get direction
	parts := strings.Split(basename, ".")
	if len(parts) != 2 {
		return 0, "", "", fmt.Errorf("invalid filename format: %s", filename)
	}
	
	direction = parts[1]
	if direction != "up" && direction != "down" {
		return 0, "", "", fmt.Errorf("invalid direction in filename: %s", filename)
	}
	
	// Split the first part to get version and name
	namePart := parts[0]
	underscoreIndex := strings.Index(namePart, "_")
	if underscoreIndex == -1 {
		return 0, "", "", fmt.Errorf("invalid filename format: %s", filename)
	}
	
	versionStr := namePart[:underscoreIndex]
	name = namePart[underscoreIndex+1:]
	
	version, err = strconv.Atoi(versionStr)
	if err != nil {
		return 0, "", "", fmt.Errorf("invalid version number in filename: %s", filename)
	}
	
	return version, name, direction, nil
} 