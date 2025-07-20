// Authentication types for the frontend application

export interface User {
  id: string;
  email: string;
  name: string;
  organizationId: string;
  role: string;
}

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface RegistrationData {
  email: string;
  password: string;
  name: string;
  organizationName: string;
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
