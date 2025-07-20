# Role-Based Component Visibility System

This document describes the comprehensive role-based component visibility system implemented in the RoutRapp frontend. The system provides robust, flexible, and performant ways to control component visibility based on user roles and permissions.

## Overview

The visibility system consists of several layers:

1. **Basic Visibility Guards** - Simple role and permission-based components
2. **Advanced Visibility Guards** - More complex components with additional features
3. **Conditional Rendering** - Components with loading states and error handling
4. **Utility Hooks** - Helper functions for complex visibility logic
5. **Convenience Components** - Pre-built components for common use cases

## Core Components

### Basic Visibility Guards

#### RoleGuard

Simple role-based visibility component.

```tsx
import { RoleGuard } from "@/components/auth";

<RoleGuard allowedRoles={["owner"]}>
  <div>Owner only content</div>
</RoleGuard>

<RoleGuard
  allowedRoles={["owner", "technician"]}
  fallback={<div>Access denied</div>}
>
  <div>Content for owners and technicians</div>
</RoleGuard>
```

#### PermissionGuard

Simple permission-based visibility component.

```tsx
import { PermissionGuard } from "@/components/auth";

<PermissionGuard requiredPermissions={["users.manage"]}>
  <div>User management content</div>
</PermissionGuard>

<PermissionGuard
  requiredPermissions={["routes.read", "routes.create"]}
  requireAllPermissions={true}
>
  <div>Content requiring both read and create permissions</div>
</PermissionGuard>
```

### Advanced Visibility Guards

#### RoleVisibility

Advanced role-based visibility with additional features.

```tsx
import { RoleVisibility } from "@/components/auth";

<RoleVisibility
  allowedRoles={["owner"]}
  fallback={<div>Access denied</div>}
  loadingFallback={<div>Loading...</div>}
  inverse={false}
  showOnError={false}
  requireAllRoles={false}
>
  <div>Owner content</div>
</RoleVisibility>;
```

**Props:**

- `allowedRoles`: Array of allowed user roles
- `fallback`: Content to show when access is denied
- `loadingFallback`: Content to show while loading
- `inverse`: If true, shows content when user DOESN'T have the roles
- `showOnError`: Whether to show content on authentication errors
- `requireAllRoles`: If true, user must have ALL roles (AND logic)

#### PermissionVisibility

Advanced permission-based visibility with additional features.

```tsx
import { PermissionVisibility } from "@/components/auth";

<PermissionVisibility
  requiredPermissions={["users.manage", "routes.manage"]}
  requireAllPermissions={false}
  fallback={<div>Insufficient permissions</div>}
  loadingFallback={<div>Checking permissions...</div>}
>
  <div>Management content</div>
</PermissionVisibility>;
```

**Props:**

- `requiredPermissions`: Array of required permissions
- `requireAllPermissions`: If true, user must have ALL permissions (AND logic)
- `fallback`: Content to show when access is denied
- `loadingFallback`: Content to show while loading
- `inverse`: If true, shows content when user DOESN'T have the permissions
- `showOnError`: Whether to show content on authentication errors

#### CombinedVisibility

Combines role and permission checks with flexible logic.

```tsx
import { CombinedVisibility } from "@/components/auth";

// AND logic: User must be owner AND have organization management
<CombinedVisibility
  allowedRoles={["owner"]}
  requiredPermissions={["organizations.manage"]}
  logic="AND"
>
  <div>Owner with org management</div>
</CombinedVisibility>

// OR logic: User must be owner OR have user management
<CombinedVisibility
  allowedRoles={["owner"]}
  requiredPermissions={["users.manage"]}
  logic="OR"
>
  <div>Owner or user manager</div>
</CombinedVisibility>
```

**Props:**

- `allowedRoles`: Array of allowed user roles
- `requiredPermissions`: Array of required permissions
- `logic`: "AND" or "OR" to combine role and permission checks
- `requireAllRoles`: If true, user must have ALL roles
- `requireAllPermissions`: If true, user must have ALL permissions
- `fallback`: Content to show when access is denied
- `loadingFallback`: Content to show while loading
- `inverse`: If true, shows content when user DOESN'T have access
- `showOnError`: Whether to show content on authentication errors

#### ConditionalVisibility

Custom condition-based visibility.

```tsx
import { ConditionalVisibility } from "@/components/auth";

<ConditionalVisibility
  condition={() => user.role === "owner" && user.active}
  fallback={<div>User must be active owner</div>}
>
  <div>Active owner content</div>
</ConditionalVisibility>;
```

#### FeatureFlag

Feature flag component for A/B testing and gradual rollouts.

```tsx
import { FeatureFlag } from "@/components/auth";

<FeatureFlag
  feature="advanced-analytics"
  enabledFor={["owner"]}
  enabledWithPermissions={["analytics.*"]}
  fallback={<div>Feature not available</div>}
>
  <div>Advanced analytics feature</div>
</FeatureFlag>;
```

### Conditional Rendering Components

These components provide loading states and error handling for better user experience.

#### RoleConditionalRender

Role-based rendering with loading and error states.

```tsx
import { RoleConditionalRender } from "@/components/auth";

<RoleConditionalRender
  allowedRoles={["owner"]}
  loadingFallback={<div>Loading...</div>}
  errorFallback="Authentication error occurred"
  showOnError={false}
  className="my-custom-class"
>
  <div>Owner content with proper loading states</div>
</RoleConditionalRender>;
```

#### PermissionConditionalRender

Permission-based rendering with loading and error states.

```tsx
import { PermissionConditionalRender } from "@/components/auth";

<PermissionConditionalRender
  requiredPermissions={["users.manage"]}
  loadingFallback={<div>Checking permissions...</div>}
  errorFallback="Permission check failed"
  showOnError={false}
>
  <div>User management with proper loading states</div>
</PermissionConditionalRender>;
```

#### CombinedConditionalRender

Combined role and permission rendering with loading and error states.

```tsx
import { CombinedConditionalRender } from "@/components/auth";

<CombinedConditionalRender
  allowedRoles={["owner"]}
  requiredPermissions={["organizations.manage"]}
  logic="AND"
  loadingFallback={<div>Verifying access...</div>}
  errorFallback="Access verification failed"
>
  <div>Owner with org management access</div>
</CombinedConditionalRender>;
```

### Convenience Components

#### AdminOnly

Shows content only to admin/owner users.

```tsx
import { AdminOnly } from "@/components/auth";

<AdminOnly fallback={<div>Admin access required</div>}>
  <div>Admin panel content</div>
</AdminOnly>

<AdminOnly
  includeTechnicianAdmins={true}
  fallback={<div>Admin or technician admin required</div>}
>
  <div>Admin or technician admin content</div>
</AdminOnly>
```

#### ManagementOnly

Shows content to users with management permissions.

```tsx
import { ManagementOnly } from "@/components/auth";

<ManagementOnly fallback={<div>Management access required</div>}>
  <div>Management features</div>
</ManagementOnly>;
```

#### ReadOnly

Shows content to users with read permissions.

```tsx
import { ReadOnly } from "@/components/auth";

<ReadOnly fallback={<div>Read access required</div>}>
  <div>Read-only content</div>
</ReadOnly>;
```

## Utility Hooks

### usePermissions

Basic permission checking hook.

```tsx
import { usePermissions } from "@/hooks/use-permissions";

function MyComponent() {
  const {
    hasRole,
    hasAnyRole,
    hasPermission,
    hasAnyPermission,
    isOwner,
    isTechnician,
  } = usePermissions();

  if (isOwner()) {
    return <div>Owner content</div>;
  }

  if (hasPermission("users.manage")) {
    return <div>User management content</div>;
  }

  return <div>Default content</div>;
}
```

### useVisibility

Advanced visibility logic hook with comprehensive permission checks.

```tsx
import { useVisibility } from "@/hooks/use-visibility";

function MyComponent() {
  const {
    canManageUsers,
    canManageRoutes,
    canViewAnalytics,
    canExportData,
    canAccessAdmin,
    hasElevatedPrivileges,
    canPerformAction,
    canManage,
    canRead,
    canCreate,
    canUpdate,
    canDelete,
  } = useVisibility();

  if (canManageUsers()) {
    return <div>User management interface</div>;
  }

  if (canPerformAction("create", "routes")) {
    return <div>Route creation interface</div>;
  }

  return <div>Default interface</div>;
}
```

## Permission System

### Permission Format

Permissions follow the format: `resource.action`

Examples:

- `users.read` - Read user data
- `users.create` - Create new users
- `users.update` - Update user data
- `users.delete` - Delete users
- `users.manage` - Full user management
- `users.*` - All user permissions (wildcard)

### Default Role Permissions

#### Owner Role

```typescript
["organizations.*", "users.*", "technicians.*", "routes.*", "roles.*"];
```

#### Technician Role

```typescript
[
  "routes.read",
  "routes.update_status",
  "technicians.read_own",
  "technicians.update_own",
];
```

### Wildcard Permissions

The system supports wildcard permissions:

- `users.*` matches `users.read`, `users.create`, `users.update`, etc.
- `*.read` matches `users.read`, `routes.read`, `technicians.read`, etc.

## Best Practices

### 1. Use Appropriate Components

- Use `RoleGuard`/`PermissionGuard` for simple cases
- Use `RoleVisibility`/`PermissionVisibility` for advanced features
- Use `RoleConditionalRender`/`PermissionConditionalRender` when loading states matter
- Use convenience components for common patterns

### 2. Provide Meaningful Fallbacks

```tsx
// Good
<RoleGuard
  allowedRoles={["owner"]}
  fallback={<div>Contact your administrator for access</div>}
>
  <div>Admin content</div>
</RoleGuard>

// Avoid
<RoleGuard allowedRoles={["owner"]}>
  <div>Admin content</div>
</RoleGuard>
```

### 3. Handle Loading States

```tsx
// Good
<RoleConditionalRender
  allowedRoles={["owner"]}
  loadingFallback={<div>Verifying access...</div>}
>
  <div>Admin content</div>
</RoleConditionalRender>
```

### 4. Use Specific Permissions

```tsx
// Good
<PermissionGuard requiredPermissions={["users.manage"]}>
  <div>User management</div>
</PermissionGuard>

// Avoid
<RoleGuard allowedRoles={["owner"]}>
  <div>User management</div>
</RoleGuard>
```

### 5. Combine Logic Appropriately

```tsx
// Good - Clear intent
<CombinedVisibility
  allowedRoles={["owner"]}
  requiredPermissions={["organizations.manage"]}
  logic="AND"
>
  <div>Owner with org management</div>
</CombinedVisibility>

// Good - Flexible access
<CombinedVisibility
  allowedRoles={["owner"]}
  requiredPermissions={["users.manage"]}
  logic="OR"
>
  <div>Owner or user manager</div>
</CombinedVisibility>
```

## Performance Considerations

### 1. Memoization

The hooks use `useMemo` to prevent unnecessary recalculations.

### 2. Lazy Loading

Use `SuspenseConditionalRender` for components that need to load additional data.

### 3. Conditional Imports

Consider lazy loading components that are only used by specific roles.

```tsx
const AdminPanel = lazy(() => import("./AdminPanel"));

<RoleGuard allowedRoles={["owner"]}>
  <Suspense fallback={<div>Loading admin panel...</div>}>
    <AdminPanel />
  </Suspense>
</RoleGuard>;
```

## Error Handling

### 1. Authentication Errors

Components handle authentication errors gracefully:

- Show loading states while checking authentication
- Provide fallback content on errors
- Allow custom error messages

### 2. Permission Errors

- Clear fallback messages for permission denials
- Log permission failures for debugging
- Provide guidance on how to gain access

### 3. Network Errors

- Retry logic for permission checks
- Graceful degradation when services are unavailable
- Clear error messages for users

## Testing

### 1. Unit Tests

Test individual components with different user states:

```tsx
// Test with owner user
<RoleGuard allowedRoles={["owner"]}>
  <div>Owner content</div>
</RoleGuard>

// Test with technician user
<RoleGuard allowedRoles={["technician"]}>
  <div>Technician content</div>
</RoleGuard>
```

### 2. Integration Tests

Test complete flows with different user roles:

```tsx
// Test admin dashboard access
test("admin can access admin dashboard", () => {
  // Setup admin user
  // Navigate to admin page
  // Verify admin content is visible
  // Verify technician content is hidden
});
```

### 3. E2E Tests

Test real user scenarios:

```tsx
// Test complete user journey
test("technician workflow", () => {
  // Login as technician
  // Verify technician dashboard
  // Verify admin features are hidden
  // Test route management
});
```

## Migration Guide

### From Basic Guards

If you're using the basic `RoleGuard` and `PermissionGuard`:

```tsx
// Old
<RoleGuard allowedRoles={["owner"]}>
  <div>Content</div>
</RoleGuard>

// New - Same functionality
<RoleVisibility allowedRoles={["owner"]}>
  <div>Content</div>
</RoleVisibility>

// New - With loading states
<RoleConditionalRender allowedRoles={["owner"]}>
  <div>Content</div>
</RoleConditionalRender>
```

### From Custom Logic

If you have custom visibility logic:

```tsx
// Old
{user?.role === "owner" && (
  <div>Owner content</div>
)}

// New
<RoleVisibility allowedRoles={["owner"]}>
  <div>Owner content</div>
</RoleVisibility>

// Or for complex logic
<ConditionalVisibility
  condition={() => user?.role === "owner" && user?.active}
>
  <div>Active owner content</div>
</ConditionalVisibility>
```

## Troubleshooting

### Common Issues

1. **Component not showing/hiding correctly**
   - Check user role and permissions
   - Verify permission format (resource.action)
   - Check for typos in role names

2. **Loading states not working**
   - Use `*ConditionalRender` components instead of `*Visibility`
   - Ensure authentication state is properly managed

3. **Performance issues**
   - Use `useMemo` for complex permission calculations
   - Consider lazy loading for large components
   - Avoid nested visibility components

4. **TypeScript errors**
   - Ensure proper imports from `@/components/auth`
   - Check that role types match `UserRole` type
   - Verify permission strings match expected format

### Debug Tools

Use the `useVisibility` hook to debug permission issues:

```tsx
function DebugComponent() {
  const visibility = useVisibility();

  console.log("Permission checks:", {
    canManageUsers: visibility.canManageUsers(),
    canManageRoutes: visibility.canManageRoutes(),
    hasElevatedPrivileges: visibility.hasElevatedPrivileges(),
  });

  return <div>Debug info in console</div>;
}
```

## Conclusion

The role-based component visibility system provides a comprehensive, flexible, and performant solution for controlling component access based on user roles and permissions. By following the best practices outlined in this document, you can create secure, user-friendly interfaces that adapt to different user capabilities.

For more examples, see the `VisibilityExamples` component in `frontend/src/components/examples/visibility-examples.tsx`.
