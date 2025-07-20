# Authentication Endpoints API Documentation

This document provides comprehensive documentation for the authentication endpoints in the RoutrApp API.

## Overview

The authentication system provides three main endpoints for user authentication:

- **Login**: Authenticate users and receive JWT tokens
- **Logout**: Invalidate user session and clear refresh tokens
- **Refresh Token**: Obtain new access tokens using refresh tokens

All endpoints follow RESTful conventions and return consistent JSON responses with proper error handling and security measures.

## Base URL

```
http://localhost:8080/api/v1/auth
```

---

## Endpoints

### 1. Login

Authenticate a user with email and password to receive access and refresh tokens.

#### Request

```http
POST /api/v1/auth/login
Content-Type: application/json
```

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

#### Request Schema

| Field      | Type   | Required | Validation                        | Description          |
| ---------- | ------ | -------- | --------------------------------- | -------------------- |
| `email`    | string | Yes      | Valid email format, max 100 chars | User's email address |
| `password` | string | Yes      | Minimum 8 characters              | User's password      |

#### Success Response

**Status Code:** `200 OK`

```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "active": true,
      "role": "owner",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 900
  },
  "message": "Login successful"
}
```

#### Response Schema

| Field                  | Type    | Description                            |
| ---------------------- | ------- | -------------------------------------- |
| `success`              | boolean | Always true for successful responses   |
| `data.user`            | object  | User information object                |
| `data.user.id`         | integer | Unique user identifier                 |
| `data.user.email`      | string  | User's email address                   |
| `data.user.first_name` | string  | User's first name                      |
| `data.user.last_name`  | string  | User's last name                       |
| `data.user.active`     | boolean | Whether user account is active         |
| `data.user.role`       | string  | User's role ("owner" or "technician")  |
| `data.user.created_at` | string  | ISO 8601 timestamp of account creation |
| `data.user.updated_at` | string  | ISO 8601 timestamp of last update      |
| `data.access_token`    | string  | JWT access token (15 minutes expiry)   |
| `data.refresh_token`   | string  | JWT refresh token (7 days expiry)      |
| `data.token_type`      | string  | Token type, always "Bearer"            |
| `data.expires_in`      | integer | Access token expiry time in seconds    |
| `message`              | string  | Success message                        |

#### Error Responses

**400 Bad Request - Validation Error**

```json
{
  "error": {
    "status": 400,
    "message": "Invalid request data: validation failed",
    "details": {
      "code": "VALIDATION_ERROR"
    }
  }
}
```

**401 Unauthorized - Invalid Credentials**

```json
{
  "error": {
    "status": 401,
    "message": "Invalid credentials",
    "details": {
      "code": "INVALID_CREDENTIALS"
    }
  }
}
```

**401 Unauthorized - Account Disabled**

```json
{
  "error": {
    "status": 401,
    "message": "Account is disabled",
    "details": {
      "code": "ACCOUNT_DISABLED"
    }
  }
}
```

**500 Internal Server Error**

```json
{
  "error": {
    "status": 500,
    "message": "Failed to generate access token",
    "details": {
      "code": "TOKEN_GENERATION_ERROR"
    }
  }
}
```

---

### 2. Logout

Invalidate the user's session by clearing their refresh token from the database.

#### Request

```http
POST /api/v1/auth/logout
Authorization: Bearer {access_token}
Content-Type: application/json
```

#### Headers

| Header          | Required | Description                          |
| --------------- | -------- | ------------------------------------ |
| `Authorization` | Yes      | Bearer token with valid access token |

#### Success Response

**Status Code:** `200 OK`

```json
{
  "success": true,
  "message": "Logout successful"
}
```

#### Error Responses

**401 Unauthorized - Missing Authorization**

```json
{
  "error": {
    "status": 401,
    "message": "Authorization header is required",
    "details": {
      "code": "MISSING_AUTH_HEADER"
    }
  }
}
```

**401 Unauthorized - Invalid Token**

```json
{
  "error": {
    "status": 401,
    "message": "Invalid or expired token",
    "details": {
      "code": "INVALID_TOKEN"
    }
  }
}
```

**401 Unauthorized - Wrong Token Type**

```json
{
  "error": {
    "status": 401,
    "message": "Access token required",
    "details": {
      "code": "INVALID_TOKEN_TYPE"
    }
  }
}
```

**500 Internal Server Error**

```json
{
  "error": {
    "status": 500,
    "message": "Failed to logout",
    "details": {
      "code": "LOGOUT_ERROR"
    }
  }
}
```

---

### 3. Refresh Token

Obtain a new access token using a valid refresh token. The system first validates the token, then checks if the user is still active before generating a new access token.

#### Request

```http
POST /api/v1/auth/refresh
Content-Type: application/json
```

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Request Schema

| Field           | Type   | Required | Description             |
| --------------- | ------ | -------- | ----------------------- |
| `refresh_token` | string | Yes      | Valid JWT refresh token |

#### Success Response

**Status Code:** `200 OK`

```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 900
  },
  "message": "Token refreshed successfully"
}
```

#### Response Schema

| Field               | Type    | Description                              |
| ------------------- | ------- | ---------------------------------------- |
| `success`           | boolean | Always true for successful responses     |
| `data.access_token` | string  | New JWT access token (15 minutes expiry) |
| `data.token_type`   | string  | Token type, always "Bearer"              |
| `data.expires_in`   | integer | Access token expiry time in seconds      |
| `message`           | string  | Success message                          |

#### Error Responses

**400 Bad Request - Validation Error**

```json
{
  "error": {
    "status": 400,
    "message": "Invalid request data: validation failed",
    "details": {
      "code": "VALIDATION_ERROR"
    }
  }
}
```

**401 Unauthorized - Invalid Refresh Token**

```json
{
  "error": {
    "status": 401,
    "message": "Invalid refresh token",
    "details": {
      "code": "INVALID_REFRESH_TOKEN"
    }
  }
}
```

**401 Unauthorized - Wrong Token Type**

```json
{
  "error": {
    "status": 401,
    "message": "Invalid token type",
    "details": {
      "code": "INVALID_TOKEN_TYPE"
    }
  }
}
```

**401 Unauthorized - Account Disabled**

```json
{
  "error": {
    "status": 401,
    "message": "Account is disabled",
    "details": {
      "code": "ACCOUNT_DISABLED"
    }
  }
}
```

---

## Authentication Flow

### Standard Login Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant A as Auth API
    participant D as Database

    C->>A: POST /auth/login (email, password)
    A->>D: Find user by email
    D-->>A: User data
    A->>A: Check user is active
    A->>A: Verify password (bcrypt)
    A->>A: Generate access & refresh tokens
    A->>D: Update user with refresh token & last login
    D-->>A: Success
    A-->>C: Login response with tokens

    Note over C: Store tokens securely
    C->>A: API Request with Authorization: Bearer {access_token}
    A->>A: Validate access token
    A-->>C: Protected resource response
```

### Token Refresh Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant A as Auth API
    participant D as Database

    Note over C: Access token expired
    C->>A: POST /auth/refresh (refresh_token)
    A->>A: Validate refresh token JWT
    A->>A: Check token type is "refresh"
    A->>D: Find user by user_id
    D-->>A: User data
    A->>A: Check user is still active
    A->>A: Verify refresh token matches stored token
    A->>A: Generate new access token
    A-->>C: New access token

    Note over C: Update stored access token
    C->>A: API Request with new access token
    A-->>C: Protected resource response
```

### Logout Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant A as Auth API
    participant D as Database

    C->>A: POST /auth/logout (Authorization: Bearer {access_token})
    A->>A: Validate access token via middleware
    A->>A: Extract user_id from token claims
    A->>D: Clear refresh token for user
    D-->>A: Success
    A-->>C: Logout successful

    Note over C: Clear stored tokens
```

---

## Security Considerations

### Password Security

- Passwords are hashed using bcrypt with cost factor 12
- Minimum password length: 8 characters
- Password validation occurs at both request and business logic levels
- Plain text passwords are never stored or logged

### Token Security

- **Access tokens**: Short-lived (15 minutes) to limit exposure
- **Refresh tokens**: Long-lived (7 days) but stored securely in database
- **Token validation**: HMAC-SHA256 signature verification
- **Token revocation**: Refresh tokens are cleared on logout
- **Token type validation**: Prevents access tokens from being used as refresh tokens

### Multi-Tenant Security

- All tokens include `organization_id` for tenant isolation
- User access is automatically scoped to their organization
- Cross-tenant access is prevented at the token level
- Organization context is validated on every request

### Session Management

- **Last login tracking**: User's last login time is updated on successful authentication
- **Concurrent sessions**: Multiple devices can be logged in simultaneously
- **Session invalidation**: Logout immediately invalidates all refresh tokens
- **Token storage**: Refresh tokens are stored securely in the database

### Rate Limiting

- Consider implementing rate limiting on auth endpoints
- Recommended: 5 login attempts per minute per IP
- Recommended: 10 refresh attempts per minute per user
- Recommended: 20 logout attempts per minute per user

---

## Common Error Codes

| Code                      | HTTP Status | Description                        | Resolution                                          |
| ------------------------- | ----------- | ---------------------------------- | --------------------------------------------------- |
| `VALIDATION_ERROR`        | 400         | Request validation failed          | Check request format and required fields            |
| `INVALID_CREDENTIALS`     | 401         | Email/password combination invalid | Verify credentials                                  |
| `ACCOUNT_DISABLED`        | 401         | User account is disabled           | Contact administrator                               |
| `MISSING_AUTH_HEADER`     | 401         | Authorization header missing       | Include Bearer token                                |
| `INVALID_AUTH_HEADER`     | 401         | Authorization header malformed     | Use format: "Bearer {token}"                        |
| `INVALID_TOKEN`           | 401         | Token validation failed            | Refresh or re-authenticate                          |
| `INVALID_TOKEN_TYPE`      | 401         | Wrong token type used              | Use access token for API, refresh token for refresh |
| `INVALID_REFRESH_TOKEN`   | 401         | Refresh token invalid/expired      | Re-authenticate                                     |
| `AUTHENTICATION_REQUIRED` | 401         | Authentication required            | Include valid access token                          |
| `TOKEN_GENERATION_ERROR`  | 500         | Server error generating tokens     | Retry or contact support                            |
| `LOGOUT_ERROR`            | 500         | Server error during logout         | Retry or contact support                            |
| `INTERNAL_ERROR`          | 500         | Internal server error              | Contact support                                     |

---

## Frontend Integration Examples

### React/TypeScript Example

```typescript
interface LoginRequest {
  email: string;
  password: string;
}

interface LoginResponse {
  success: boolean;
  data: {
    user: User;
    access_token: string;
    refresh_token: string;
    token_type: string;
    expires_in: number;
  };
  message: string;
}

class AuthService {
  private baseUrl = "http://localhost:8080/api/v1/auth";

  async login(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await fetch(`${this.baseUrl}/login`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(credentials),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error?.message || "Login failed");
    }

    return response.json();
  }

  async logout(accessToken: string): Promise<void> {
    const response = await fetch(`${this.baseUrl}/logout`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error?.message || "Logout failed");
    }
  }

  async refreshToken(refreshToken: string): Promise<{ access_token: string }> {
    const response = await fetch(`${this.baseUrl}/refresh`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error?.message || "Token refresh failed");
    }

    const result = await response.json();
    return result.data;
  }
}
```

### cURL Examples

**Login:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

**Logout:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Refresh Token:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }'
```

---

## Testing

### Manual Testing Checklist

#### Login Endpoint

- [ ] Valid credentials return success with tokens
- [ ] Invalid email returns validation error
- [ ] Invalid password format returns validation error
- [ ] Non-existent user returns invalid credentials
- [ ] Wrong password returns invalid credentials
- [ ] Inactive user returns account disabled
- [ ] Database error returns internal server error
- [ ] Last login time is updated on successful login

#### Logout Endpoint

- [ ] Valid access token clears refresh token
- [ ] Missing authorization header returns error
- [ ] Invalid token format returns error
- [ ] Expired token returns error
- [ ] Refresh token used instead of access token returns error
- [ ] Database error returns internal server error

#### Refresh Token Endpoint

- [ ] Valid refresh token returns new access token
- [ ] Invalid refresh token returns error
- [ ] Expired refresh token returns error
- [ ] Access token used instead of refresh token returns error
- [ ] Token for inactive user returns account disabled
- [ ] Token not in database returns error
- [ ] User not found returns error

### Automated Testing

See the integration tests in `backend/internal/api/auth_test.go` for comprehensive automated testing coverage.

**Test Coverage Includes:**

- ✅ All success scenarios
- ✅ All validation error cases
- ✅ All authentication error cases
- ✅ Security edge cases
- ✅ Concurrent login scenarios
- ✅ Full authentication flow testing
- ✅ Token invalidation after logout

---

## Changelog

| Version | Date       | Changes                                                                          |
| ------- | ---------- | -------------------------------------------------------------------------------- |
| 1.1.0   | 2024-07-20 | Improved refresh token logic, added user status checking before token validation |
| 1.0.0   | 2024-01-15 | Initial implementation of auth endpoints                                         |

**Recent Changes (v1.1.0):**

- Enhanced refresh token validation to check user status first
- Improved error handling for inactive users during token refresh
- Added comprehensive test coverage for all scenarios
- Fixed concurrent login test isolation issues
- Updated error codes and response formats for consistency

---

## Support

For issues or questions regarding the authentication endpoints:

1. Check the [troubleshooting section](jwt_authentication.md#troubleshooting) in the JWT documentation
2. Review the error codes and their resolutions above
3. Check the server logs for detailed error information
4. Ensure your JWT secret is properly configured in production
5. Verify that all required environment variables are set correctly
