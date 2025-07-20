"use client";

import React, { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/contexts/auth-context";
import { usePermissions } from "@/hooks/use-permissions";
import { ProtectedRouteProps, UserRole } from "@/types/auth";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { AlertTriangle, Shield, Lock, RefreshCw } from "lucide-react";

/**
 * ProtectedRoute Component
 *
 * A comprehensive route protection component that supports:
 * - Authentication checking
 * - Role-based access control
 * - Permission-based access control
 * - Custom fallback components
 * - Automatic redirects
 * - Loading states
 * - Error handling
 */
export function ProtectedRoute({
  children,
  requireAuth = true,
  allowedRoles,
  requiredPermissions,
  redirectTo = "/login",
  fallback: FallbackComponent,
}: ProtectedRouteProps) {
  const { isAuthenticated, isLoading, user, refreshUser } = useAuth();
  const permissions = usePermissions();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading && requireAuth && !isAuthenticated) {
      router.push(redirectTo);
    }
  }, [isAuthenticated, isLoading, router, redirectTo, requireAuth]);

  // Show loading state while checking authentication
  if (isLoading) {
    return <LoadingState />;
  }

  // If authentication is not required and user is not authenticated, show children
  if (!requireAuth && !isAuthenticated) {
    return <>{children}</>;
  }

  // If authentication is required but user is not authenticated, show loading state
  // (user will be redirected by the useEffect above)
  if (requireAuth && !isAuthenticated) {
    return <LoadingState message="Redirecting to login..." />;
  }

  // Check role-based access
  if (allowedRoles && allowedRoles.length > 0) {
    if (!permissions.hasAnyRole(allowedRoles)) {
      if (FallbackComponent) {
        return <FallbackComponent />;
      }
      return (
        <AccessDeniedState
          type="role"
          allowedRoles={allowedRoles}
          userRole={user?.role as UserRole}
          onRefresh={refreshUser}
        />
      );
    }
  }

  // Check permission-based access
  if (requiredPermissions && requiredPermissions.length > 0) {
    if (!permissions.hasAnyPermission(requiredPermissions)) {
      if (FallbackComponent) {
        return <FallbackComponent />;
      }
      return (
        <AccessDeniedState
          type="permission"
          requiredPermissions={requiredPermissions}
          onRefresh={refreshUser}
        />
      );
    }
  }

  // All checks passed, render children
  return <>{children}</>;
}

/**
 * Loading state component for authentication checks
 */
function LoadingState({
  message = "Checking authentication...",
}: {
  message?: string;
}) {
  return (
    <div className="min-h-screen flex items-center justify-center bg-background">
      <Card className="w-full max-w-md">
        <CardContent className="pt-6">
          <div className="flex items-center justify-center space-x-3">
            <div className="h-5 w-5 rounded-full bg-primary animate-pulse" />
            <span className="text-sm text-muted-foreground">{message}</span>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

/**
 * Access denied state component with different messages for different scenarios
 */
interface AccessDeniedStateProps {
  type: "role" | "permission";
  allowedRoles?: UserRole[];
  userRole?: UserRole;
  requiredPermissions?: string[];
  onRefresh: () => Promise<void>;
}

function AccessDeniedState({
  type,
  allowedRoles,
  userRole,
  requiredPermissions,
  onRefresh,
}: AccessDeniedStateProps) {
  const router = useRouter();
  const [isRefreshing, setIsRefreshing] = React.useState(false);

  const handleRefresh = async () => {
    setIsRefreshing(true);
    try {
      await onRefresh();
    } finally {
      setIsRefreshing(false);
    }
  };

  const getTitle = () => {
    switch (type) {
      case "role":
        return "Insufficient Role Permissions";
      case "permission":
        return "Access Denied";
      default:
        return "Access Denied";
    }
  };

  const getMessage = () => {
    switch (type) {
      case "role":
        return (
          <div className="space-y-2">
            <p className="text-sm text-muted-foreground">
              Your current role ({userRole || "unknown"}) does not have access
              to this page.
            </p>
            {allowedRoles && allowedRoles.length > 0 && (
              <p className="text-sm text-muted-foreground">
                Required roles: {allowedRoles.join(", ")}
              </p>
            )}
          </div>
        );
      case "permission":
        return (
          <div className="space-y-2">
            <p className="text-sm text-muted-foreground">
              You don&apos;t have the required permissions to access this page.
            </p>
            {requiredPermissions && requiredPermissions.length > 0 && (
              <p className="text-xs text-muted-foreground">
                Required permissions: {requiredPermissions.join(", ")}
              </p>
            )}
          </div>
        );
      default:
        return (
          <p className="text-sm text-muted-foreground">
            You don&apos;t have permission to access this page.
          </p>
        );
    }
  };

  const getIcon = () => {
    switch (type) {
      case "role":
        return <Shield className="h-12 w-12 text-destructive" />;
      case "permission":
        return <Lock className="h-12 w-12 text-destructive" />;
      default:
        return <AlertTriangle className="h-12 w-12 text-destructive" />;
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-background p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <div className="flex justify-center mb-4">{getIcon()}</div>
          <CardTitle className="text-xl text-destructive">
            {getTitle()}
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {getMessage()}
          <div className="flex flex-col sm:flex-row gap-2">
            <Button
              variant="outline"
              onClick={() => router.push("/")}
              className="flex-1"
            >
              Go Home
            </Button>
            <Button
              variant="outline"
              onClick={handleRefresh}
              disabled={isRefreshing}
              className="flex-1"
            >
              {isRefreshing ? (
                <>
                  <RefreshCw className="h-4 w-4 mr-2 animate-spin" />
                  Refreshing...
                </>
              ) : (
                <>
                  <RefreshCw className="h-4 w-4 mr-2" />
                  Refresh
                </>
              )}
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
