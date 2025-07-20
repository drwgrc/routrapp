-- Migration: fix_user_email_index
-- Version: 5
-- Created: 2025-01-17 15:00:00
-- Direction: DOWN

-- Remove migration record
DELETE FROM schema_migrations WHERE version = 5;

-- Drop the composite unique index
DROP INDEX IF EXISTS idx_users_org_email;

-- Note: We don't recreate the incorrect global email index as it would violate multi-tenancy
-- The correct behavior is to have no global email uniqueness constraint 