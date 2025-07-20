"use client";

import { useEffect, useState } from "react";
import { MainLayout } from "@/components/layout";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useTheme } from "@/components/theme-provider";
import apiClient, { ApiError } from "@/lib/api/api-client";
import { AxiosResponse } from "axios";

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
  const { systemTheme, actualTheme } = useTheme();

  const getStatusColor = (status: number | null): string => {
    if (status === null) return "text-gray-500";
    if (status >= 200 && status < 300) return "text-green-500";
    if (status >= 400 && status < 500) return "text-yellow-500";
    if (status >= 500) return "text-red-500";
    return "text-gray-500";
  };

  const getStatusText = (status: number | null): string => {
    if (status === null) return "Unknown";
    if (status >= 200 && status < 300) return "Healthy";
    if (status >= 400 && status < 500) return "Warning";
    if (status >= 500) return "Error";
    return "Unknown";
  };

  useEffect(() => {
    const checkApiStatus = async () => {
      try {
        setApiStatus(prev => ({ ...prev, loading: true, error: null }));
        const response = (await apiClient.get(
          "/health",
          undefined,
          true
        )) as AxiosResponse;
        setApiStatus({
          status: response.status,
          loading: false,
          error: null,
        });
      } catch (error) {
        const apiError = error as ApiError;
        setApiStatus({
          status: apiError.status || null,
          loading: false,
          error: apiError.message || "Unknown error",
        });
      }
    };

    checkApiStatus();
  }, []);

  return (
    <MainLayout>
      <div className="space-y-6">
        {/* Welcome Section */}
        <div className="space-y-2">
          <h1 className="text-3xl font-bold tracking-tight">
            Welcome to RoutrApp
          </h1>
          <p className="text-muted-foreground">
            Your route optimization platform for utility and trade companies.
          </p>
        </div>

        {/* Status Cards */}
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {/* API Status Card */}
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">API Status</CardTitle>
              <div
                className={`h-2 w-2 rounded-full ${
                  apiStatus.loading
                    ? "bg-gray-400 animate-pulse"
                    : apiStatus.status &&
                        apiStatus.status >= 200 &&
                        apiStatus.status < 300
                      ? "bg-green-500"
                      : "bg-red-500"
                }`}
              />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {apiStatus.loading ? (
                  <span className="text-gray-500">Checking...</span>
                ) : (
                  <span className={getStatusColor(apiStatus.status)}>
                    {getStatusText(apiStatus.status)}
                  </span>
                )}
              </div>
              <p className="text-xs text-muted-foreground">
                {apiStatus.error
                  ? `Error: ${apiStatus.error}`
                  : apiStatus.status
                    ? `HTTP ${apiStatus.status}`
                    : "Backend connection status"}
              </p>
            </CardContent>
          </Card>

          {/* Theme Status Card */}
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                Theme Status
              </CardTitle>
              <div
                className={`h-2 w-2 rounded-full ${actualTheme === "dark" ? "bg-gray-700" : "bg-yellow-400"}`}
              />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold capitalize">{actualTheme}</div>
              <p className="text-xs text-muted-foreground">
                System preference: {systemTheme}
              </p>
            </CardContent>
          </Card>

          {/* System Info Card */}
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Platform</CardTitle>
              <div className="h-2 w-2 rounded-full bg-blue-500" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">Web</div>
              <p className="text-xs text-muted-foreground">
                Next.js 15 with TypeScript
              </p>
            </CardContent>
          </Card>
        </div>

        {/* Quick Actions */}
        <Card>
          <CardHeader>
            <CardTitle>Quick Actions</CardTitle>
            <CardDescription>
              Get started with RoutrApp by exploring these key features.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              <div className="flex flex-col items-center p-4 border rounded-lg hover:bg-accent transition-colors cursor-pointer">
                <div className="text-2xl mb-2">🗺️</div>
                <h3 className="font-semibold">View Routes</h3>
                <p className="text-sm text-muted-foreground text-center">
                  Manage and optimize your delivery routes
                </p>
              </div>
              <div className="flex flex-col items-center p-4 border rounded-lg hover:bg-accent transition-colors cursor-pointer">
                <div className="text-2xl mb-2">👥</div>
                <h3 className="font-semibold">Manage Technicians</h3>
                <p className="text-sm text-muted-foreground text-center">
                  Add and manage your field technicians
                </p>
              </div>
              <div className="flex flex-col items-center p-4 border rounded-lg hover:bg-accent transition-colors cursor-pointer">
                <div className="text-2xl mb-2">📊</div>
                <h3 className="font-semibold">View Analytics</h3>
                <p className="text-sm text-muted-foreground text-center">
                  Track performance and efficiency metrics
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </MainLayout>
  );
}
