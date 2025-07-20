// Authentication types for the frontend application

export interface User {
  id: number;
  email: string;
  first_name: string;
  last_name: string;
  active: boolean;
  role: string;
  created_at: string;
  updated_at: string;
}

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface RegistrationData {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
  organizationName: string;
  organizationEmail: string;
  subDomain: string;
}

export interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
}

// Enhanced auth response types
export interface LoginResponse {
  user: User;
  access_token: string;
  refresh_token: string;
  token_type: string;
  expires_in: number;
}

export interface TokenInfo {
  accessToken: string | null;
  isExpired: boolean;
  expiresAt: Date | null;
  timeUntilExpiry: number | null;
}

// Auth service response types
export interface AuthServiceLoginResponse {
  token: string; // access_token for backwards compatibility
  refreshToken: string;
  user: User;
}

// Enhanced auth context value with token management
export interface AuthContextValue extends AuthState {
  login: (credentials: LoginCredentials) => Promise<void>;
  logout: () => Promise<void>;
  register: (data: RegistrationData) => Promise<void>;
  clearError: () => void;
  refreshUser: () => Promise<void>;
}

// Extended auth context with token utilities
export interface ExtendedAuthContextValue extends AuthContextValue {
  getTokenInfo: () => Promise<TokenInfo>;
  refreshToken: () => Promise<boolean>;
}

// Token storage related types
export interface TokenStorage {
  setToken: (key: string, value: string) => Promise<void>;
  getToken: (key: string) => Promise<string | null>;
  removeToken: (key: string) => Promise<void>;
  clear: () => Promise<void>;
  isAvailable: () => boolean;
}

// Token manager types
export interface TokenManagerConfig {
  storage?: TokenStorage;
  refreshEndpoint?: string;
  refreshThresholdMs?: number;
  maxRetries?: number;
  retryDelayMs?: number;
  onTokenRefreshed?: (tokens: { accessToken: string; expiresAt: Date }) => void;
  onRefreshFailed?: (error: Error) => void;
  onTokenExpired?: () => void;
}

// JWT payload structure (for client-side parsing only)
export interface JWTPayload {
  sub: string; // user_id
  email: string;
  role: string;
  organization_id: string;
  exp: number; // expiry timestamp
  iat: number; // issued at timestamp
  type: "access" | "refresh";
}

// Auth service error types
export interface AuthError {
  message: string;
  code: string;
  status?: number;
  originalError?: unknown;
}

// API error response format
export interface APIErrorResponse {
  error: {
    status: number;
    message: string;
    details?: {
      code: string;
      [key: string]: unknown;
    };
  };
}

// Auth status for better state management
export type AuthStatus = "checking" | "authenticated" | "unauthenticated";

// Token refresh response from API
export interface TokenRefreshAPIResponse {
  success: boolean;
  data: {
    access_token: string;
    token_type: string;
    expires_in: number;
  };
  message: string;
}
