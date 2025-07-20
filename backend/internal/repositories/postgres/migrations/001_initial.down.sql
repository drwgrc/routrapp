-- Migration: initial
-- Version: 1
-- Created: 2025-01-17 12:00:00
-- Direction: DOWN

-- Drop tables in reverse order of creation to handle dependencies
DROP TABLE IF EXISTS routes CASCADE;
DROP TABLE IF EXISTS technicians CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS organizations CASCADE; 