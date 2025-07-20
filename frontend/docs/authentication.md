# Authentication Context Documentation

## Overview

The authentication system for the RouteApp frontend is built using React Context API combined with TanStack Query for optimal state management, caching, and server synchronization. This implementation provides a centralized authentication state that can be accessed throughout the application.

## Architecture

The authentication system consists of several key components:

1. **Authentication Types** (`src/types/auth.ts`) - TypeScript interfaces for type safety
2. **Query Client Setup** (`src/lib/query-client.ts`) - TanStack Query configuration
3. **Authentication Context** (`src/contexts/auth-context.tsx`) - Main context implementation
4. **App Providers** (`src/providers/app-providers.tsx`) - Provider wrapper for the application
5. **Authentication Service** (`src/services/auth-service.ts`) - API service layer

## Key Features

### 🚀 TanStack Query Integration

- **Automatic Caching**: User data is cached and automatically invalidated when needed
- **Background Refetching**: Keeps user data fresh without blocking the UI
- **Optimistic Updates**: Immediate UI updates with automatic rollback on errors
- **Smart Retries**: Automatic retry logic with exponential backoff for failed requests
- **Error Handling**: Centralized error management with automatic 401/403 handling

### 🔒 Security Features

- **Automatic Token Management**: Handles JWT token storage and refresh
- **Session Validation**: Validates user sessions on app initialization
- **Secure Logout**: Clears all cached data and tokens on logout
- **Error Recovery**: Automatic logout on authentication errors

### 🎯 Developer Experience

- **Type Safety**: Full TypeScript support with strict typing
- **Easy Integration**: Simple hook-based API for components
- **Predictable State**: Centralized state management with clear data flow
- **Error Boundaries**: Graceful error handling with user-friendly messages

## Usage

### 1. Setup Providers

Wrap your application with the `AppProviders` component:

```tsx
// app/layout.tsx
import { AppProviders } from "@/providers/app-providers";

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>
        <AppProviders>{children}</AppProviders>
      </body>
    </html>
  );
}
```

### 2. Using the Authentication Hook

```tsx
// components/LoginForm.tsx
import { useAuth } from "@/contexts";

export function LoginForm() {
  const { login, user, isLoading, error } = useAuth();

  const handleLogin = async (credentials: LoginCredentials) => {
    try {
      await login(credentials);
      // User is automatically updated in context
    } catch (error) {
      // Error is automatically handled and available in context
      console.error("Login failed:", error);
    }
  };

  if (isLoading) return <div>Loading...</div>;
  if (user) return <div>Welcome, {user.name}!</div>;

  return (
    <form onSubmit={handleLogin}>
      {error && <div className="error">{error}</div>}
      {/* Login form fields */}
    </form>
  );
}
```

### 3. Protected Routes

```tsx
// components/ProtectedRoute.tsx
import { useAuth } from "@/contexts";
import { redirect } from "next/navigation";

export function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) return <div>Loading...</div>;
  if (!isAuthenticated) redirect("/login");

  return <>{children}</>;
}
```

## API Reference

### AuthContextValue Interface

```typescript
interface AuthContextValue {
  // State
  user: User | null; // Current authenticated user
  isAuthenticated: boolean; // Authentication status
  isLoading: boolean; // Loading state for all operations
  error: string | null; // Current error message

  // Methods
  login: (credentials: LoginCredentials) => Promise<void>;
  logout: () => Promise<void>;
  register: (data: RegistrationData) => Promise<void>;
  clearError: () => void; // Clear current error
  refreshUser: () => Promise<void>; // Refresh user data from server
}
```

### Available Methods

#### `login(credentials: LoginCredentials)`

- Authenticates user with email/password
- Updates user state on success
- Handles token storage automatically
- Throws error on failure (also available in `error` state)

#### `logout()`

- Logs out current user
- Clears all cached data and tokens
- Redirects to login page (if implemented)

#### `register(data: RegistrationData)`

- Registers new user account
- Does not automatically log in user
- Returns success status

#### `refreshUser()`

- Invalidates and refetches current user data
- Useful after profile updates
- Handles authentication errors gracefully

#### `clearError()`

- Clears current error state
- Resets all mutation error states

## Configuration

### Query Client Settings

The authentication system uses the following TanStack Query configuration:

```typescript
{
  staleTime: 5 * 60 * 1000,        // 5 minutes - data considered fresh
  gcTime: 10 * 60 * 1000,          // 10 minutes - garbage collection
  refetchOnWindowFocus: false,      // Don't refetch on window focus
  refetchOnReconnect: true,         // Refetch on network reconnection
  retry: (failureCount, error) => { // Smart retry logic
    // Don't retry on auth errors (401/403)
    // Retry up to 3 times for other errors
  }
}
```

### Error Handling

The system automatically handles common authentication errors:

- **401 Unauthorized**: Automatically logs out user
- **403 Forbidden**: Automatically logs out user
- **Network Errors**: Retries up to 3 times with exponential backoff
- **Other Errors**: Displays user-friendly error messages

## Best Practices

### 1. Error Handling

```tsx
const { login, error, clearError } = useAuth();

useEffect(() => {
  if (error) {
    // Show toast notification
    toast.error(error);
    clearError();
  }
}, [error, clearError]);
```

### 2. Loading States

```tsx
const { isLoading, user } = useAuth();

if (isLoading) {
  return <LoadingSpinner />;
}
```

### 3. Conditional Rendering

```tsx
const { isAuthenticated, user } = useAuth();

return (
  <div>
    {isAuthenticated ? <DashboardLayout user={user} /> : <AuthenticationFlow />}
  </div>
);
```

## Migration Guide

If you're migrating from a basic Context implementation:

1. **Install TanStack Query**: `npm install @tanstack/react-query`
2. **Replace AuthProvider**: Use the new TanStack Query-powered provider
3. **Update imports**: Change import paths to use the new structure
4. **Add AppProviders**: Wrap your app with the new provider structure
5. **Test thoroughly**: Verify all authentication flows work correctly

## Performance Considerations

- **Automatic Caching**: User data is cached to reduce API calls
- **Background Updates**: Data is refreshed in the background without blocking UI
- **Optimistic Updates**: UI updates immediately for better user experience
- **Memory Management**: Old cached data is automatically garbage collected

## Troubleshooting

### Common Issues

1. **"useAuth must be used within an AuthProvider"**
   - Ensure your component is wrapped with `AppProviders`
   - Check the provider hierarchy in your app structure

2. **Authentication state not persisting**
   - Verify token storage in localStorage
   - Check network requests in browser DevTools
   - Ensure API endpoints are returning expected data

3. **Infinite loading states**
   - Check if `isAuthenticated()` method in auth service works correctly
   - Verify API endpoints are accessible
   - Check network connectivity

### Debug Mode

To enable debug mode for TanStack Query:

```typescript
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';

// Add to your AppProviders component
<ReactQueryDevtools initialIsOpen={false} />
```

## Future Enhancements

- **Refresh Token Rotation**: Implement automatic token refresh
- **Multi-Factor Authentication**: Add MFA support
- **Role-Based Access Control**: Enhanced permission system
- **Session Management**: Advanced session handling
- **Offline Support**: Cache authentication state for offline use

## Contributing

When making changes to the authentication system:

1. **Update Types**: Ensure all TypeScript interfaces are updated
2. **Test Thoroughly**: Test all authentication flows
3. **Update Documentation**: Keep this documentation current
4. **Consider Migration**: Provide migration guides for breaking changes
5. **Performance**: Consider caching and performance implications
