import axios from "axios";

// Environment variables with fallback values
const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080/api";
const API_VERSION = process.env.NEXT_PUBLIC_API_VERSION || "v1";

// Create axios instance with default configuration
const axiosInstance = axios.create({
  baseURL: `${API_BASE_URL}/${API_VERSION}`,
  headers: {
    "Content-Type": "application/json",
    Accept: "application/json",
  },
  timeout: 10000, // 10 seconds
});

export default axiosInstance;
