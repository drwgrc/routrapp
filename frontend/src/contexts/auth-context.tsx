"use client";

import React, {
  createContext,
  useContext,
  ReactNode,
  useState,
  useEffect,
} from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import authService from "@/services/auth-service";
import { queryKeys } from "@/lib/query-client";
import { defaultTokenManager } from "@/lib/token-manager";
import {
  AuthContextValue,
  LoginCredentials,
  RegistrationData,
  User,
} from "@/types/auth";

// Create authentication context
const AuthContext = createContext<AuthContextValue | undefined>(undefined);

// Authentication provider props
interface AuthProviderProps {
  children: ReactNode;
}

// Authentication provider component
export function AuthProvider({ children }: AuthProviderProps) {
  const queryClient = useQueryClient();
  const [isMounted, setIsMounted] = useState(false);
  const [isInitialized, setIsInitialized] = useState(false);

  // Set mounted state after hydration to prevent SSR mismatch
  useEffect(() => {
    setIsMounted(true);
  }, []);

  // Initialize auth state and token management
  useEffect(() => {
    if (!isMounted) return;

    const initializeAuth = async () => {
      try {
        // Set up token manager event handlers
        const tokenManager = defaultTokenManager;

        // Handle token refresh events
        const originalOnTokenRefreshed = tokenManager["onTokenRefreshed"];
        tokenManager["onTokenRefreshed"] = tokens => {
          console.log("Token refreshed successfully");
          originalOnTokenRefreshed?.(tokens);
        };

        // Handle token refresh failures
        const originalOnRefreshFailed = tokenManager["onRefreshFailed"];
        tokenManager["onRefreshFailed"] = error => {
          console.warn("Token refresh failed:", error);
          // Clear user data from cache
          queryClient.setQueryData(queryKeys.auth.user, null);
          originalOnRefreshFailed?.(error);
        };

        // Handle token expiration
        const originalOnTokenExpired = tokenManager["onTokenExpired"];
        tokenManager["onTokenExpired"] = () => {
          console.warn("Token expired");
          // Clear user data and potentially redirect
          queryClient.setQueryData(queryKeys.auth.user, null);
          originalOnTokenExpired?.();
        };

        setIsInitialized(true);
      } catch (error) {
        console.error("Failed to initialize auth:", error);
        setIsInitialized(true); // Still mark as initialized to prevent blocking
      }
    };

    initializeAuth();
  }, [isMounted, queryClient]);

  // Query for current user data
  const userQuery = useQuery({
    queryKey: queryKeys.auth.user,
    queryFn: async (): Promise<User | null> => {
      if (!isMounted || !isInitialized) {
        return null;
      }

      try {
        // Check authentication status first
        const isAuth = await authService.isAuthenticated();
        if (!isAuth) {
          return null;
        }

        return await authService.getCurrentUser();
      } catch (error) {
        // Handle different error types
        const isAxiosError =
          error && typeof error === "object" && "status" in error;
        const status = isAxiosError
          ? (error as { status?: number }).status
          : undefined;

        if (status === 401) {
          // Token issues - let token manager handle this
          console.warn("Authentication failed during user fetch");
          return null;
        }

        // For other errors, log and return null
        console.error("Failed to get user data:", error);
        return null;
      }
    },
    retry: (failureCount, error) => {
      // Don't retry on auth errors
      const isAxiosError =
        error && typeof error === "object" && "response" in error;
      const status = isAxiosError
        ? (error as { response?: { status?: number } }).response?.status
        : undefined;
      return status !== 401 && status !== 403 && failureCount < 2;
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
    enabled: isMounted && isInitialized, // Only run if mounted and initialized
    refetchOnWindowFocus: false, // Disable to reduce unnecessary requests
    refetchOnReconnect: true, // Re-fetch when connection is restored
  });

  // Login mutation
  const loginMutation = useMutation({
    mutationFn: async (credentials: LoginCredentials) => {
      const response = await authService.login(credentials);
      return response.user;
    },
    onSuccess: user => {
      // Update the user query cache
      queryClient.setQueryData(queryKeys.auth.user, user);
      // Invalidate to ensure fresh data
      queryClient.invalidateQueries({ queryKey: queryKeys.auth.user });
    },
    onError: () => {
      // Clear user data on login failure
      queryClient.setQueryData(queryKeys.auth.user, null);
    },
  });

  // Logout mutation
  const logoutMutation = useMutation({
    mutationFn: async () => {
      await authService.logout();
    },
    onSuccess: () => {
      // Clear all auth-related cached data on successful logout
      queryClient.removeQueries({ queryKey: ["auth"] });
    },
    onError: error => {
      // Log error but still clear local cache
      console.error("Logout failed:", error);
      // Clear cache anyway for security
      queryClient.removeQueries({ queryKey: ["auth"] });
    },
  });

  // Registration mutation
  const registerMutation = useMutation({
    mutationFn: async (data: RegistrationData) => {
      return await authService.register(data);
    },
  });

  // Derived state from queries and mutations
  const user = userQuery.data || null;

  // Enhanced authentication check that considers async nature
  const [authStatus, setAuthStatus] = useState<
    "checking" | "authenticated" | "unauthenticated"
  >("checking");

  useEffect(() => {
    if (!isMounted || !isInitialized) {
      setAuthStatus("checking");
      return;
    }

    const checkAuth = async () => {
      try {
        const isAuth = await authService.isAuthenticated();
        setAuthStatus(isAuth ? "authenticated" : "unauthenticated");
      } catch (error) {
        console.warn("Error checking auth status:", error);
        setAuthStatus("unauthenticated");
      }
    };

    checkAuth();
  }, [isMounted, isInitialized, user]);

  const isAuthenticated = authStatus === "authenticated";

  const isLoading =
    !isMounted ||
    !isInitialized ||
    authStatus === "checking" ||
    userQuery.isLoading ||
    loginMutation.isPending ||
    logoutMutation.isPending ||
    registerMutation.isPending;

  // Combine all possible errors
  const error =
    userQuery.error ||
    loginMutation.error ||
    logoutMutation.error ||
    registerMutation.error;
  const errorMessage = error instanceof Error ? error.message : null;

  // Authentication methods
  const login = async (credentials: LoginCredentials) => {
    await loginMutation.mutateAsync(credentials);
  };

  const logout = async () => {
    await logoutMutation.mutateAsync();
  };

  const register = async (data: RegistrationData) => {
    await registerMutation.mutateAsync(data);
  };

  // Clear error function - mutation reset functions are stable
  const clearError = () => {
    loginMutation.reset();
    logoutMutation.reset();
    registerMutation.reset();
  };

  const refreshUser = async () => {
    await queryClient.invalidateQueries({ queryKey: queryKeys.auth.user });
  };

  // Enhanced context value with additional utilities
  const contextValue: AuthContextValue & {
    getTokenInfo: () => Promise<{
      accessToken: string | null;
      isExpired: boolean;
      expiresAt: Date | null;
      timeUntilExpiry: number | null;
    }>;
    refreshToken: () => Promise<boolean>;
  } = {
    user,
    isAuthenticated,
    isLoading,
    error: errorMessage,
    login,
    logout,
    register,
    clearError,
    refreshUser,
    // Additional utilities
    getTokenInfo: authService.getTokenInfo,
    refreshToken: authService.refreshToken,
  };

  return (
    <AuthContext.Provider value={contextValue}>{children}</AuthContext.Provider>
  );
}

// Custom hook for using authentication context
export function useAuth(): AuthContextValue {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
