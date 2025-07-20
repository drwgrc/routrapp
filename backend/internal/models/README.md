# Models Package

This package contains all database models for the RouteApp system. The models implement a multi-tenant architecture where all tenant-scoped entities embed the `Base` struct which includes `OrganizationID` for tenant isolation.

## Model Structure

### Core Models

#### Organization

- **Purpose**: Root entity for multi-tenancy
- **Note**: Does not embed `Base` since it's the tenant root
- **Relationships**: Has many Users, Technicians, Routes

#### User

- **Purpose**: System users (owners and technicians)
- **Tenant-scoped**: Yes (embeds `Base`)
- **Relationships**: Belongs to Organization and Role, has one Technician (optional)

#### Role

- **Purpose**: Role-Based Access Control (RBAC)
- **Tenant-scoped**: Yes (embeds `Base`)
- **Types**: Owner, Technician
- **Features**: JSON permissions field for granular access control

#### Technician

- **Purpose**: Field technicians with location tracking
- **Tenant-scoped**: Yes (embeds `Base`)
- **Relationships**: Belongs to User and Organization, has many Routes

#### Route

- **Purpose**: Route optimization and management
- **Tenant-scoped**: Yes (embeds `Base`)
- **Relationships**: Belongs to Organization and Technician, has many RouteStops

#### RouteStop

- **Purpose**: Individual stops/waypoints in a route
- **Tenant-scoped**: Yes (embeds `Base`)
- **Features**: GPS coordinates, time windows, completion tracking

#### RouteActivity

- **Purpose**: Activity tracking during route execution
- **Tenant-scoped**: Yes (embeds `Base`)
- **Features**: GPS tracking, photos, notes

### Base Struct

All tenant-scoped models embed the `Base` struct:

```go
type Base struct {
    ID             uint           `gorm:"primaryKey" json:"id"`
    OrganizationID uint           `gorm:"not null;index" json:"organization_id"`
    CreatedAt      time.Time      `json:"created_at"`
    UpdatedAt      time.Time      `json:"updated_at"`
    DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
```

## Multi-Tenancy

- All models (except Organization) include `organization_id` for tenant isolation
- Database indexes ensure efficient tenant-scoped queries
- Row-level security can be implemented at the database level
- Foreign key constraints maintain referential integrity within tenants

## Role-Based Access Control (RBAC)

### Role Types

- **Owner**: Full access to organization resources
- **Technician**: Limited access to own routes and profile

### Permissions

Permissions are stored as JSON arrays in the `permissions` field:

#### Owner Permissions

```json
["organizations.*", "users.*", "technicians.*", "routes.*", "roles.*"]
```

#### Technician Permissions

```json
[
  "routes.read",
  "routes.update_status",
  "technicians.read_own",
  "technicians.update_own"
]
```

## Model Relationships

```
Organization (1) --> (N) Users
Organization (1) --> (N) Roles
Organization (1) --> (N) Technicians
Organization (1) --> (N) Routes

User (N) --> (1) Role
User (1) --> (1) Technician (optional)

Route (N) --> (1) Technician
Route (1) --> (N) RouteStops
Route (1) --> (N) RouteActivities

Technician (N) --> (1) User
```

## Usage

### Importing Models

```go
import "your-project/internal/models"

// Use individual models
user := &models.User{}
org := &models.Organization{}

// Or use the central export
type User = models.UserModel
```

### Auto-Migration

```go
import "your-project/internal/models"

// Migrate all models
db.AutoMigrate(models.AllModels()...)
```

### Querying with Tenant Scope

```go
// Always include organization_id in queries for tenant-scoped models
var users []models.User
db.Where("organization_id = ?", orgID).Find(&users)

// Use preloading for relationships
var user models.User
db.Preload("Role").Preload("Technician").First(&user, userID)
```

## Status Enums

### TechnicianStatus

- `active`: Available for route assignment
- `inactive`: Not available
- `on_route`: Currently executing a route
- `on_break`: On scheduled break
- `off_duty`: End of work day

### RouteStatus

- `pending`: Created but not assigned
- `assigned`: Assigned to technician
- `started`: Route execution started
- `completed`: All stops completed
- `cancelled`: Route cancelled
- `paused`: Temporarily paused

## Database Migrations

Models are created through database migrations in `repositories/postgres/migrations/`:

1. **001_initial.up.sql**: Organizations, Users, Technicians, Routes
2. **002_add_route_stops.up.sql**: RouteStops and RouteActivities
3. **003_add_user_sessions.up.sql**: User session management
4. **004_add_roles_table.up.sql**: Roles and RBAC

## Validation

- Email uniqueness per organization
- Role name uniqueness per organization
- Proper foreign key constraints
- Not null constraints on required fields
- Default values for status fields

## Indexes

- Primary keys on all tables
- Foreign key indexes for relationships
- Composite indexes for tenant isolation
- Soft delete indexes (deleted_at)
- Query optimization indexes (status, dates)

## Security Considerations

- Passwords are excluded from JSON serialization
- Refresh tokens are excluded from JSON
- Soft deletes prevent accidental data loss
- Foreign key constraints prevent orphaned records
- Tenant isolation through organization_id
