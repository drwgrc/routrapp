import {
  AxiosInstance,
  AxiosError,
  InternalAxiosRequestConfig,
  AxiosResponse,
} from "axios";

/**
 * Configure request interceptors for axios instance
 * @param axiosInstance - The axios instance to configure
 */
export const setupRequestInterceptors = (
  axiosInstance: AxiosInstance
): void => {
  axiosInstance.interceptors.request.use(
    (config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
      // Get token from local storage or other secure storage
      const token = localStorage.getItem("auth_token");

      // If token exists, add to Authorization header
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }

      return config;
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
      const originalRequest = error.config;

      // Handle 401 Unauthorized errors - token expired
      if (error.response?.status === 401 && originalRequest) {
        try {
          // Attempt to refresh the token - implement your token refresh logic here
          // const refreshToken = localStorage.getItem('refresh_token');
          // Call your refresh token endpoint
          // Update the tokens in storage

          // Retry the original request with new token
          // const token = localStorage.getItem('auth_token');
          // originalRequest.headers.Authorization = `Bearer ${token}`;
          // return axiosInstance(originalRequest);

          // For now, just redirect to login
          window.location.href = "/login";
          return Promise.reject(error);
        } catch (refreshError) {
          // If refresh token fails, redirect to login
          localStorage.removeItem("auth_token");
          localStorage.removeItem("refresh_token");
          window.location.href = "/login";
          return Promise.reject(refreshError);
        }
      }

      // Format error response for consistent handling
      return Promise.reject({
        message:
          error.response?.data &&
          typeof error.response.data === "object" &&
          "message" in error.response.data
            ? (error.response.data as { message: string }).message
            : "Something went wrong",
        status: error.response?.status,
        data: error.response?.data,
        originalError: error,
      });
    }
  );
};
