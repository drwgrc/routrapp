-- Migration: fix_user_email_index
-- Version: 5
-- Created: 2025-01-17 15:00:00
-- Direction: UP

-- Drop any existing incorrect email unique indexes
-- This handles the case where GORM might have created a global unique index on email
DROP INDEX IF EXISTS idx_users_email;

-- Ensure the correct composite unique index exists
-- This index enforces email uniqueness per organization (multi-tenancy constraint)
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_org_email ON users(organization_id, email) WHERE deleted_at IS NULL;

-- Insert migration record
INSERT INTO schema_migrations (version, description) 
VALUES (5, 'Fix user email index to enforce multi-tenancy constraints')
ON CONFLICT (version) DO NOTHING; 