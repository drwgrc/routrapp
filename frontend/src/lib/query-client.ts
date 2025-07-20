import { QueryClient, DefaultOptions } from "@tanstack/react-query";

// Default query options for the application
const queryConfig: DefaultOptions = {
  queries: {
    retry: (failureCount, error: unknown) => {
      // Don't retry on authentication errors
      const isAxiosError =
        error && typeof error === "object" && "response" in error;
      const status = isAxiosError
        ? (error as { response?: { status?: number } }).response?.status
        : undefined;

      if (status === 401 || status === 403) {
        return false;
      }
      // Retry up to 3 times for other errors
      return failureCount < 3;
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 10 * 60 * 1000, // 10 minutes (formerly cacheTime)
    refetchOnWindowFocus: false,
    refetchOnReconnect: true,
  },
  mutations: {
    retry: false, // Don't retry mutations by default
  },
};

// Create and export query client instance
export const queryClient = new QueryClient({
  defaultOptions: queryConfig,
});

// Query keys for consistent cache management
export const queryKeys = {
  auth: {
    user: ["auth", "user"] as const,
  },
  // Add more query keys as the application grows
} as const;
