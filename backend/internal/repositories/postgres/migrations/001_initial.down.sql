-- Drop indices
DROP INDEX IF EXISTS idx_routes_deleted_at;
DROP INDEX IF EXISTS idx_technicians_deleted_at;
DROP INDEX IF EXISTS idx_tenants_deleted_at;
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_routes_date;
DROP INDEX IF EXISTS idx_routes_tenant_id;
DROP INDEX IF EXISTS idx_routes_technician_id;
DROP INDEX IF EXISTS idx_technicians_tenant_id;
DROP INDEX IF EXISTS idx_technicians_user_id;
DROP INDEX IF EXISTS idx_users_email;

-- Drop tables (in reverse order of creation to handle references)
DROP TABLE IF EXISTS routes;
DROP TABLE IF EXISTS technicians;
DROP TABLE IF EXISTS tenants;
DROP TABLE IF EXISTS users; 