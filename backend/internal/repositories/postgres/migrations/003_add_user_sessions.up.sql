-- Migration: add user sessions
-- Version: 3
-- Created: 2025-07-19 00:08:59
-- Direction: UP

-- Create user sessions table for authentication
CREATE TABLE IF NOT EXISTS user_sessions (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token VARCHAR(255) NOT NULL,
    access_token_hash VARCHAR(255) NOT NULL,
    ip_address VARCHAR(50),
    user_agent TEXT,
    is_valid BOOLEAN DEFAULT true,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_user_sessions_organization_id ON user_sessions(organization_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions(user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_sessions_refresh_token ON user_sessions(refresh_token);
CREATE INDEX IF NOT EXISTS idx_user_sessions_is_valid ON user_sessions(is_valid);
CREATE INDEX IF NOT EXISTS idx_user_sessions_expires_at ON user_sessions(expires_at);

-- Create route activities table for tracking activities during routes
CREATE TABLE IF NOT EXISTS route_activities (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    route_id INTEGER NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
    route_stop_id INTEGER REFERENCES route_stops(id) ON DELETE SET NULL,
    technician_id INTEGER NOT NULL REFERENCES technicians(id) ON DELETE CASCADE,
    activity_type VARCHAR(50) NOT NULL,
    notes TEXT,
    lat DECIMAL(10, 8),
    lng DECIMAL(11, 8),
    photo_url VARCHAR(255),
    timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for route activities
CREATE INDEX IF NOT EXISTS idx_route_activities_organization_id ON route_activities(organization_id);
CREATE INDEX IF NOT EXISTS idx_route_activities_route_id ON route_activities(route_id);
CREATE INDEX IF NOT EXISTS idx_route_activities_route_stop_id ON route_activities(route_stop_id);
CREATE INDEX IF NOT EXISTS idx_route_activities_technician_id ON route_activities(technician_id);
CREATE INDEX IF NOT EXISTS idx_route_activities_timestamp ON route_activities(timestamp);

-- Insert migration record
INSERT INTO schema_migrations (version, description) 
VALUES (3, 'Add user sessions and route activities tables with multi-tenant support')
ON CONFLICT (version) DO NOTHING;
