# Authentication Service Migration Guide

## Breaking Change: `authService.isAuthenticated()` is now Async

### Overview

The `authService.isAuthenticated()` method has been changed from synchronous to asynchronous to support more robust token validation and refresh capabilities. This is a **breaking change** that requires code updates.

### Before (Synchronous - DEPRECATED)

```typescript
// ❌ This no longer works correctly
if (authService.isAuthenticated()) {
  // This code will always execute because isAuthenticated() now returns a Promise
  // The Promise object is truthy, even if the actual authentication status is false
  console.log("User is authenticated");
}

// ❌ This will assign a Promise object instead of a boolean
const isAuth = authService.isAuthenticated();
```

### After (Asynchronous - REQUIRED)

```typescript
// ✅ Correct async usage
if (await authService.isAuthenticated()) {
  console.log("User is authenticated");
}

// ✅ Or assign the awaited result
const isAuth = await authService.isAuthenticated();
if (isAuth) {
  console.log("User is authenticated");
}
```

## Migration Strategies

### 1. Standard Migration (Recommended)

Update your code to use `await` with the async method:

```typescript
// Before
function checkUserAccess() {
  if (authService.isAuthenticated()) {
    return "Access granted";
  }
  return "Access denied";
}

// After
async function checkUserAccess() {
  if (await authService.isAuthenticated()) {
    return "Access granted";
  }
  return "Access denied";
}
```

### 2. Temporary Sync Fallback (For Legacy Code)

If you cannot immediately convert to async, use the deprecated sync method:

```typescript
// ⚠️ Temporary solution only - plan to migrate to async
import { authService } from "@/services";

if (authService.isAuthenticatedSync()) {
  console.log("User might be authenticated (basic token check only)");
}
```

**Important**: `isAuthenticatedSync()` only checks if a token exists in storage. It does NOT validate the token or check if it's expired.

### 3. Migration Utility (Transition Helper)

Use the migration utility for gradual migration:

```typescript
import { authMigrationUtils } from "@/services";

// Async usage (recommended)
const isAuth = await authMigrationUtils.checkAuth();

// Sync usage (deprecated, for legacy code only)
const isAuth = authMigrationUtils.checkAuth(true);

// Helper for conditional logic
await authMigrationUtils.withAuth(async () => {
  console.log("This runs only if user is authenticated");
});
```

## Common Migration Patterns

### React Components

```typescript
// Before
function MyComponent() {
  const isAuth = authService.isAuthenticated(); // Returns Promise!

  if (isAuth) { // This always executes because Promise is truthy
    return <AuthenticatedView />;
  }
  return <LoginForm />;
}

// After - Option 1: Use useAuth hook (recommended)
function MyComponent() {
  const { isAuthenticated } = useAuth(); // Already handles async logic

  if (isAuthenticated) {
    return <AuthenticatedView />;
  }
  return <LoginForm />;
}

// After - Option 2: Use useEffect for async checking
function MyComponent() {
  const [isAuth, setIsAuth] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const checkAuth = async () => {
      const authStatus = await authService.isAuthenticated();
      setIsAuth(authStatus);
      setLoading(false);
    };
    checkAuth();
  }, []);

  if (loading) return <LoadingSpinner />;
  if (isAuth) return <AuthenticatedView />;
  return <LoginForm />;
}
```

### Route Guards

```typescript
// Before
function routeGuard(to, from, next) {
  if (authService.isAuthenticated()) {
    next();
  } else {
    next("/login");
  }
}

// After
async function routeGuard(to, from, next) {
  if (await authService.isAuthenticated()) {
    next();
  } else {
    next("/login");
  }
}
```

### Conditional API Calls

```typescript
// Before
function makeApiCall() {
  if (authService.isAuthenticated()) {
    return apiClient.get("/protected-data");
  }
  throw new Error("Not authenticated");
}

// After
async function makeApiCall() {
  if (await authService.isAuthenticated()) {
    return apiClient.get("/protected-data");
  }
  throw new Error("Not authenticated");
}
```

## Why This Change Was Made

1. **Token Validation**: The new async method validates tokens properly, including expiration checks
2. **Automatic Refresh**: Integrates with the token refresh system
3. **Better Security**: Ensures tokens are actually valid, not just present
4. **SSR Compatibility**: Works correctly in server-side rendering contexts
5. **Error Handling**: Provides better error handling for token-related issues

## TypeScript Updates

If you're using TypeScript, the method signatures have been updated:

```typescript
interface AuthService {
  // New signature
  isAuthenticated(): Promise<boolean>;

  // Temporary fallback
  isAuthenticatedSync(): boolean; // @deprecated
}
```

## Testing Considerations

Update your tests to handle the async nature:

```typescript
// Before
test("should check authentication", () => {
  const isAuth = authService.isAuthenticated();
  expect(isAuth).toBe(true);
});

// After
test("should check authentication", async () => {
  const isAuth = await authService.isAuthenticated();
  expect(isAuth).toBe(true);
});
```

## Timeline for Migration

1. **Immediate**: Use `isAuthenticatedSync()` for critical fixes
2. **Within 1 week**: Convert all synchronous usage to async
3. **Within 2 weeks**: Remove all usage of `isAuthenticatedSync()`
4. **Future**: The `isAuthenticatedSync()` method will be removed

## Need Help?

If you encounter issues during migration:

1. Check console warnings for specific guidance
2. Use the migration utilities as temporary bridges
3. Refer to existing components that use the `useAuth()` hook
4. Contact the development team for assistance

## Validation Checklist

After migration, verify:

- [ ] No console warnings about deprecated sync methods
- [ ] Authentication checks work correctly in all scenarios
- [ ] No TypeScript errors related to Promise/boolean mismatches
- [ ] Tests pass with async authentication checks
- [ ] User experience remains smooth during auth checks
