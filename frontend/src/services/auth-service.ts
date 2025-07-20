import apiClient from "../lib/api/api-client";
import { defaultTokenManager } from "../lib/token-manager";

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

// Enhanced auth service implementation
const authService = {
  // User login
  login: async (credentials: LoginCredentials): Promise<LoginResponse> => {
    try {
      const response = await apiClient.post<{
        user: UserData;
        access_token: string;
        refresh_token: string;
        token_type: string;
        expires_in: number;
      }>("/auth/login", {
        email: credentials.email,
        password: credentials.password,
      });

      // Store tokens using token manager
      await defaultTokenManager.setTokens(
        response.access_token,
        response.refresh_token
      );

      // Return the expected format for backwards compatibility
      return {
        token: response.access_token,
        refreshToken: response.refresh_token,
        user: response.user,
      };
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
      // Continue with local cleanup even if server call fails
    } finally {
      // Clear tokens using token manager
      await defaultTokenManager.clearTokens();
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

  // Check if user is authenticated
  isAuthenticated: async (): Promise<boolean> => {
    try {
      return await defaultTokenManager.isAuthenticated();
    } catch (error) {
      console.warn("Error checking authentication status:", error);
      return false;
    }
  },

  // Get current user data
  getCurrentUser: async (): Promise<UserData | null> => {
    try {
      // Check if user is authenticated first
      const isAuth = await authService.isAuthenticated();
      if (!isAuth) {
        return null;
      }

      const userData = await apiClient.get<UserData>("/auth/me");
      return userData as UserData;
    } catch (error) {
      // Handle different error types appropriately
      const isAxiosError =
        error && typeof error === "object" && "status" in error;
      const status = isAxiosError
        ? (error as { status?: number }).status
        : undefined;

      // For 401 errors, token might be expired or invalid
      if (status === 401) {
        console.warn("Authentication required or token expired");
        // Let token manager handle token cleanup if needed
        return null;
      }

      // Don't log 401 errors as they are expected when tokens are invalid
      if (status !== 401) {
        console.error("Failed to get user data:", error);
      }

      return null;
    }
  },

  // Get access token (useful for debugging or manual API calls)
  getAccessToken: async (): Promise<string | null> => {
    try {
      return await defaultTokenManager.getAccessToken();
    } catch (error) {
      console.warn("Failed to get access token:", error);
      return null;
    }
  },

  // Get refresh token (useful for debugging)
  getRefreshToken: async (): Promise<string | null> => {
    try {
      return await defaultTokenManager.getRefreshToken();
    } catch (error) {
      console.warn("Failed to get refresh token:", error);
      return null;
    }
  },

  // Manually refresh token (useful for testing or explicit refresh)
  refreshToken: async (): Promise<boolean> => {
    try {
      const newToken = await defaultTokenManager.refreshTokenIfNeeded();
      return !!newToken;
    } catch (error) {
      console.error("Manual token refresh failed:", error);
      return false;
    }
  },

  // Get token information (useful for debugging)
  getTokenInfo: async () => {
    try {
      return await defaultTokenManager.getTokenInfo();
    } catch (error) {
      console.warn("Failed to get token info:", error);
      return {
        accessToken: null,
        isExpired: true,
        expiresAt: null,
        timeUntilExpiry: null,
      };
    }
  },

  // Clear all auth data (useful for debugging or manual cleanup)
  clearAuthData: async (): Promise<void> => {
    try {
      await defaultTokenManager.clearTokens();
    } catch (error) {
      console.warn("Failed to clear auth data:", error);
    }
  },
};

export default authService;
