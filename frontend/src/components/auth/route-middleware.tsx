"use client";

import React, { useEffect, useMemo } from "react";
import { useRouter, usePathname } from "next/navigation";
import { useAuth } from "@/contexts/auth-context";
import { usePermissions } from "@/hooks/use-permissions";
import { UserRole } from "@/types/auth";

interface RouteMiddlewareProps {
  children: React.ReactNode;
}

/**
 * AuthMiddleware Component
 *
 * Global route middleware that handles authentication state changes
 * and provides automatic redirects based on authentication status
 */
export function AuthMiddleware({ children }: RouteMiddlewareProps) {
  const { isAuthenticated, isLoading } = useAuth();
  const router = useRouter();
  const pathname = usePathname();

  // Define public routes that don't require authentication
  const publicRoutes = useMemo(() => ["/login", "/register", "/"], []);

  // Define routes that should redirect authenticated users
  const guestOnlyRoutes = useMemo(() => ["/login", "/register"], []);

  useEffect(() => {
    if (isLoading) return;

    const isPublicRoute = publicRoutes.includes(pathname);
    const isGuestOnlyRoute = guestOnlyRoutes.includes(pathname);

    // Redirect authenticated users away from guest-only routes
    if (isAuthenticated && isGuestOnlyRoute) {
      router.push("/dashboard");
      return;
    }

    // Redirect unauthenticated users from protected routes
    if (!isAuthenticated && !isPublicRoute) {
      router.push("/login");
      return;
    }
  }, [
    isAuthenticated,
    isLoading,
    pathname,
    router,
    publicRoutes,
    guestOnlyRoutes,
  ]);

  return <>{children}</>;
}

/**
 * RoleRedirectMiddleware Component
 *
 * Middleware that automatically redirects users to role-appropriate pages
 */
export function RoleRedirectMiddleware({ children }: RouteMiddlewareProps) {
  const { isAuthenticated, user } = useAuth();
  const router = useRouter();
  const pathname = usePathname();

  useEffect(() => {
    if (!isAuthenticated || !user) return;

    const userRole = user.role as UserRole;

    // Define role-specific default routes
    const roleDefaultRoutes: Record<UserRole, string> = {
      owner: "/admin/dashboard",
      technician: "/technician/dashboard",
    };

    // Redirect to role-appropriate page if on root
    if (pathname === "/" || pathname === "/dashboard") {
      const defaultRoute = roleDefaultRoutes[userRole];
      if (defaultRoute) {
        router.push(defaultRoute);
      }
    }
  }, [isAuthenticated, user, pathname, router]);

  return <>{children}</>;
}

/**
 * PermissionMiddleware Component
 *
 * Middleware that checks permissions for the current route
 * and redirects if user doesn't have required permissions
 */
interface PermissionMiddlewareProps extends RouteMiddlewareProps {
  routePermissions?: Record<string, string[]>;
}

export function PermissionMiddleware({
  children,
  routePermissions = {},
}: PermissionMiddlewareProps) {
  const { isAuthenticated } = useAuth();
  const { hasAnyPermission } = usePermissions();
  const router = useRouter();
  const pathname = usePathname();

  useEffect(() => {
    if (!isAuthenticated) return;

    const requiredPermissions = routePermissions[pathname];
    if (requiredPermissions && !hasAnyPermission(requiredPermissions)) {
      router.push("/unauthorized");
    }
  }, [isAuthenticated, pathname, routePermissions, hasAnyPermission, router]);

  return <>{children}</>;
}

/**
 * SessionTimeoutMiddleware Component
 *
 * Middleware that handles session timeouts and token refresh
 */
export function SessionTimeoutMiddleware({ children }: RouteMiddlewareProps) {
  const { isAuthenticated, refreshUser, logout } = useAuth();

  useEffect(() => {
    if (!isAuthenticated) return;

    let timeoutId: NodeJS.Timeout;

    const checkSession = async () => {
      try {
        await refreshUser();
        // Schedule next check in 5 minutes
        timeoutId = setTimeout(checkSession, 5 * 60 * 1000);
      } catch (error) {
        // If refresh fails, logout the user
        console.error("Session refresh failed:", error);
        await logout();
      }
    };

    // Initial check after 5 minutes
    timeoutId = setTimeout(checkSession, 5 * 60 * 1000);

    return () => {
      if (timeoutId) {
        clearTimeout(timeoutId);
      }
    };
  }, [isAuthenticated, refreshUser, logout]);

  return <>{children}</>;
}

/**
 * CombinedMiddleware Component
 *
 * Combines all middleware components into a single wrapper
 */
interface CombinedMiddlewareProps extends RouteMiddlewareProps {
  routePermissions?: Record<string, string[]>;
  enableSessionTimeout?: boolean;
  enableRoleRedirect?: boolean;
}

export function CombinedMiddleware({
  children,
  routePermissions,
  enableSessionTimeout = true,
  enableRoleRedirect = true,
}: CombinedMiddlewareProps) {
  let content = children;

  // Wrap with permission middleware if route permissions are provided
  if (routePermissions) {
    content = (
      <PermissionMiddleware routePermissions={routePermissions}>
        {content}
      </PermissionMiddleware>
    );
  }

  // Wrap with role redirect middleware if enabled
  if (enableRoleRedirect) {
    content = <RoleRedirectMiddleware>{content}</RoleRedirectMiddleware>;
  }

  // Wrap with session timeout middleware if enabled
  if (enableSessionTimeout) {
    content = <SessionTimeoutMiddleware>{content}</SessionTimeoutMiddleware>;
  }

  // Always wrap with auth middleware as the outermost layer
  return <AuthMiddleware>{content}</AuthMiddleware>;
}
