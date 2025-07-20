# Database Migrations

This directory contains database migration files for the RoutRapp application.

## Migration Structure

Migrations are stored as numbered SQL files with both "up" and "down" versions:

- `001_migration_name.up.sql` - Applies the migration
- `001_migration_name.down.sql` - Reverts the migration

## Usage

### Using the Shell Script (Recommended)

```bash
# Navigate to the backend directory
cd backend

# Create a new migration
./scripts/migrate.sh create "add user sessions"

# Validate all migrations
./scripts/migrate.sh validate

# Check migration status
./scripts/migrate.sh status
```

### Using the Go Tool Directly

```bash
# Create a new migration
go run scripts/migrate.go -action=create -name="add user sessions"

# Validate all migrations
go run scripts/migrate.go -action=validate

# Check migration status
go run scripts/migrate.go -action=status
```

## Migration Naming Convention

- Use descriptive names that explain what the migration does
- Use snake_case for migration names
- Examples:
  - `create_users_table`
  - `add_indexes_to_routes`
  - `update_technician_status_enum`

## File Structure

Each migration consists of two files:

### Up Migration (`XXX_name.up.sql`)

Contains SQL statements to apply the changes:

```sql
-- Migration: add_user_sessions
-- Version: 3
-- Created: 2025-01-17 12:00:00
-- Direction: UP

CREATE TABLE user_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_sessions_token ON user_sessions(token);
CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);

-- Insert migration record
INSERT INTO schema_migrations (version, description)
VALUES (3, 'Add user sessions table for authentication')
ON CONFLICT (version) DO NOTHING;
```

### Down Migration (`XXX_name.down.sql`)

Contains SQL statements to revert the changes:

```sql
-- Migration: add_user_sessions
-- Version: 3
-- Created: 2025-01-17 12:00:00
-- Direction: DOWN

-- Remove migration record
DELETE FROM schema_migrations WHERE version = 3;

-- Drop the user_sessions table
DROP TABLE IF EXISTS user_sessions;
```

## Best Practices

1. **Always create both up and down migrations** - This allows for rollbacks
2. **Test your migrations** - Validate before applying to production
3. **Use IF NOT EXISTS / IF EXISTS** - Makes migrations idempotent
4. **Add indexes for performance** - Include relevant database indexes
5. **Document your changes** - Include comments explaining complex migrations
6. **Keep migrations small** - One logical change per migration

## Current Migrations

| Version | Name              | Description                                                                  |
| ------- | ----------------- | ---------------------------------------------------------------------------- |
| 001     | initial           | Creates base multi-tenant schema (organizations, users, technicians, routes) |
| 002     | add_route_stops   | Adds route_stops table for multi-stop routes with constraints                |
| 003     | add_user_sessions | Adds user_sessions and route_activities tables                               |

## Migration Issues Fixed (2025-01-17)

The migration system was cleaned up to resolve several issues:

1. **Duplicate version conflicts** - Removed conflicting `001_initial_schema.up.sql` that had the same version as `001_initial.up.sql`
2. **Schema inconsistencies** - Fixed column name mismatches (`is_active` vs `active`, `password` vs `password_hash`)
3. **Table creation conflicts** - Reorganized migrations to prevent duplicate table creation
4. **Timestamp standardization** - Updated all migrations to use `TIMESTAMP WITH TIME ZONE`
5. **Security improvements** - Removed hardcoded admin password from migration
6. **Index improvements** - Added proper `deleted_at` indexes for soft deletes
7. **Foreign key consistency** - Ensured proper CASCADE behaviors for multi-tenancy

The migration system now follows a clean incremental approach:

- **001_initial**: Core multi-tenant foundation (organizations, users, technicians, routes)
- **002_add_route_stops**: Route stops functionality
- **003_add_user_sessions**: Authentication and activity tracking

## Future Development

When database connection is implemented, additional actions will be available:

- `up` - Apply pending migrations
- `down` - Rollback last migration
- `reset` - Rollback all migrations
- `redo` - Rollback and reapply last migration

## Database Schema Tracking

The system uses a `schema_migrations` table to track which migrations have been applied:

```sql
CREATE TABLE schema_migrations (
    version INTEGER PRIMARY KEY,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    description TEXT
);
```

## Environment Variables

Database configuration can be controlled via environment variables:

- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_NAME` - Database name (default: routrapp)
- `DB_USER` - Database user (default: routrapp)
- `DB_PASSWORD` - Database password (default: routrapp_password)
- `DB_SSL_MODE` - SSL mode (default: disable)
- `DB_MIGRATIONS_PATH` - Migration files path
- `DB_AUTO_MIGRATE` - Auto-run migrations on startup (default: false)

## Troubleshooting

### Common Issues

1. **"Failed to read migrations directory"**

   - Ensure you're running commands from the backend directory
   - Check that the migrations directory exists

2. **"Migration X is missing up/down SQL"**

   - Ensure both .up.sql and .down.sql files exist for each migration
   - Check file naming convention is correct

3. **"Invalid filename format"**
   - Migration files must follow: `001_name.up.sql` / `001_name.down.sql`
   - Version numbers must be sequential integers

### Getting Help

Run the migration tool without arguments to see usage information:

```bash
./scripts/migrate.sh
# or
go run scripts/migrate.go
```
