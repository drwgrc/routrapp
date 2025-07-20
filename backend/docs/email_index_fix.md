# Email Index Fix for Multi-Tenancy

## Problem Description

The User model's email unique index was incorrectly configured, causing a multi-tenancy constraint violation. The issue was:

1. **Incorrect GORM Tag**: The User model had `uniqueIndex:idx_users_org_email` on the email field, but this created a global unique constraint on email instead of a per-organization constraint.

2. **Schema Inconsistency**: The database migration (001_initial.up.sql) correctly created a composite unique index `idx_users_org_email` on `(organization_id, email)`, but the GORM model definition didn't match this schema.

3. **Multi-Tenancy Violation**: This would prevent different organizations from having users with the same email address, which violates the multi-tenancy principle.

## Solution Implemented

### 1. Fixed User Model (`backend/internal/models/user.go`)

**Before:**

```go
Email string `gorm:"type:varchar(100);uniqueIndex:idx_users_org_email" json:"email"`
```

**After:**

```go
Email string `gorm:"type:varchar(100)" json:"email"`
```

And added a proper composite unique index:

```go
// Indexes returns the database indexes for the User model
func (User) Indexes() []string {
    return []string{
        "CREATE UNIQUE INDEX IF NOT EXISTS idx_users_org_email ON users(organization_id, email) WHERE deleted_at IS NULL",
    }
}
```

### 2. Created Migration (`005_fix_user_email_index.up.sql`)

```sql
-- Drop any existing incorrect email unique indexes
DROP INDEX IF EXISTS idx_users_email;

-- Ensure the correct composite unique index exists
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_org_email ON users(organization_id, email) WHERE deleted_at IS NULL;
```

### 3. Fixed User Model Methods

Updated the role checking methods to use role name instead of role ID for better testability:

**Before:**

```go
func (u *User) IsOwner() bool {
    if u.Role.ID == 0 {
        return false // Role not loaded
    }
    return u.Role.Name == RoleTypeOwner
}
```

**After:**

```go
func (u *User) IsOwner() bool {
    if u.Role.Name == "" {
        return false // Role not loaded
    }
    return u.Role.Name == RoleTypeOwner
}
```

### 4. Added Comprehensive Tests (`backend/internal/models/user_test.go`)

Created tests to verify:

- User role checking (IsOwner, IsTechnician)
- Permission checking (HasPermission)
- Multi-tenancy constraints work correctly

## Multi-Tenancy Constraint Behavior

With this fix, the system now correctly enforces:

✅ **Allowed**: Different organizations can have users with the same email address

```
Organization A: user1@example.com
Organization B: user1@example.com  // ✅ Allowed
```

❌ **Not Allowed**: Same organization cannot have duplicate email addresses

```
Organization A: user1@example.com
Organization A: user1@example.com  // ❌ Constraint violation
```

## Database Schema

The correct index structure:

```sql
CREATE UNIQUE INDEX idx_users_org_email ON users(organization_id, email) WHERE deleted_at IS NULL;
```

This index:

- Enforces uniqueness per organization (multi-tenancy)
- Excludes soft-deleted records (`WHERE deleted_at IS NULL`)
- Allows efficient queries by organization and email

## Testing

Run the tests to verify the fix:

```bash
cd backend
go test ./internal/models -v
```

## Migration Instructions

To apply this fix to existing databases:

1. Run the new migration:

```bash
# Using the migration script
go run scripts/migrate.go up
```

2. Verify the index was created correctly:

```sql
-- Check existing indexes
SELECT indexname, indexdef FROM pg_indexes WHERE tablename = 'users';
```

## Impact

- **Positive**: Proper multi-tenancy support for email addresses
- **Positive**: Consistent schema between GORM model and database
- **Positive**: Better testability of user role methods
- **Neutral**: No breaking changes to existing functionality

## Related Files

- `backend/internal/models/user.go` - Fixed User model
- `backend/internal/models/user_test.go` - Added comprehensive tests
- `backend/internal/repositories/postgres/migrations/005_fix_user_email_index.up.sql` - Migration to fix database schema
- `backend/internal/repositories/postgres/migrations/005_fix_user_email_index.down.sql` - Rollback migration
