package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"routrapp-api/internal/repositories/postgres"
)

func main() {
	var (
		action         = flag.String("action", "", "Action to perform: create, validate, status")
		migrationName  = flag.String("name", "", "Name for new migration (required for create action)")
		migrationsPath = flag.String("path", "internal/repositories/postgres/migrations", "Path to migrations directory")
	)
	flag.Parse()

	if *action == "" {
		printUsage()
		os.Exit(1)
	}

	// Get absolute path for migrations
	absPath, err := filepath.Abs(*migrationsPath)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	manager := postgres.NewMigrationManager(absPath)

	switch *action {
	case "create":
		if *migrationName == "" {
			fmt.Println("Error: migration name is required for create action")
			printUsage()
			os.Exit(1)
		}
		createMigration(manager, *migrationName)
	case "validate":
		validateMigrations(manager)
	case "status":
		showMigrationStatus(manager)
	default:
		fmt.Printf("Error: unknown action '%s'\n", *action)
		printUsage()
		os.Exit(1)
	}
}

func createMigration(manager *postgres.MigrationManager, name string) {
	fmt.Printf("Creating new migration: %s\n", name)
	
	upPath, downPath, err := manager.CreateMigration(name)
	if err != nil {
		log.Fatalf("Failed to create migration: %v", err)
	}
	
	fmt.Printf("✅ Created migration files:\n")
	fmt.Printf("   UP:   %s\n", upPath)
	fmt.Printf("   DOWN: %s\n", downPath)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("1. Edit the SQL files to define your schema changes\n")
	fmt.Printf("2. Validate your migrations: go run scripts/migrate.go -action=validate\n")
	fmt.Printf("3. When database connection is implemented, run: go run scripts/migrate.go -action=up\n")
}

func validateMigrations(manager *postgres.MigrationManager) {
	fmt.Println("Validating migrations...")
	
	err := manager.ValidateMigrations()
	if err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
		os.Exit(1)
	}
	
	migrations, err := manager.LoadMigrations()
	if err != nil {
		log.Fatalf("Failed to load migrations: %v", err)
	}
	
	fmt.Printf("✅ All migrations are valid!\n")
	fmt.Printf("Found %d migration(s):\n", len(migrations))
	
	for _, migration := range migrations {
		fmt.Printf("   %03d: %s\n", migration.Version, migration.Name)
	}
}

func showMigrationStatus(manager *postgres.MigrationManager) {
	fmt.Println("Migration Status")
	fmt.Println("================")
	
	migrations, err := manager.GetMigrationStatus()
	if err != nil {
		log.Fatalf("Failed to get migration status: %v", err)
	}
	
	if len(migrations) == 0 {
		fmt.Println("No migrations found.")
		return
	}
	
	fmt.Printf("Available migrations:\n")
	for _, migration := range migrations {
		status := "📄 Available"
		// TODO: When database connection is implemented, check if migration is applied
		// if migration.Applied {
		//     status = "✅ Applied"
		// }
		
		fmt.Printf("   %03d: %-30s %s\n", migration.Version, migration.Name, status)
	}
	
	fmt.Printf("\nTotal: %d migrations\n", len(migrations))
	fmt.Println("\nNote: Database connection not implemented yet.")
	fmt.Println("Once connected, this will show which migrations have been applied.")
}

func printUsage() {
	fmt.Println("Database Migration Tool")
	fmt.Println("======================")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run scripts/migrate.go -action=<action> [options]")
	fmt.Println()
	fmt.Println("Actions:")
	fmt.Println("  create    Create a new migration file pair (up/down)")
	fmt.Println("  validate  Validate all migration files")
	fmt.Println("  status    Show migration status")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -name     Migration name (required for create)")
	fmt.Println("  -path     Path to migrations directory (default: internal/repositories/postgres/migrations)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run scripts/migrate.go -action=create -name=\"add user sessions\"")
	fmt.Println("  go run scripts/migrate.go -action=validate")
	fmt.Println("  go run scripts/migrate.go -action=status")
	fmt.Println()
	fmt.Println("Future actions (when database connection is implemented):")
	fmt.Println("  up        Apply pending migrations")
	fmt.Println("  down      Rollback last migration")
	fmt.Println("  reset     Rollback all migrations")
} 