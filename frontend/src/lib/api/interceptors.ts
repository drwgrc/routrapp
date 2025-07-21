import {
  AxiosInstance,
  AxiosError,
  InternalAxiosRequestConfig,
  AxiosResponse,
} from "axios";
import { defaultTokenManager } from "../token-manager";

// Track ongoing refresh to avoid multiple refresh attempts
let isRefreshing = false;
let failedQueue: Array<{
  resolve: (token: string) => void;
  reject: (error: AxiosError) => void;
}> = [];

/**
 * Process queued requests after token refresh
 */
const processQueue = (
  error: AxiosError | null,
  token: string | null = null
) => {
  failedQueue.forEach(({ resolve, reject }) => {
    if (error) {
      reject(error);
    } else if (token) {
      resolve(token);
    }
  });

  failedQueue = [];
};

/**
 * Configure request interceptors for axios instance
 * @param axiosInstance - The axios instance to configure
 */
export const setupRequestInterceptors = (
  axiosInstance: AxiosInstance
): void => {
  axiosInstance.interceptors.request.use(
    async (
      config: InternalAxiosRequestConfig
    ): Promise<InternalAxiosRequestConfig> => {
      try {
        // Get token from token manager (will handle refresh if needed)
        const token = await defaultTokenManager.getAccessToken();

        // If token exists, add to Authorization header
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }

        return config;
      } catch (error) {
        // If token retrieval fails, proceed without auth header
        console.warn("Failed to get access token for request:", error);
        return config;
      }
    },
    (error: AxiosError) => {
      return Promise.reject(error);
    }
  );
};

/**
 * Configure response interceptors for axios instance
 * @param axiosInstance - The axios instance to configure
 */
export const setupResponseInterceptors = (
  axiosInstance: AxiosInstance
): void => {
  axiosInstance.interceptors.response.use(
    (response: AxiosResponse) => {
      return response;
    },
    async (error: AxiosError) => {
      const originalRequest = error.config as InternalAxiosRequestConfig & {
        _retry?: boolean;
      };

      // Handle 401 Unauthorized errors - token expired or invalid
      if (
        error.response?.status === 401 &&
        originalRequest &&
        !originalRequest._retry
      ) {
        // Avoid infinite retry loops
        originalRequest._retry = true;

        // Special handling for auth endpoints
        const isAuthEndpoint = originalRequest.url?.includes("/auth/");
        const isRefreshEndpoint =
          originalRequest.url?.includes("/auth/refresh");
        const isLogoutEndpoint = originalRequest.url?.includes("/auth/logout");

        // Don't try to refresh for certain auth endpoints
        if (isRefreshEndpoint) {
          // Refresh endpoint failed - clear tokens and redirect
          console.warn("Refresh token is invalid or expired");
          await defaultTokenManager.clearTokens();

          if (typeof window !== "undefined") {
            window.location.href = "/login";
          }
          return Promise.reject(error);
        }

        // For logout endpoint, just proceed with the error
        if (isLogoutEndpoint) {
          return Promise.reject(error);
        }

        // For the /auth/me endpoint, don't automatically redirect
        // Let the auth context handle this gracefully
        if (originalRequest.url?.includes("/auth/me")) {
          return Promise.reject(error);
        }

        // For other 401 errors, attempt token refresh
        if (isRefreshing) {
          // If already refreshing, queue this request
          return new Promise((resolve, reject) => {
            failedQueue.push({
              resolve: (token: string) => {
                originalRequest.headers.Authorization = `Bearer ${token}`;
                resolve(axiosInstance(originalRequest));
              },
              reject: (err: AxiosError) => {
                reject(err);
              },
            });
          });
        }

        isRefreshing = true;

        try {
          // Attempt to refresh the token
          const newToken = await defaultTokenManager.refreshTokenIfNeeded();

          if (newToken) {
            // Update the original request with new token
            originalRequest.headers.Authorization = `Bearer ${newToken}`;

            // Process queued requests
            processQueue(null, newToken);

            // Retry the original request
            return axiosInstance(originalRequest);
          } else {
            throw new Error("No token received from refresh");
          }
        } catch (refreshError) {
          console.warn("Token refresh failed:", refreshError);

          // Process failed queue
          processQueue(error);

          // Clear tokens
          await defaultTokenManager.clearTokens();

          // For non-auth endpoints, redirect to login
          if (!isAuthEndpoint && typeof window !== "undefined") {
            window.location.href = "/login";
          }

          return Promise.reject(error);
        } finally {
          isRefreshing = false;
        }
      }

      // Handle other error types
      if (error.response?.status === 403) {
        // Forbidden - user doesn't have permission
        console.warn("Access forbidden:", error.response.data);
      }

      // Format error response for consistent handling
      const formattedError = {
        message:
          error.response?.data &&
          typeof error.response.data === "object" &&
          "error" in error.response.data &&
          typeof error.response.data.error === "object" &&
          error.response.data.error !== null &&
          "message" in error.response.data.error
            ? (error.response.data.error as { message: string }).message
            : error.response?.data &&
                typeof error.response.data === "object" &&
                "message" in error.response.data
              ? (error.response.data as { message: string }).message
              : "Something went wrong",
        status: error.response?.status,
        data: error.response?.data,
        originalError: error,
      };

      return Promise.reject(formattedError);
    }
  );
};
