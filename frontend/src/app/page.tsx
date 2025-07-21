"use client";

import { useEffect, useState } from "react";
import { MainLayout } from "@/components/layout";
import { useAuth } from "@/contexts/auth-context";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useTheme } from "@/components/theme-provider";
import apiClient, { ApiError } from "@/lib/api/api-client";
import { AxiosResponse } from "axios";
import Link from "next/link";
import {
  ArrowRight,
  MapPin,
  Users,
  BarChart3,
  Shield,
  Zap,
  Clock,
} from "lucide-react";

interface ApiStatus {
  status: number | null;
  loading: boolean;
  error: string | null;
}

export default function Home() {
  const { isAuthenticated, user } = useAuth();
  const [apiStatus, setApiStatus] = useState<ApiStatus>({
    status: null,
    loading: true,
    error: null,
  });
  const { actualTheme } = useTheme();

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

  // If user is authenticated, show a welcome back message with quick actions
  // The RoleRedirectMiddleware will handle automatic redirects to role-specific pages
  if (isAuthenticated && user) {
    return (
      <MainLayout>
        <div className="container mx-auto px-4 py-8">
          <div className="text-center space-y-6">
            <div className="space-y-2">
              <h1 className="text-4xl font-bold tracking-tight">
                Welcome back, {user.first_name}!
              </h1>
              <p className="text-lg text-muted-foreground">
                Your route optimization platform is ready to go.
              </p>
            </div>

            {/* Quick Access Cards */}
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3 max-w-4xl mx-auto">
              <Card className="hover:shadow-lg transition-shadow cursor-pointer">
                <CardContent className="p-6 text-center">
                  <MapPin className="w-12 h-12 mx-auto mb-4 text-primary" />
                  <h3 className="text-lg font-semibold mb-2">Routes</h3>
                  <p className="text-sm text-muted-foreground mb-4">
                    Manage and optimize your delivery routes
                  </p>
                  <Button asChild size="sm">
                    <Link href="/routes">
                      View Routes <ArrowRight className="w-4 h-4 ml-2" />
                    </Link>
                  </Button>
                </CardContent>
              </Card>

              <Card className="hover:shadow-lg transition-shadow cursor-pointer">
                <CardContent className="p-6 text-center">
                  <Users className="w-12 h-12 mx-auto mb-4 text-primary" />
                  <h3 className="text-lg font-semibold mb-2">Technicians</h3>
                  <p className="text-sm text-muted-foreground mb-4">
                    Manage your field technician team
                  </p>
                  <Button asChild size="sm">
                    <Link href="/technicians">
                      Manage Team <ArrowRight className="w-4 h-4 ml-2" />
                    </Link>
                  </Button>
                </CardContent>
              </Card>

              <Card className="hover:shadow-lg transition-shadow cursor-pointer">
                <CardContent className="p-6 text-center">
                  <BarChart3 className="w-12 h-12 mx-auto mb-4 text-primary" />
                  <h3 className="text-lg font-semibold mb-2">Analytics</h3>
                  <p className="text-sm text-muted-foreground mb-4">
                    Track performance and efficiency
                  </p>
                  <Button asChild size="sm">
                    <Link href="/analytics">
                      View Reports <ArrowRight className="w-4 h-4 ml-2" />
                    </Link>
                  </Button>
                </CardContent>
              </Card>
            </div>

            {/* System Status */}
            <Card className="max-w-md mx-auto">
              <CardHeader>
                <CardTitle className="text-sm">System Status</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="flex items-center justify-between">
                  <span className="text-sm">API Health</span>
                  <div className="flex items-center gap-2">
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
                    <span
                      className={`text-sm ${getStatusColor(apiStatus.status)}`}
                    >
                      {apiStatus.loading
                        ? "Checking..."
                        : getStatusText(apiStatus.status)}
                    </span>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </MainLayout>
    );
  }

  // Public landing page for unauthenticated users
  return (
    <MainLayout>
      <div className="container mx-auto px-4 py-8">
        {/* Hero Section */}
        <section className="text-center space-y-8 py-12">
          <div className="space-y-4">
            <h1 className="text-5xl md:text-6xl font-bold tracking-tight">
              Welcome to{" "}
              <span className="bg-gradient-to-r from-primary to-blue-600 bg-clip-text text-transparent">
                RoutrApp
              </span>
            </h1>
            <p className="text-xl text-muted-foreground max-w-3xl mx-auto">
              The complete route optimization platform for utility and trade
              companies. Streamline your operations, manage technicians, and
              optimize routes with ease.
            </p>
          </div>

          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Button asChild size="lg" className="text-lg px-8">
              <Link href="/register">
                Get Started <ArrowRight className="w-5 h-5 ml-2" />
              </Link>
            </Button>
            <Button
              asChild
              variant="outline"
              size="lg"
              className="text-lg px-8"
            >
              <Link href="/login">Sign In</Link>
            </Button>
          </div>
        </section>

        {/* Features Section */}
        <section className="py-16">
          <div className="text-center mb-12">
            <h2 className="text-3xl font-bold mb-4">
              Everything you need to optimize routes
            </h2>
            <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
              Powerful features designed specifically for utility and trade
              companies to manage technicians and optimize delivery routes.
            </p>
          </div>

          <div className="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
            <Card className="text-center p-6">
              <CardContent className="space-y-4">
                <div className="w-16 h-16 mx-auto bg-primary/10 rounded-full flex items-center justify-center">
                  <MapPin className="w-8 h-8 text-primary" />
                </div>
                <h3 className="text-xl font-semibold">
                  Smart Route Optimization
                </h3>
                <p className="text-muted-foreground">
                  Automatically optimize routes using advanced algorithms to
                  reduce travel time and fuel costs.
                </p>
              </CardContent>
            </Card>

            <Card className="text-center p-6">
              <CardContent className="space-y-4">
                <div className="w-16 h-16 mx-auto bg-primary/10 rounded-full flex items-center justify-center">
                  <Users className="w-8 h-8 text-primary" />
                </div>
                <h3 className="text-xl font-semibold">Technician Management</h3>
                <p className="text-muted-foreground">
                  Manage your field team with real-time tracking, skill
                  assignments, and availability scheduling.
                </p>
              </CardContent>
            </Card>

            <Card className="text-center p-6">
              <CardContent className="space-y-4">
                <div className="w-16 h-16 mx-auto bg-primary/10 rounded-full flex items-center justify-center">
                  <BarChart3 className="w-8 h-8 text-primary" />
                </div>
                <h3 className="text-xl font-semibold">Analytics & Insights</h3>
                <p className="text-muted-foreground">
                  Get detailed insights into performance metrics, efficiency
                  gains, and cost savings.
                </p>
              </CardContent>
            </Card>

            <Card className="text-center p-6">
              <CardContent className="space-y-4">
                <div className="w-16 h-16 mx-auto bg-primary/10 rounded-full flex items-center justify-center">
                  <Shield className="w-8 h-8 text-primary" />
                </div>
                <h3 className="text-xl font-semibold">Enterprise Security</h3>
                <p className="text-muted-foreground">
                  Multi-tenant architecture with role-based access control and
                  enterprise-grade security.
                </p>
              </CardContent>
            </Card>

            <Card className="text-center p-6">
              <CardContent className="space-y-4">
                <div className="w-16 h-16 mx-auto bg-primary/10 rounded-full flex items-center justify-center">
                  <Zap className="w-8 h-8 text-primary" />
                </div>
                <h3 className="text-xl font-semibold">Real-time Updates</h3>
                <p className="text-muted-foreground">
                  Get instant notifications and updates on route changes, job
                  completions, and technician status.
                </p>
              </CardContent>
            </Card>

            <Card className="text-center p-6">
              <CardContent className="space-y-4">
                <div className="w-16 h-16 mx-auto bg-primary/10 rounded-full flex items-center justify-center">
                  <Clock className="w-8 h-8 text-primary" />
                </div>
                <h3 className="text-xl font-semibold">Time Tracking</h3>
                <p className="text-muted-foreground">
                  Track time spent on jobs, monitor productivity, and generate
                  accurate billing reports.
                </p>
              </CardContent>
            </Card>
          </div>
        </section>

        {/* CTA Section */}
        <section className="py-16 text-center">
          <Card className="max-w-2xl mx-auto">
            <CardHeader>
              <CardTitle className="text-2xl">
                Ready to optimize your routes?
              </CardTitle>
              <CardDescription className="text-lg">
                Join thousands of companies that trust RoutrApp for their route
                optimization needs.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="flex flex-col sm:flex-row gap-4 justify-center">
                <Button asChild size="lg">
                  <Link href="/register">
                    Start Free Trial <ArrowRight className="w-4 h-4 ml-2" />
                  </Link>
                </Button>
                <Button asChild variant="outline" size="lg">
                  <Link href="/login">I already have an account</Link>
                </Button>
              </div>

              {/* System Status for transparency */}
              <div className="pt-4 border-t">
                <div className="flex items-center justify-center gap-4 text-sm text-muted-foreground">
                  <div className="flex items-center gap-2">
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
                    <span>
                      System Status:{" "}
                      {apiStatus.loading
                        ? "Checking..."
                        : getStatusText(apiStatus.status)}
                    </span>
                  </div>
                  <span>•</span>
                  <span>Theme: {actualTheme}</span>
                </div>
              </div>
            </CardContent>
          </Card>
        </section>
      </div>
    </MainLayout>
  );
}
