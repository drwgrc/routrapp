import axiosInstance from "./axios-instance";
import {
  setupRequestInterceptors,
  setupResponseInterceptors,
} from "./interceptors";

// Configure interceptors
setupRequestInterceptors(axiosInstance);
setupResponseInterceptors(axiosInstance);

// Type definitions for API responses
export interface ApiResponse<T = unknown> {
  data: T;
  message?: string;
  status: number;
}

export interface ApiError {
  message: string;
  status?: number;
  data?: unknown;
  originalError?: unknown;
}

// Generic API client with common methods
const apiClient = {
  // GET request
  get: async <T>(url: string, params?: Record<string, unknown>): Promise<T> => {
    const response = await axiosInstance.get<ApiResponse<T>>(url, { params });
    return response.data.data;
  },

  // POST request
  post: async <T, D = Record<string, unknown>>(
    url: string,
    data: D
  ): Promise<T> => {
    const response = await axiosInstance.post<ApiResponse<T>>(url, data);
    return response.data.data;
  },

  // PUT request
  put: async <T, D = Record<string, unknown>>(
    url: string,
    data: D
  ): Promise<T> => {
    const response = await axiosInstance.put<ApiResponse<T>>(url, data);
    return response.data.data;
  },

  // PATCH request
  patch: async <T, D = Record<string, unknown>>(
    url: string,
    data: D
  ): Promise<T> => {
    const response = await axiosInstance.patch<ApiResponse<T>>(url, data);
    return response.data.data;
  },

  // DELETE request
  delete: async <T>(url: string): Promise<T> => {
    const response = await axiosInstance.delete<ApiResponse<T>>(url);
    return response.data.data;
  },
};

export default apiClient;
