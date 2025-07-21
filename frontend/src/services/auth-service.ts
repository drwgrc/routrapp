import apiClient from "../lib/api/api-client";

// Auth service types
interface LoginCredentials {
  email: string;
  password: string;
}

interface LoginResponse {
  token: string;
  refreshToken: string;
  user: UserData;
}

interface UserData {
  id: number;
  email: string;
  first_name: string;
  last_name: string;
  active: boolean;
  role: string;
  created_at: string;
  updated_at: string;
}

interface RegistrationData {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
  organizationName: string;
  organizationEmail: string;
  subDomain: string;
}

interface ProfileUpdateData {
  firstName?: string;
  lastName?: string;
}

// Helper function to safely access localStorage
const getFromStorage = (key: string): string | null => {
  if (typeof window === "undefined") return null;
  try {
    return localStorage.getItem(key);
  } catch {
    return null;
  }
};

const setToStorage = (key: string, value: string): void => {
  if (typeof window === "undefined") return;
  try {
    localStorage.setItem(key, value);
  } catch {
    // Silently fail if localStorage is not available
  }
};

const removeFromStorage = (key: string): void => {
  if (typeof window === "undefined") return;
  try {
    localStorage.removeItem(key);
  } catch {
    // Silently fail if localStorage is not available
  }
};

// Auth service implementation
const authService = {
  // User login
  login: async (credentials: LoginCredentials): Promise<LoginResponse> => {
    try {
      const response = await apiClient.post<LoginResponse>("/auth/login", {
        email: credentials.email,
        password: credentials.password,
      });

      // Store tokens in localStorage or secure storage
      setToStorage("auth_token", response.token);
      setToStorage("refresh_token", response.refreshToken);

      return response;
    } catch (error) {
      console.error("Login failed:", error);
      throw error;
    }
  },

  // User logout
  logout: async (): Promise<void> => {
    try {
      // Call logout endpoint to invalidate token on server
      await apiClient.post<{ success: boolean }>("/auth/logout", {});
    } catch (error) {
      console.error("Logout error:", error);
    } finally {
      // Clear local storage regardless of API call success
      removeFromStorage("auth_token");
      removeFromStorage("refresh_token");
    }
  },

  // User registration
  register: async (data: RegistrationData): Promise<{ success: boolean }> => {
    try {
      const response = await apiClient.post<{ success: boolean }>(
        "/auth/register",
        {
          email: data.email,
          password: data.password,
          first_name: data.firstName,
          last_name: data.lastName,
          organization_name: data.organizationName,
          organization_email: data.organizationEmail,
          sub_domain: data.subDomain,
        }
      );
      return response;
    } catch (error) {
      console.error("Registration failed:", error);
      throw error;
    }
  },

  // Update user profile
  updateProfile: async (data: ProfileUpdateData): Promise<UserData> => {
    try {
      const response = await apiClient.put<UserData>("/users/profile", {
        first_name: data.firstName,
        last_name: data.lastName,
      });
      return response;
    } catch (error) {
      console.error("Profile update failed:", error);
      throw error;
    }
  },

  // Check if user is authenticated
  isAuthenticated: (): boolean => {
    return !!getFromStorage("auth_token");
  },

  // Get current user data
  getCurrentUser: async (): Promise<UserData | null> => {
    try {
      if (!authService.isAuthenticated()) {
        return null;
      }

      const userData = await apiClient.get<UserData>("/auth/me");
      return userData as UserData;
    } catch (error) {
      // Don't log 401 errors as they are expected when tokens are invalid
      const isAxiosError =
        error && typeof error === "object" && "status" in error;
      const status = isAxiosError
        ? (error as { status?: number }).status
        : undefined;

      if (status !== 401) {
        console.error("Failed to get user data:", error);
      }

      // Don't clear tokens on 401 errors - let the user stay logged in
      // The token might be valid but there could be server issues
      // Only clear tokens for other types of errors
      if (status && status !== 401) {
        removeFromStorage("auth_token");
        removeFromStorage("refresh_token");
      }

      return null;
    }
  },
};

export default authService;
