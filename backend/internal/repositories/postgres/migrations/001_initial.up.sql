-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

-- Create tenants table
CREATE TABLE IF NOT EXISTS tenants (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

-- Create technicians table
CREATE TABLE IF NOT EXISTS technicians (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    user_id INTEGER NOT NULL REFERENCES users(id),
    tenant_id INTEGER NOT NULL REFERENCES tenants(id),
    phone VARCHAR(20),
    availability_start TIME,
    availability_end TIME,
    max_daily_jobs INTEGER DEFAULT 8,
    skills TEXT[]
);

-- Create routes table
CREATE TABLE IF NOT EXISTS routes (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    technician_id INTEGER REFERENCES technicians(id),
    tenant_id INTEGER NOT NULL REFERENCES tenants(id),
    date DATE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'scheduled',
    start_location JSON,
    end_location JSON,
    total_distance NUMERIC(10, 2),
    total_duration INTEGER
);

-- Create indices for performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_technicians_user_id ON technicians(user_id);
CREATE INDEX IF NOT EXISTS idx_technicians_tenant_id ON technicians(tenant_id);
CREATE INDEX IF NOT EXISTS idx_routes_technician_id ON routes(technician_id);
CREATE INDEX IF NOT EXISTS idx_routes_tenant_id ON routes(tenant_id);
CREATE INDEX IF NOT EXISTS idx_routes_date ON routes(date);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
CREATE INDEX IF NOT EXISTS idx_tenants_deleted_at ON tenants(deleted_at);
CREATE INDEX IF NOT EXISTS idx_technicians_deleted_at ON technicians(deleted_at);
CREATE INDEX IF NOT EXISTS idx_routes_deleted_at ON routes(deleted_at); 

-- Migration: initial
-- Version: 1
-- Created: 2025-01-17 12:00:00
-- Direction: UP

-- Create organizations (tenants) table
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

-- Create first admin organization
INSERT INTO organizations (name, sub_domain, contact_email, active, plan_type, created_at, updated_at)
VALUES ('Admin Organization', 'admin', 'admin@routrapp.com', true, 'enterprise', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Create admin user
INSERT INTO users (organization_id, email, password_hash, first_name, last_name, role, active)
VALUES (
    (SELECT id FROM organizations WHERE sub_domain = 'admin'), 
    'admin@routrapp.com', 
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', -- password: admin123
    'Admin',
    'User',
    'owner',
    true
); 