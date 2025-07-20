-- Migration: add_route_stops
-- Version: 2
-- Created: 2025-01-17 12:00:00
-- Direction: DOWN

-- Remove constraint
ALTER TABLE IF EXISTS route_stops DROP CONSTRAINT IF EXISTS unique_route_stop_sequence;

-- Drop indices
DROP INDEX IF EXISTS idx_route_stops_completed;
DROP INDEX IF EXISTS idx_route_stops_sequence;
DROP INDEX IF EXISTS idx_route_stops_route_id;
DROP INDEX IF EXISTS idx_route_stops_organization_id;

-- Drop table
DROP TABLE IF EXISTS route_stops CASCADE; 