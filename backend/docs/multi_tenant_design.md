# Multi-Tenant Database Design

This document outlines the multi-tenant architecture implemented in the RoutrApp system. The database schema is designed to isolate data between different organizations (tenants) while maintaining efficient data access patterns.

## Core Principles

1. **Tenant Isolation**: Each organization's data is isolated from others
2. **Shared Database, Separate Tables**: All tenants share the same database but with organizational isolation
3. **Organization ID Field**: Every tenant-specific table includes an `organization_id` field
4. **Row-Level Security**: All queries are filtered by organization_id
5. **Middleware-Based Access Control**: Tenant context extracted from JWT or subdomain

## Schema Structure

### Organizations (Tenants)

The organizations table is the core of our multi-tenant architecture:

```sql
CREATE TABLE organizations (
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
```

### Tenant-Scoped Tables

All tenant-scoped tables include:

- `organization_id` (foreign key to organizations table)
- Indexes on `organization_id` for efficient filtering
- Composite indexes when necessary (e.g., unique email per organization)

Example of tenant-scoped users table:

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    -- Other fields
    -- ...
);

-- Unique email per organization
CREATE UNIQUE INDEX idx_users_org_email ON users(organization_id, email);
```

## Entity Relationships

```
organizations (1) --- (*) users
organizations (1) --- (*) technicians
organizations (1) --- (*) routes
organizations (1) --- (*) route_stops
organizations (1) --- (*) route_activities
organizations (1) --- (*) user_sessions
```

## Access Control Implementation

1. **Tenant Middleware**: Extracts tenant context from:

   - JWT claims (`organization_id`)
   - Subdomain parsing
   - Query parameters (for testing/development)

2. **Context Injection**: Sets tenant context in request context

3. **Repository Layer Filtering**: All database queries automatically filter by the tenant context

## Security Considerations

1. **No Cross-Tenant Access**: Queries explicitly filter by organization_id
2. **Soft Deletes**: Most entities use soft deletes (`deleted_at` timestamp)
3. **Cascade Deletion**: When organization is deleted, all related data is removed
4. **JWT Claims**: Organization ID stored in JWT claims for authentication

## Data Migration and Backup

1. **Organization-Based Exports**: Data can be exported per organization
2. **Schema Migrations**: Applied to all tenants simultaneously
3. **Tenant-Specific Backup**: Ability to backup/restore single tenant data

## Performance Optimizations

1. **Indexed Organization ID**: All tenant-scoped tables index the organization_id
2. **Composite Indexes**: Used for common query patterns
3. **Efficient Subdomain Lookup**: Quick tenant resolution from subdomain

## Future Considerations

1. **Horizontal Partitioning**: For large tenants, consider partitioning by organization_id
2. **Redis Caching**: Implement per-tenant caching strategy
3. **Connection Pooling**: Configure connection pools per tenant for heavy usage
