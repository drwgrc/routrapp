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

export interface AuthContextValue extends AuthState {
  login: (credentials: LoginCredentials) => Promise<void>;
  logout: () => Promise<void>;
  register: (data: RegistrationData) => Promise<void>;
  clearError: () => void;
  refreshUser: () => Promise<void>;
}

// Role-based access control types
export type UserRole = "owner" | "technician";

export interface Permission {
  action: string;
  resource: string;
}

// Route protection types
export interface RouteProtectionOptions {
  requireAuth?: boolean;
  allowedRoles?: UserRole[];
  requiredPermissions?: string[];
  redirectTo?: string;
  fallback?: React.ComponentType;
}

export interface ProtectedRouteProps extends RouteProtectionOptions {
  children: React.ReactNode;
}

// Permission checking utilities
export interface PermissionCheck {
  hasRole: (role: UserRole) => boolean;
  hasAnyRole: (roles: UserRole[]) => boolean;
  hasPermission: (permission: string) => boolean;
  hasAnyPermission: (permissions: string[]) => boolean;
  isOwner: () => boolean;
  isTechnician: () => boolean;
}

// Auth error types
export interface AuthError {
  type: "UNAUTHORIZED" | "FORBIDDEN" | "TOKEN_EXPIRED" | "UNKNOWN";
  message: string;
  statusCode?: number;
}
