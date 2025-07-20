"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/contexts/auth-context";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription } from "@/components/ui/alert";
import Link from "next/link";

export default function LogoutPage() {
  const router = useRouter();
  const { logout, isLoading, error, isAuthenticated, user } = useAuth();
  const [logoutComplete, setLogoutComplete] = useState(false);

  // Redirect if not authenticated
  useEffect(() => {
    if (!isAuthenticated && !logoutComplete) {
      router.push("/login");
    }
  }, [isAuthenticated, router, logoutComplete]);

  const handleLogout = async () => {
    try {
      await logout();
      setLogoutComplete(true);
      // Small delay before redirect to show success message
      setTimeout(() => {
        router.push("/login");
      }, 2000);
    } catch (err) {
      console.error("Logout failed:", err);
    }
  };

  const handleCancel = () => {
    router.back();
  };

  // Show success message after logout
  if (logoutComplete) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background px-4">
        <div className="w-full max-w-md space-y-6">
          <div className="text-center space-y-2">
            <h1 className="text-3xl font-bold tracking-tight text-green-600">
              Signed Out Successfully
            </h1>
            <p className="text-muted-foreground">
              You have been logged out of your RoutrApp account
            </p>
          </div>

          <Card>
            <CardContent className="pt-6">
              <div className="text-center space-y-4">
                <div className="mx-auto w-12 h-12 bg-green-100 rounded-full flex items-center justify-center">
                  <svg
                    className="w-6 h-6 text-green-600"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M5 13l4 4L19 7"
                    />
                  </svg>
                </div>
                <div className="space-y-2">
                  <h3 className="text-lg font-semibold">
                    Thank you for using RoutrApp
                  </h3>
                  <p className="text-sm text-muted-foreground">
                    Your session has been terminated securely. Redirecting to
                    login...
                  </p>
                </div>
                <div className="flex items-center justify-center space-x-2">
                  <div className="h-4 w-4 rounded-full bg-primary animate-pulse" />
                  <span className="text-sm text-muted-foreground">
                    Redirecting...
                  </span>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    );
  }

  // If not authenticated and logout not complete, show loading
  if (!isAuthenticated) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <Card className="w-full max-w-md">
          <CardContent className="pt-6">
            <div className="flex items-center justify-center space-x-2">
              <div className="h-4 w-4 rounded-full bg-primary animate-pulse" />
              <span className="text-sm text-muted-foreground">
                Redirecting...
              </span>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-background px-4">
      <div className="w-full max-w-md space-y-6">
        {/* Header Section */}
        <div className="text-center space-y-2">
          <h1 className="text-3xl font-bold tracking-tight">Sign Out</h1>
          <p className="text-muted-foreground">
            Are you sure you want to sign out of your RoutrApp account?
          </p>
        </div>

        {/* Logout Confirmation Card */}
        <Card>
          <CardHeader className="space-y-1">
            <CardTitle className="text-2xl text-center">
              Confirm Sign Out
            </CardTitle>
            <CardDescription className="text-center">
              {user
                ? `Signed in as ${user.first_name} ${user.last_name} (${user.email})`
                : "You are currently signed in"}
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            {/* Error Alert */}
            {error && (
              <Alert variant="destructive">
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            )}

            {/* User Info */}
            {user && (
              <div className="bg-muted rounded-lg p-4 space-y-2">
                <div className="flex items-center space-x-3">
                  <div className="w-10 h-10 bg-primary rounded-full flex items-center justify-center text-primary-foreground font-semibold">
                    {(user.first_name || user.last_name || user.email || "U")
                      .charAt(0)
                      .toUpperCase()}
                  </div>
                  <div>
                    <p className="font-medium">
                      {user.first_name} {user.last_name}
                    </p>
                    <p className="text-sm text-muted-foreground">
                      {user.email}
                    </p>
                    <p className="text-xs text-muted-foreground capitalize">
                      {user.role} Role
                    </p>
                  </div>
                </div>
              </div>
            )}

            {/* Warning Message */}
            <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
              <div className="flex items-start space-x-2">
                <svg
                  className="w-5 h-5 text-yellow-600 mt-0.5"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"
                  />
                </svg>
                <div>
                  <p className="text-sm font-medium text-yellow-800">
                    Before signing out
                  </p>
                  <p className="text-sm text-yellow-700">
                    Make sure you have saved any unsaved work. You will need to
                    sign in again to access your account.
                  </p>
                </div>
              </div>
            </div>

            {/* Action Buttons */}
            <div className="flex space-x-3">
              <Button
                variant="outline"
                className="flex-1"
                onClick={handleCancel}
                disabled={isLoading}
              >
                Cancel
              </Button>
              <Button
                variant="destructive"
                className="flex-1"
                onClick={handleLogout}
                disabled={isLoading}
              >
                {isLoading ? (
                  <div className="flex items-center space-x-2">
                    <div className="h-4 w-4 rounded-full border-2 border-white border-t-transparent animate-spin" />
                    <span>Signing out...</span>
                  </div>
                ) : (
                  "Sign Out"
                )}
              </Button>
            </div>

            {/* Quick Actions */}
            <div className="pt-4 border-t">
              <p className="text-sm text-muted-foreground text-center mb-3">
                Or continue working:
              </p>
              <div className="grid grid-cols-2 gap-2">
                <Button variant="ghost" size="sm" asChild>
                  <Link href="/">
                    <span className="mr-2">🏠</span>
                    Dashboard
                  </Link>
                </Button>
                <Button variant="ghost" size="sm" asChild>
                  <Link href="/routes">
                    <span className="mr-2">🗺️</span>
                    Routes
                  </Link>
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Security Note */}
        <Card>
          <CardContent className="pt-4">
            <div className="text-center text-xs text-muted-foreground">
              <p>
                🔒 Your session will be terminated securely across all devices
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
