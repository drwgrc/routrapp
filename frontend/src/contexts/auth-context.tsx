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
  ExtendedAuthContextValue,
  LoginCredentials,
  RegistrationData,
  User,
} from "@/types/auth";

// Create authentication context
const AuthContext = createContext<ExtendedAuthContextValue | undefined>(
  undefined
);

// Authentication provider props
interface AuthProviderProps {
  children: ReactNode;
}

// Authentication provider component
export function AuthProvider({ children }: AuthProviderProps) {
  const queryClient = useQueryClient();
  const [isMounted, setIsMounted] = useState(false);
  const [isInitialized, setIsInitialized] = useState(false);
  const [tokenBasedAuth, setTokenBasedAuth] = useState<boolean | null>(null);

  // Set mounted state after hydration to prevent SSR mismatch
  useEffect(() => {
    setIsMounted(true);
  }, []);

  // Initialize auth state and token management
  useEffect(() => {
    if (!isMounted) return;

    const initializeAuth = async () => {
      try {
        // Set up token manager event handlers using the public API
        defaultTokenManager.setEventHandlers({
          onTokenRefreshed: () => {
            console.log("Token refreshed successfully");
            // Refresh user data after token refresh
            queryClient.invalidateQueries({ queryKey: queryKeys.auth.user });
          },
          onRefreshFailed: error => {
            console.warn("Token refresh failed:", error);
            // Update token-based auth state
            setTokenBasedAuth(false);
            // Clear user data from cache
            queryClient.setQueryData(queryKeys.auth.user, null);
          },
          onTokenExpired: () => {
            console.warn("Token expired");
            // Update token-based auth state
            setTokenBasedAuth(false);
            // Clear user data and potentially redirect
            queryClient.setQueryData(queryKeys.auth.user, null);
          },
        });

        // Check initial authentication state based on tokens
        const isAuth = await authService.isAuthenticated();
        setTokenBasedAuth(isAuth);
        setIsInitialized(true);
      } catch (error) {
        console.error("Failed to initialize auth:", error);
        setTokenBasedAuth(false);
        setIsInitialized(true); // Still mark as initialized to prevent blocking
      }
    };

    initializeAuth();
  }, [isMounted, queryClient]);

  // Query for current user data - this runs independently of authentication status
  const userQuery = useQuery({
    queryKey: queryKeys.auth.user,
    queryFn: async (): Promise<User | null> => {
      if (!isMounted || !isInitialized || tokenBasedAuth === false) {
        return null;
      }

      try {
        return await authService.getCurrentUser();
      } catch (error) {
        // Handle different error types
        const isAxiosError =
          error && typeof error === "object" && "status" in error;
        const status = isAxiosError
          ? (error as { status?: number }).status
          : undefined;

        if (status === 401) {
          // Token issues - update token-based auth state
          console.warn("Authentication failed during user fetch");
          setTokenBasedAuth(false);
          return null;
        }

        // For other errors (network, server issues), log and return null
        // but don't mark user as unauthenticated if we have valid tokens
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
    enabled: isMounted && isInitialized && tokenBasedAuth === true, // Only run if we have valid tokens
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
      // Update token-based auth state
      setTokenBasedAuth(true);
      // Update the user query cache
      queryClient.setQueryData(queryKeys.auth.user, user);
      // Invalidate to ensure fresh data
      queryClient.invalidateQueries({ queryKey: queryKeys.auth.user });
    },
    onError: () => {
      // Update auth states on login failure
      setTokenBasedAuth(false);
      queryClient.setQueryData(queryKeys.auth.user, null);
    },
  });

  // Logout mutation
  const logoutMutation = useMutation({
    mutationFn: async () => {
      await authService.logout();
    },
    onSuccess: () => {
      // Update token-based auth state
      setTokenBasedAuth(false);
      // Clear all auth-related cached data on successful logout
      queryClient.removeQueries({ queryKey: ["auth"] });
    },
    onError: error => {
      // Log error but still clear local cache
      console.error("Logout failed:", error);
      // Update auth states anyway for security
      setTokenBasedAuth(false);
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

  // Robust authentication state that prioritizes token validity over user data fetching success
  // This prevents false unauthenticated states when user data fetch fails but tokens are valid
  const isAuthenticated = tokenBasedAuth === true;

  // Loading state includes token-based auth check and user data fetching when appropriate
  const isLoading =
    !isMounted ||
    !isInitialized ||
    tokenBasedAuth === null ||
    (tokenBasedAuth === true && userQuery.isLoading) ||
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
    // Re-check token-based authentication first
    const isAuth = await authService.isAuthenticated();
    setTokenBasedAuth(isAuth);

    // If still authenticated, refresh user data
    if (isAuth) {
      await queryClient.invalidateQueries({ queryKey: queryKeys.auth.user });
    }
  };

  // Enhanced context value with additional utilities
  const contextValue: ExtendedAuthContextValue = {
    user,
    isAuthenticated,
    isLoading,
    error: errorMessage,
    login,
    logout,
    register,
    clearError,
    refreshUser,
    // Additional utilities - properly typed
    getTokenInfo: authService.getTokenInfo,
    refreshToken: authService.refreshToken,
  };

  return (
    <AuthContext.Provider value={contextValue}>{children}</AuthContext.Provider>
  );
}

// Custom hook for using authentication context
export function useAuth(): ExtendedAuthContextValue {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
