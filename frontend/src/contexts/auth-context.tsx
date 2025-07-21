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
import {
  AuthContextValue,
  LoginCredentials,
  RegistrationData,
  ProfileUpdateData,
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

  // Set mounted state after hydration to prevent SSR mismatch
  useEffect(() => {
    setIsMounted(true);
  }, []);

  // Query for current user data
  const userQuery = useQuery({
    queryKey: queryKeys.auth.user,
    queryFn: async (): Promise<User | null> => {
      if (!isMounted || !authService.isAuthenticated()) {
        return null;
      }
      try {
        return await authService.getCurrentUser();
      } catch (error) {
        // For 401 errors, don't clear tokens immediately
        // This could be a temporary network issue or token refresh needed
        const isAxiosError =
          error && typeof error === "object" && "status" in error;
        const status = isAxiosError
          ? (error as { status?: number }).status
          : undefined;

        if (status === 401) {
          // Don't clear tokens on 401 - let the user stay logged in
          // The token might be valid but the server might be having issues
          // or the user might be on a page that doesn't require authentication
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
      return status !== 401 && status !== 403 && failureCount < 3;
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
    enabled: isMounted && authService.isAuthenticated(), // Only run if mounted and has token
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
      // Clear only auth-related cached data on successful logout
      queryClient.removeQueries({ queryKey: ["auth"] });
    },
    onError: error => {
      // Log error but don't clear cache on failed logout
      console.error("Logout failed:", error);
    },
  });

  // Registration mutation
  const registerMutation = useMutation({
    mutationFn: async (data: RegistrationData) => {
      return await authService.register(data);
    },
  });

  // Profile update mutation
  const updateProfileMutation = useMutation({
    mutationFn: async (data: ProfileUpdateData) => {
      return await authService.updateProfile(data);
    },
    onSuccess: updatedUser => {
      // Update the user query cache with new data
      queryClient.setQueryData(queryKeys.auth.user, updatedUser);
    },
    onError: error => {
      console.error("Profile update failed:", error);
    },
  });

  // Derived state from queries and mutations
  const user = userQuery.data || null;
  // Consider user authenticated if they have a token, even if user data fetch fails
  const isAuthenticated = authService.isAuthenticated() && isMounted;
  const isLoading =
    !isMounted ||
    userQuery.isLoading ||
    loginMutation.isPending ||
    logoutMutation.isPending ||
    registerMutation.isPending ||
    updateProfileMutation.isPending;

  // Combine all possible errors
  const error =
    userQuery.error ||
    loginMutation.error ||
    logoutMutation.error ||
    registerMutation.error ||
    updateProfileMutation.error;
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

  const updateProfile = async (data: ProfileUpdateData) => {
    await updateProfileMutation.mutateAsync(data);
  };

  // Clear error function - mutation reset functions are stable
  const clearError = () => {
    loginMutation.reset();
    logoutMutation.reset();
    registerMutation.reset();
    updateProfileMutation.reset();
  };

  const refreshUser = async () => {
    await queryClient.invalidateQueries({ queryKey: queryKeys.auth.user });
  };

  const contextValue: AuthContextValue = {
    user,
    isAuthenticated,
    isLoading,
    error: errorMessage,
    login,
    logout,
    register,
    updateProfile,
    clearError,
    refreshUser,
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
