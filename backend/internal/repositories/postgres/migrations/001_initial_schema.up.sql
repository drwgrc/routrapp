-- Migration: initial_schema
-- Version: 1
-- Created: 2025-01-17 12:00:00
-- Direction: UP

-- Create migrations tracking table
CREATE TABLE IF NOT EXISTS schema_migrations (
    version INTEGER PRIMARY KEY,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    description TEXT
);

-- Organizations (Tenants) table
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
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_organizations_sub_domain ON organizations(sub_domain);
CREATE INDEX IF NOT EXISTS idx_organizations_active ON organizations(active);

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role VARCHAR(50) DEFAULT 'technician',
    refresh_token VARCHAR(255),
    last_login_at TIMESTAMP,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create index on email for faster lookups (unique per organization)
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_org_email ON users(organization_id, email);
CREATE INDEX IF NOT EXISTS idx_users_organization_id ON users(organization_id);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- Technicians table
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
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_technicians_user_id ON technicians(user_id);
CREATE INDEX IF NOT EXISTS idx_technicians_organization_id ON technicians(organization_id);
CREATE INDEX IF NOT EXISTS idx_technicians_status ON technicians(status);

-- Routes table
CREATE TABLE IF NOT EXISTS routes (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    technician_id INTEGER REFERENCES technicians(id) ON DELETE SET NULL,
    status VARCHAR(20) DEFAULT 'pending',
    scheduled_date TIMESTAMP,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    is_optimized BOOLEAN DEFAULT false,
    total_distance DECIMAL(10, 2) DEFAULT 0,
    total_duration INTEGER DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_routes_organization_id ON routes(organization_id);
CREATE INDEX IF NOT EXISTS idx_routes_technician_id ON routes(technician_id);
CREATE INDEX IF NOT EXISTS idx_routes_status ON routes(status);
CREATE INDEX IF NOT EXISTS idx_routes_scheduled_date ON routes(scheduled_date);

-- Insert initial migration record
INSERT INTO schema_migrations (version, description) 
VALUES (1, 'Initial schema setup with multi-tenant structure')
ON CONFLICT (version) DO NOTHING; 