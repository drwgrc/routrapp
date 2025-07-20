import apiClient from "./api-client";
import axiosInstance from "./axios-instance";
import {
  setupRequestInterceptors,
  setupResponseInterceptors,
} from "./interceptors";
import type { ApiResponse, ApiError } from "./api-client";

export {
  apiClient,
  axiosInstance,
  setupRequestInterceptors,
  setupResponseInterceptors,
  ApiResponse,
  ApiError,
};

export default apiClient;
