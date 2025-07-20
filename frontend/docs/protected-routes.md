# Protected Route System

This document describes the comprehensive protected route system implemented in the frontend application. The system provides multiple layers of authentication and authorization controls to secure your application.

## Overview

The protected route system consists of several components that work together to provide:

- **Authentication checking** - Verify users are logged in
- **Role-based access control (RBAC)** - Control access based on user roles
- **Permission-based access control** - Fine-grained permissions for specific actions
- **Automatic redirects** - Seamless user experience with smart routing
- **Error handling** - Graceful handling of auth errors
- **Loading states** - Smooth UI during authentication checks

## Components

### 1. ProtectedRoute

The core component for route protection with comprehensive options.

```tsx
import { ProtectedRoute } from "@/components/auth";

<ProtectedRoute
  requireAuth={true} // Require authentication
  allowedRoles={["owner", "technician"]} // Allow specific roles
  requiredPermissions={["routes.read"]} // Require specific permissions
  redirectTo="/login" // Where to redirect unauthorized users
  fallback={CustomFallbackComponent} // Custom fallback component
>
  <YourProtectedContent />
</ProtectedRoute>;
```

### 2. Page Guards

Convenient wrappers for entire pages with built-in layouts.

```tsx
import { OwnerPage, TechnicianPage, AuthenticatedPage } from "@/components/auth";

// Owner-only page
<OwnerPage title="Admin Dashboard" description="Manage your organization">
  <AdminContent />
</OwnerPage>

// Technician-only page
<TechnicianPage title="Route Management" description="Manage your routes">
  <TechnicianContent />
</TechnicianPage>

// Any authenticated user
<AuthenticatedPage title="Profile" description="Manage your profile">
  <ProfileContent />
</AuthenticatedPage>
```

### 3. Component-Level Guards

For conditional rendering within components.

```tsx
import { RoleGuard, PermissionGuard, OwnerOnly, TechnicianOnly } from "@/components/auth";

// Role-based rendering
<RoleGuard allowedRoles={["owner"]}>
  <AdminOnlyButton />
</RoleGuard>

// Permission-based rendering
<PermissionGuard requiredPermissions={["routes.create"]}>
  <CreateRouteButton />
</PermissionGuard>

// Convenience components
<OwnerOnly>
  <DeleteUserButton />
</OwnerOnly>

<TechnicianOnly fallback={<p>Technicians only</p>}>
  <UpdateStatusButton />
</TechnicianOnly>
```

### 4. Route Middleware

Global middleware components for application-wide protection.

```tsx
import { CombinedMiddleware } from "@/components/auth";

// In your root layout
<CombinedMiddleware
  routePermissions={routePermissionMap}
  enableSessionTimeout={true}
  enableRoleRedirect={true}
>
  {children}
</CombinedMiddleware>;
```

### 5. Error Boundary

Handles authentication errors gracefully.

```tsx
import { AuthErrorBoundary } from "@/components/auth";

<AuthErrorBoundary
  onError={error => console.error("Auth error:", error)}
  onLogout={() => authService.logout()}
>
  <YourApp />
</AuthErrorBoundary>;
```

## Roles and Permissions

### Available Roles

- **owner** - Full access to all organization features
- **technician** - Limited access to route and technician features

### Permission System

The system uses a hierarchical permission structure with wildcard support:

#### Owner Permissions

```
organizations.*  // All organization operations
users.*         // All user operations
technicians.*   // All technician operations
routes.*        // All route operations
roles.*         // All role operations
```

#### Technician Permissions

```
routes.read           // View routes
routes.update_status  // Update route status
technicians.read_own  // View own profile
technicians.update_own // Update own profile
```

### Wildcard Matching

Permissions support wildcard matching:

- `routes.*` matches `routes.read`, `routes.create`, `routes.update`, etc.
- `routes.read` only matches exact permission

## Usage Examples

### Basic Authentication Check

```tsx
import { ProtectedRoute } from "@/components/auth";

export default function Dashboard() {
  return (
    <ProtectedRoute>
      <h1>Dashboard</h1>
      <p>Only authenticated users can see this.</p>
    </ProtectedRoute>
  );
}
```

### Role-Based Page Protection

```tsx
import { OwnerPage } from "@/components/auth";

export default function AdminPanel() {
  return (
    <OwnerPage title="Admin Panel">
      <UserManagement />
      <SystemSettings />
    </OwnerPage>
  );
}
```

### Permission-Based Feature Access

```tsx
import { PermissionGuard } from "@/components/auth";
import { Button } from "@/components/ui/button";

export default function RouteList() {
  return (
    <div>
      <h1>Routes</h1>
      <RouteTable />

      <PermissionGuard requiredPermissions={["routes.create"]}>
        <Button>Create New Route</Button>
      </PermissionGuard>
    </div>
  );
}
```

### Conditional UI Based on Role

```tsx
import { RoleGuard } from "@/components/auth";

export default function Navigation() {
  return (
    <nav>
      <Link href="/dashboard">Dashboard</Link>

      <RoleGuard allowedRoles={["owner"]}>
        <Link href="/admin">Admin Panel</Link>
      </RoleGuard>

      <RoleGuard allowedRoles={["technician"]}>
        <Link href="/technician">My Routes</Link>
      </RoleGuard>
    </nav>
  );
}
```

### Custom Error Handling

```tsx
import { ProtectedRoute } from "@/components/auth";

function CustomUnauthorized() {
  return (
    <div>
      <h1>Oops!</h1>
      <p>You need special permissions for this page.</p>
    </div>
  );
}

export default function SpecialPage() {
  return (
    <ProtectedRoute allowedRoles={["owner"]} fallback={CustomUnauthorized}>
      <VeryImportantContent />
    </ProtectedRoute>
  );
}
```

## Hook Usage

### usePermissions Hook

```tsx
import { usePermissions } from "@/hooks/use-permissions";

export default function MyComponent() {
  const {
    hasRole,
    hasAnyRole,
    hasPermission,
    hasAnyPermission,
    isOwner,
    isTechnician,
  } = usePermissions();

  if (isOwner()) {
    return <AdminView />;
  }

  if (isTechnician()) {
    return <TechnicianView />;
  }

  return <GuestView />;
}
```

### useAuthErrorBoundary Hook

```tsx
import { useAuthErrorBoundary } from "@/components/auth";

export default function ApiComponent() {
  const { throwAuthError } = useAuthErrorBoundary();

  const handleApiCall = async () => {
    try {
      await api.someProtectedCall();
    } catch (error) {
      if (error.status === 403) {
        throwAuthError({
          type: "FORBIDDEN",
          message: "Access denied to this resource",
          statusCode: 403,
        });
      }
    }
  };

  return <button onClick={handleApiCall}>Make API Call</button>;
}
```

## Configuration

### Route Permissions Map

Define route-specific permissions in your layout:

```tsx
const routePermissions: Record<string, string[]> = {
  "/admin": ["organizations.*"],
  "/admin/users": ["users.*"],
  "/admin/technicians": ["technicians.*"],
  "/admin/routes": ["routes.*"],
  "/technician": ["routes.read"],
  "/technician/routes": ["routes.read"],
};

<CombinedMiddleware routePermissions={routePermissions}>
  {children}
</CombinedMiddleware>;
```

### Public Routes

Define which routes don't require authentication:

```tsx
// In AuthMiddleware component
const publicRoutes = ["/", "/login", "/register", "/about"];
const guestOnlyRoutes = ["/login", "/register"];
```

## Best Practices

1. **Use Page Guards for entire pages** - Cleaner and more maintainable than wrapping everything in ProtectedRoute
2. **Use Component Guards for UI elements** - Hide/show buttons and features based on permissions
3. **Combine role and permission checks** - Use roles for broad access, permissions for specific features
4. **Handle loading states** - Always provide good UX during authentication checks
5. **Provide fallbacks** - Custom error messages for better user experience
6. **Test all access levels** - Ensure both positive and negative test cases work

## Error States

The system handles various error states:

- **Unauthorized (401)** - User needs to log in
- **Forbidden (403)** - User doesn't have permission
- **Token Expired** - Session has expired
- **Unknown** - Generic error fallback

Each error type has appropriate UI and recovery options.

## Integration with Backend

The frontend permission system mirrors the backend RBAC implementation:

- Role definitions match backend enum values
- Permission strings match backend permission constants
- Wildcard matching logic is consistent
- Error responses are handled appropriately

This ensures consistent behavior across the entire application stack.
