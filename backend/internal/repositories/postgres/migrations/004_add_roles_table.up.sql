-- Migration: add_roles_table
-- Version: 4
-- Created: 2025-01-17 14:00:00
-- Direction: UP

-- Create roles table
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(20) NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    description TEXT,
    permissions TEXT, -- JSON array of permissions
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for roles
CREATE UNIQUE INDEX IF NOT EXISTS idx_role_org_name ON roles(organization_id, name) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_roles_organization_id ON roles(organization_id);
CREATE INDEX IF NOT EXISTS idx_roles_name ON roles(name);
CREATE INDEX IF NOT EXISTS idx_roles_deleted_at ON roles(deleted_at);

-- Add role_id column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS role_id INTEGER REFERENCES roles(id) ON DELETE RESTRICT;

-- Create index for role_id in users table
CREATE INDEX IF NOT EXISTS idx_users_role_id ON users(role_id);

-- Insert default roles for each existing organization
DO $$
DECLARE
    org_record RECORD;
    owner_role_id INTEGER;
    tech_role_id INTEGER;
BEGIN
    FOR org_record IN SELECT id FROM organizations WHERE deleted_at IS NULL LOOP
        -- Insert Owner role
        INSERT INTO roles (organization_id, name, display_name, description, permissions, active)
        VALUES (
            org_record.id,
            'owner',
            'Organization Owner',
            'Full access to all organization resources and settings',
            '["organizations.*", "users.*", "technicians.*", "routes.*", "roles.*"]',
            true
        )
        RETURNING id INTO owner_role_id;

        -- Insert Technician role
        INSERT INTO roles (organization_id, name, display_name, description, permissions, active)
        VALUES (
            org_record.id,
            'technician',
            'Technician',
            'Access to routes and personal information',
            '["routes.read", "routes.update_status", "technicians.read_own", "technicians.update_own"]',
            true
        )
        RETURNING id INTO tech_role_id;

        -- Update existing users with appropriate roles based on their current role field
        UPDATE users 
        SET role_id = owner_role_id 
        WHERE organization_id = org_record.id 
        AND role = 'owner' 
        AND deleted_at IS NULL;

        UPDATE users 
        SET role_id = tech_role_id 
        WHERE organization_id = org_record.id 
        AND role = 'technician' 
        AND deleted_at IS NULL;

        -- Set default role for any users without a role
        UPDATE users 
        SET role_id = tech_role_id 
        WHERE organization_id = org_record.id 
        AND role_id IS NULL 
        AND deleted_at IS NULL;
    END LOOP;
END
$$;

-- Make role_id NOT NULL after populating it
ALTER TABLE users ALTER COLUMN role_id SET NOT NULL;

-- Remove the old role column (keep it for now to maintain backward compatibility)
-- ALTER TABLE users DROP COLUMN IF EXISTS role;

-- Insert migration record
INSERT INTO schema_migrations (version, description) 
VALUES (4, 'Add roles table and update user role relationships')
ON CONFLICT (version) DO NOTHING; 