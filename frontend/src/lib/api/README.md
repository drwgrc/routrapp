# API Client

This directory contains the API client implementation for the RoutrApp frontend.

## Overview

The API client is built on Axios and provides a consistent interface for making API requests to the backend. It handles:

- Base URL configuration
- Authentication with JWT
- Request/response interceptors
- Error standardization
- Type-safe API methods

## Usage

### Basic Usage

```typescript
import apiClient from "../lib/api";

// Simple GET request
const userData = await apiClient.get("/users/me");

// POST request with data
const createResult = await apiClient.post("/routes", {
  name: "Monday Route",
  technicianId: 123,
});
```

### With Services

For most use cases, you should use the service layer which builds upon this API client:

```typescript
import { authService } from "../services";

// Login
await authService.login({
  email: "user@example.com",
  password: "password",
});

// Get current user
const user = await authService.getCurrentUser();
```

## Error Handling

All API errors are standardized to follow this format:

```typescript
{
  message: string;      // Human-readable error message
  status?: number;      // HTTP status code
  data?: unknown;       // Original response data
  originalError?: any;  // Original error object
}
```

Example error handling:

```typescript
try {
  const data = await apiClient.get("/routes");
  // Handle success
} catch (error: any) {
  if (error.status === 404) {
    // Handle not found
  } else {
    // Handle other errors
    console.error(error.message);
  }
}
```

## Architecture

The API client is structured in layers:

1. `axios-instance.ts` - Base Axios configuration
2. `interceptors.ts` - Request/response interceptors
3. `api-client.ts` - Generic API methods
4. `index.ts` - Export barrel

## Adding New Functionality

### Creating a New Service

When adding a new feature, create a service in `src/services/` following this pattern:

```typescript
// src/services/route-service.ts
import apiClient from "../lib/api";

interface Route {
  id: string;
  name: string;
  // ...other properties
}

const routeService = {
  getRoutes: async (): Promise<Route[]> => {
    return apiClient.get("/routes");
  },

  getRoute: async (id: string): Promise<Route> => {
    return apiClient.get(`/routes/${id}`);
  },

  createRoute: async (data: Omit<Route, "id">): Promise<Route> => {
    return apiClient.post("/routes", data);
  },

  // ... other methods
};

export default routeService;
```

Then add it to the service index:

```typescript
// src/services/index.ts
import routeService from "./route-service";

export {
  // ...other services
  routeService,
};
```

## Configuration

The API client uses environment variables for configuration:

- `NEXT_PUBLIC_API_BASE_URL` - Base API URL (defaults to 'http://localhost:8080/api')
- `NEXT_PUBLIC_API_VERSION` - API version (defaults to 'v1')

These can be set in `.env.local` or other appropriate environment files.

## Extending Interceptors

If you need to add custom interceptors, modify the `interceptors.ts` file:

```typescript
// Example: Adding a custom header to all requests
export const setupCustomInterceptors = (axiosInstance: AxiosInstance): void => {
  axiosInstance.interceptors.request.use(config => {
    config.headers["X-Custom-Header"] = "value";
    return config;
  });
};
```

Then update the initialization in `api-client.ts`.

## See Also

For more detailed documentation, see the [API Client Architecture](../../docs/architecture/frontend-api-client.md) document.
