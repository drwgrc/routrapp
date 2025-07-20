-- Migration: add_route_stops
-- Version: 2
-- Created: 2025-01-17 12:00:00
-- Direction: UP

-- Add route stops table for managing multiple stops in a route
CREATE TABLE IF NOT EXISTS route_stops (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    route_id INTEGER REFERENCES routes(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    address VARCHAR(255) NOT NULL,
    lat DECIMAL(10, 8) NOT NULL,
    lng DECIMAL(11, 8) NOT NULL,
    sequence_num INTEGER NOT NULL,
    stop_type VARCHAR(20),
    duration INTEGER, -- estimated time at stop in minutes
    notes TEXT,
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    is_completed BOOLEAN DEFAULT false,
    completed_at TIMESTAMP,
    photos_count INTEGER DEFAULT 0,
    notes_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_route_stops_organization_id ON route_stops(organization_id);
CREATE INDEX IF NOT EXISTS idx_route_stops_route_id ON route_stops(route_id);
CREATE INDEX IF NOT EXISTS idx_route_stops_sequence ON route_stops(route_id, sequence_num);
CREATE INDEX IF NOT EXISTS idx_route_stops_completed ON route_stops(is_completed);

-- Add constraint to ensure stop sequence is unique within a route
ALTER TABLE route_stops ADD CONSTRAINT unique_route_stop_sequence 
UNIQUE (route_id, sequence_num);

-- Insert migration record
INSERT INTO schema_migrations (version, description) 
VALUES (2, 'Add route_stops table with multi-tenant support')
ON CONFLICT (version) DO NOTHING; 