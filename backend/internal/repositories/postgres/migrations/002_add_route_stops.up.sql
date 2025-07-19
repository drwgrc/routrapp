-- Migration: add_route_stops
-- Version: 2
-- Created: 2025-01-17 12:00:00
-- Direction: UP

-- Add route stops table for managing multiple stops in a route
CREATE TABLE IF NOT EXISTS route_stops (
    id SERIAL PRIMARY KEY,
    route_id INTEGER REFERENCES routes(id) ON DELETE CASCADE,
    stop_order INTEGER NOT NULL,
    location_name VARCHAR(255) NOT NULL,
    location_lat DECIMAL(10, 8) NOT NULL,
    location_lng DECIMAL(11, 8) NOT NULL,
    estimated_arrival_time TIMESTAMP,
    actual_arrival_time TIMESTAMP,
    notes TEXT,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_route_stops_route_id ON route_stops(route_id);
CREATE INDEX IF NOT EXISTS idx_route_stops_order ON route_stops(route_id, stop_order);
CREATE INDEX IF NOT EXISTS idx_route_stops_status ON route_stops(status);

-- Add constraint to ensure stop order is unique within a route
ALTER TABLE route_stops ADD CONSTRAINT unique_route_stop_order 
UNIQUE (route_id, stop_order);

-- Insert migration record
INSERT INTO schema_migrations (version, description) 
VALUES (2, 'Add route_stops table for managing multiple stops per route')
ON CONFLICT (version) DO NOTHING; 