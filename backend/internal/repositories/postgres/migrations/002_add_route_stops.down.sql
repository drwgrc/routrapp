-- Migration: add_route_stops
-- Version: 2
-- Created: 2025-01-17 12:00:00
-- Direction: DOWN

-- Remove migration record
DELETE FROM schema_migrations WHERE version = 2;

-- Drop the route_stops table
DROP TABLE IF EXISTS route_stops; 