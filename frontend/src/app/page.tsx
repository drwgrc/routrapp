"use client";

import { useEffect, useState } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/styles/components/ui/card";

interface ApiStatus {
  status: number | null;
  loading: boolean;
  error: string | null;
}

export default function Home() {
  const [apiStatus, setApiStatus] = useState<ApiStatus>({
    status: null,
    loading: true,
    error: null,
  });

  const getStatusColor = (status: number | null): string => {
    if (!status) return "text-gray-500";

    if (status >= 200 && status < 300) return "text-green-500";
    if (status >= 400 && status < 500) {
      if (status === 404) return "text-yellow-500";
      return "text-red-500";
    }
    if (status >= 500) return "text-red-500";
    return "text-blue-500";
  };

  const getStatusText = (status: number | null): string => {
    if (!status) return "Unknown";

    if (status >= 200 && status < 300) return "Healthy";
    if (status >= 400 && status < 500) return "Client Error";
    if (status >= 500) return "Server Error";
    return "Info";
  };

  useEffect(() => {
    const checkApiStatus = async () => {
      try {
        const apiBaseUrl =
          process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";

        const response = await fetch(apiBaseUrl, {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
          },
        });

        setApiStatus({
          status: response.status,
          loading: false,
          error: null,
        });
      } catch (error) {
        setApiStatus({
          status: null,
          loading: false,
          error: error instanceof Error ? error.message : "Unknown error",
        });
      }
    };

    checkApiStatus();
  }, []);

  return (
    <div className="min-h-screen flex items-center justify-center p-4 bg-gray-50 dark:bg-gray-900">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <CardTitle>API Status</CardTitle>
          <CardDescription>
            Checking connection to backend service
          </CardDescription>
        </CardHeader>
        <CardContent className="text-center space-y-4">
          {apiStatus.loading ? (
            <div className="flex items-center justify-center">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 dark:border-gray-100"></div>
            </div>
          ) : apiStatus.error ? (
            <div className="space-y-2">
              <div className="text-red-500 font-semibold text-lg">
                Connection Failed
              </div>
              <div className="text-sm text-gray-600 dark:text-gray-400">
                {apiStatus.error}
              </div>
            </div>
          ) : (
            <div className="space-y-2">
              <div
                className={`font-bold text-3xl ${getStatusColor(
                  apiStatus.status
                )}`}
              >
                {apiStatus.status}
              </div>
              <div
                className={`font-semibold ${getStatusColor(apiStatus.status)}`}
              >
                {getStatusText(apiStatus.status)}
              </div>
              <div className="text-sm text-gray-600 dark:text-gray-400">
                {process.env.NEXT_PUBLIC_API_BASE_URL ||
                  "http://localhost:8080"}
              </div>
            </div>
          )}

          <button
            onClick={() => window.location.reload()}
            className="mt-4 px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 transition-colors"
          >
            Refresh Status
          </button>
        </CardContent>
      </Card>
    </div>
  );
}
