-- Migration: initial
-- Version: 1
-- Created: 2025-01-17 12:00:00
-- Direction: DOWN

-- Remove initial migration record
DELETE FROM schema_migrations WHERE version = 1;

-- Drop tables in reverse order to handle foreign key constraints
-- Note: route_activities and route_stops are handled by subsequent migrations
-- Drop routes (depend on technicians)
DROP TABLE IF EXISTS routes CASCADE;

-- Drop technicians (depend on users)
DROP TABLE IF EXISTS technicians CASCADE;

-- Drop users (depend on organizations)
DROP TABLE IF EXISTS users CASCADE;

-- Drop organizations (root entity)
DROP TABLE IF EXISTS organizations CASCADE;

-- Drop the migrations tracking table last
DROP TABLE IF EXISTS schema_migrations CASCADE; 