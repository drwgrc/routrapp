-- Migration: add_roles_table
-- Version: 4
-- Created: 2025-01-17 14:00:00
-- Direction: DOWN

-- Remove migration record
DELETE FROM schema_migrations WHERE version = 4;

-- Add back the role column if it was dropped
ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(20) DEFAULT 'technician';

-- Populate role column from role_id relationship
UPDATE users 
SET role = roles.name 
FROM roles 
WHERE users.role_id = roles.id 
AND users.deleted_at IS NULL;

-- Remove role_id column from users
ALTER TABLE users DROP COLUMN IF EXISTS role_id;

-- Drop indexes
DROP INDEX IF EXISTS idx_users_role_id;
DROP INDEX IF EXISTS idx_role_org_name;
DROP INDEX IF EXISTS idx_roles_organization_id;
DROP INDEX IF EXISTS idx_roles_name;
DROP INDEX IF EXISTS idx_roles_deleted_at;

-- Drop roles table
DROP TABLE IF EXISTS roles; 