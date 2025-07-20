-- Migration: add user sessions
-- Version: 3
-- Created: 2025-07-19 00:08:59
-- Direction: DOWN

-- Drop route activities table
DROP INDEX IF EXISTS idx_route_activities_timestamp;
DROP INDEX IF EXISTS idx_route_activities_technician_id;
DROP INDEX IF EXISTS idx_route_activities_route_stop_id;
DROP INDEX IF EXISTS idx_route_activities_route_id;
DROP INDEX IF EXISTS idx_route_activities_organization_id;
DROP TABLE IF EXISTS route_activities CASCADE;

-- Drop user sessions table
DROP INDEX IF EXISTS idx_user_sessions_expires_at;
DROP INDEX IF EXISTS idx_user_sessions_is_valid;
DROP INDEX IF EXISTS idx_user_sessions_refresh_token;
DROP INDEX IF EXISTS idx_user_sessions_user_id;
DROP INDEX IF EXISTS idx_user_sessions_organization_id;
DROP TABLE IF EXISTS user_sessions CASCADE;
