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
  id: string;
  email: string;
  name: string;
  organizationId: string;
  role: string;
}

interface RegistrationData {
  email: string;
  password: string;
  name: string;
  organizationName: string;
}

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
      localStorage.setItem("auth_token", response.token);
      localStorage.setItem("refresh_token", response.refreshToken);

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
      localStorage.removeItem("auth_token");
      localStorage.removeItem("refresh_token");
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
          name: data.name,
          organizationName: data.organizationName,
        }
      );
      return response;
    } catch (error) {
      console.error("Registration failed:", error);
      throw error;
    }
  },

  // Check if user is authenticated
  isAuthenticated: (): boolean => {
    return !!localStorage.getItem("auth_token");
  },

  // Get current user data
  getCurrentUser: async (): Promise<UserData | null> => {
    try {
      if (!authService.isAuthenticated()) {
        return null;
      }

      const userData = await apiClient.get<UserData>("/auth/me");
      return userData;
    } catch (error) {
      console.error("Failed to get user data:", error);
      return null;
    }
  },
};

export default authService;
