# Frontend API Client Architecture

## Overview

The frontend API client provides a consistent, type-safe interface for communicating with the backend API. It's built around Axios with added features for authentication, error handling, and request/response standardization.

## Structure

```
frontend/src/lib/api/
├── axios-instance.ts    # Base axios configuration
├── interceptors.ts      # Request/response interceptors
├── api-client.ts        # Generic API methods
└── index.ts             # Export barrel

frontend/src/services/
├── auth-service.ts      # Authentication service
└── index.ts             # Service exports
```

## Key Features

- JWT authentication with automatic token inclusion
- Token refresh handling on 401 responses
- Standardized error formatting
- Type-safe request/response handling
- Offline request queueing (planned)

## Core API Client

The API client follows a layered architecture:

1. **Base Configuration (`axios-instance.ts`)**:

   - Configures base URL from environment variables
   - Sets default headers and timeout
   - Uses versioned API endpoints (`/api/v1/`)

2. **Request/Response Interceptors (`interceptors.ts`)**:

   - Automatically injects JWT authentication tokens
   - Handles authentication errors and token refresh
   - Standardizes error responses

3. **Generic API Methods (`api-client.ts`)**:
   - Provides typed wrapper methods for HTTP verbs (GET, POST, etc.)
   - Handles response unwrapping
   - Enforces consistent error handling

## Error Handling

Follows the project's standardized error format:

```typescript
{
  message: string;
  status?: number;
  data?: unknown;
  originalError?: unknown;
}
```

Error responses from the backend are automatically parsed and formatted for consistent handling throughout the application.

## Authentication Flow

1. User credentials sent to `/auth/login` endpoint
2. JWT token stored in localStorage
3. Token automatically included in subsequent requests
4. 401 errors trigger token refresh or redirect to login

## Service Layer

API requests are organized into service modules:

- `authService`: Handles user authentication and session management
- Future services will be added for routes, technicians, etc.

## Usage Example

```typescript
import { authService } from "../services";

// Login with credentials
const handleLogin = async () => {
  try {
    const result = await authService.login({
      email: "user@example.com",
      password: "password",
    });
    // Handle success
  } catch (error) {
    // Handle error
  }
};

// Get current user data
const getUserData = async () => {
  const userData = await authService.getCurrentUser();
  if (userData) {
    // User is authenticated
  }
};
```

## Offline Support

The API client architecture will support offline operations:

- **Planned**: Request queueing for offline operations
- **Planned**: Background synchronization when connectivity is restored
- **Planned**: Conflict resolution for offline data changes

## Security Considerations

- Tokens are stored in localStorage (consider more secure alternatives in the future)
- All requests are made over HTTPS
- Authentication headers are automatically managed
- Sensitive data is not logged

## Future Enhancements

- Implement React Query integration for data fetching and caching
- Add request cancellation support
- Add request rate limiting
- Implement more sophisticated offline capabilities
- Consider using HTTP-only cookies for token storage
