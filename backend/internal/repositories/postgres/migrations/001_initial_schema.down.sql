-- Migration: initial_schema
-- Version: 1
-- Created: 2025-01-17 12:00:00
-- Direction: DOWN

-- Remove migration record
DELETE FROM schema_migrations WHERE version = 1;

-- Drop tables in reverse order to handle foreign key constraints
DROP TABLE IF EXISTS routes;
DROP TABLE IF EXISTS technicians;
DROP TABLE IF EXISTS users;

-- Drop the migrations tracking table last
DROP TABLE IF EXISTS schema_migrations; 