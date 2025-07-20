-- Migration: initial
-- Version: 1
-- Created: 2025-01-17 12:00:00
-- Direction: UP

-- Create migrations tracking table
CREATE TABLE IF NOT EXISTS schema_migrations (
    version INTEGER PRIMARY KEY,
    applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    description TEXT
);

-- Organizations (Tenants) table - Base entity for multi-tenancy
-- Note: Organizations don't have organization_id since they are the root entity
CREATE TABLE IF NOT EXISTS organizations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    sub_domain VARCHAR(100) NOT NULL,
    contact_email VARCHAR(100) NOT NULL,
    contact_phone VARCHAR(20),
    logo_url VARCHAR(255),
    primary_color VARCHAR(20),
    secondary_color VARCHAR(20),
    active BOOLEAN DEFAULT true,
    plan_type VARCHAR(20) DEFAULT 'basic',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for organizations
CREATE UNIQUE INDEX IF NOT EXISTS idx_organizations_sub_domain ON organizations(sub_domain);
CREATE INDEX IF NOT EXISTS idx_organizations_active ON organizations(active);
CREATE INDEX IF NOT EXISTS idx_organizations_deleted_at ON organizations(deleted_at);

-- Users table - All users belong to an organization
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role VARCHAR(20) DEFAULT 'technician',
    last_login_at TIMESTAMP WITH TIME ZONE,
    refresh_token VARCHAR(255),
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for users
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_org_email ON users(organization_id, email) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_organization_id ON users(organization_id);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- Technicians table - Technical staff with location tracking
CREATE TABLE IF NOT EXISTS technicians (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) DEFAULT 'inactive',
    phone_number VARCHAR(20),
    current_lat DECIMAL(10, 8),
    current_lng DECIMAL(11, 8),
    last_location_at BIGINT,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for technicians
CREATE UNIQUE INDEX IF NOT EXISTS idx_technicians_user_id ON technicians(user_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_technicians_organization_id ON technicians(organization_id);
CREATE INDEX IF NOT EXISTS idx_technicians_status ON technicians(status);
CREATE INDEX IF NOT EXISTS idx_technicians_deleted_at ON technicians(deleted_at);

-- Routes table - Main routing entity
CREATE TABLE IF NOT EXISTS routes (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    technician_id INTEGER REFERENCES technicians(id) ON DELETE SET NULL,
    status VARCHAR(20) DEFAULT 'pending',
    scheduled_date TIMESTAMP WITH TIME ZONE,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    cancelled_at TIMESTAMP WITH TIME ZONE,
    is_optimized BOOLEAN DEFAULT false,
    total_distance DECIMAL(10, 2) DEFAULT 0,
    total_duration INTEGER DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for routes
CREATE INDEX IF NOT EXISTS idx_routes_organization_id ON routes(organization_id);
CREATE INDEX IF NOT EXISTS idx_routes_technician_id ON routes(technician_id);
CREATE INDEX IF NOT EXISTS idx_routes_status ON routes(status);
CREATE INDEX IF NOT EXISTS idx_routes_scheduled_date ON routes(scheduled_date);
CREATE INDEX IF NOT EXISTS idx_routes_deleted_at ON routes(deleted_at);

-- Note: route_stops and route_activities tables will be created in subsequent migrations
-- to maintain proper migration history and incremental schema changes

-- Insert initial migration record
INSERT INTO schema_migrations (version, description) 
VALUES (1, 'Initial schema setup with multi-tenant structure')
ON CONFLICT (version) DO NOTHING; 