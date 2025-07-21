import apiClient from "../lib/api/api-client";
import { defaultTokenManager } from "../lib/token-manager";
import type { AuthService } from "../types/auth";

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
const authService: AuthService = {
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

  // Check if user is authenticated (async - preferred method)
  isAuthenticated: async (): Promise<boolean> => {
    try {
      return await defaultTokenManager.isAuthenticated();
    } catch (error) {
      console.warn("Error checking authentication status:", error);
      return false;
    }
  },

  // DEPRECATED: Synchronous authentication check for backward compatibility
  // This method provides a fallback for code that hasn't been migrated to async yet
  // It performs a basic token existence check without validation
  // TODO: Remove this method after all code is migrated to use async isAuthenticated()
  isAuthenticatedSync: (): boolean => {
    console.warn(
      "authService.isAuthenticatedSync() is deprecated. " +
        "Please use 'await authService.isAuthenticated()' instead. " +
        "This synchronous method only checks token presence, not validity."
    );

    try {
      // Basic synchronous check - only verifies token exists, not if it's valid
      if (typeof window === "undefined") return false;

      const accessToken =
        localStorage.getItem("access_token") ||
        localStorage.getItem("auth_token");
      return !!accessToken;
    } catch (error) {
      console.warn("Error in synchronous auth check:", error);
      return false;
    }
  },

  // Get current user data
  getCurrentUser: async (): Promise<UserData | null> => {
    try {
      // Check if user is authenticated first
      const isAuth = await defaultTokenManager.isAuthenticated();
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

// Migration utility for transitioning from sync to async authentication checks
export const authMigrationUtils = {
  /**
   * Helper function to safely check authentication with fallback for legacy code
   * This function attempts async authentication first, then falls back to sync if in a non-async context
   *
   * @param preferSync - If true, uses the deprecated sync method (not recommended)
   * @returns Promise<boolean> for async contexts, boolean for sync contexts
   *
   * @example
   * // Preferred async usage:
   * const isAuth = await authMigrationUtils.checkAuth();
   *
   * // Legacy sync usage (not recommended):
   * const isAuth = authMigrationUtils.checkAuth(true);
   */
  checkAuth: (preferSync = false): boolean | Promise<boolean> => {
    if (preferSync) {
      console.warn(
        "Using synchronous authentication check. " +
          "Please migrate to async: await authService.isAuthenticated()"
      );
      return authService.isAuthenticatedSync();
    }

    return defaultTokenManager.isAuthenticated();
  },

  /**
   * Migration helper that wraps authentication logic to handle both sync and async patterns
   * This is a temporary helper to ease the migration process
   *
   * @param callback - Function to execute if authenticated
   * @param useAsync - Whether to use async authentication check (recommended)
   */
  withAuth: async (
    callback: () => void | Promise<void>,
    useAsync = true
  ): Promise<void> => {
    let isAuth: boolean;

    if (useAsync) {
      isAuth = await defaultTokenManager.isAuthenticated();
    } else {
      console.warn("Using deprecated sync auth check in withAuth helper");
      isAuth = authService.isAuthenticatedSync();
    }

    if (isAuth) {
      await callback();
    } else {
      console.warn("User not authenticated - callback not executed");
    }
  },
};
