"use client";

import React, { ReactNode, Suspense } from "react";
import { useAuth } from "@/contexts/auth-context";
import { usePermissions } from "@/hooks/use-permissions";
import { UserRole } from "@/types/auth";

// Loading component for visibility checks
interface LoadingFallbackProps {
  children?: ReactNode;
  className?: string;
}

function LoadingFallback({ children, className = "" }: LoadingFallbackProps) {
  return (
    <div className={`animate-pulse ${className}`}>
      {children || <div className="h-4 bg-gray-200 rounded w-3/4"></div>}
    </div>
  );
}

// Error component for visibility checks
interface ErrorFallbackProps {
  message?: string;
  className?: string;
}

function ErrorFallback({
  message = "Access denied",
  className = "",
}: ErrorFallbackProps) {
  return <div className={`text-red-600 text-sm ${className}`}>{message}</div>;
}

// Base interface for conditional render props
interface BaseConditionalRenderProps {
  children: ReactNode;
  fallback?: ReactNode;
  loadingFallback?: ReactNode;
  errorFallback?: ReactNode;
  showOnError?: boolean;
  className?: string;
}

// Role-based conditional render props
interface RoleConditionalRenderProps extends BaseConditionalRenderProps {
  allowedRoles: UserRole[];
  requireAllRoles?: boolean;
}

// Permission-based conditional render props
interface PermissionConditionalRenderProps extends BaseConditionalRenderProps {
  requiredPermissions: string[];
  requireAllPermissions?: boolean;
}

// Combined conditional render props
interface CombinedConditionalRenderProps extends BaseConditionalRenderProps {
  allowedRoles?: UserRole[];
  requiredPermissions?: string[];
  requireAllRoles?: boolean;
  requireAllPermissions?: boolean;
  logic?: "AND" | "OR";
}

// Custom condition render props
interface CustomConditionalRenderProps extends BaseConditionalRenderProps {
  condition: () => boolean;
}

/**
 * RoleConditionalRender Component
 *
 * Conditionally renders content based on user roles with loading and error states.
 * Provides smooth transitions and proper fallbacks.
 */
export function RoleConditionalRender({
  children,
  allowedRoles,
  fallback = null,
  loadingFallback,
  errorFallback,
  showOnError = false,
  className = "",
  requireAllRoles = false,
}: RoleConditionalRenderProps) {
  const { isLoading, error } = useAuth();
  const { hasAnyRole, hasRole } = usePermissions();

  // Show loading state while auth is loading
  if (isLoading) {
    return (
      <LoadingFallback className={className}>{loadingFallback}</LoadingFallback>
    );
  }

  // Show error state if there's an auth error and showOnError is false
  if (error && !showOnError) {
    return (
      <ErrorFallback
        message={(errorFallback as string) || "Authentication error"}
        className={className}
      />
    );
  }

  // Check if user has required roles
  const hasRequiredRoles = requireAllRoles
    ? allowedRoles.every(role => hasRole(role))
    : hasAnyRole(allowedRoles);

  if (hasRequiredRoles) {
    return <>{children}</>;
  }

  return <>{fallback}</>;
}

/**
 * PermissionConditionalRender Component
 *
 * Conditionally renders content based on user permissions with loading and error states.
 * Provides smooth transitions and proper fallbacks.
 */
export function PermissionConditionalRender({
  children,
  requiredPermissions,
  fallback = null,
  loadingFallback,
  errorFallback,
  showOnError = false,
  className = "",
  requireAllPermissions = false,
}: PermissionConditionalRenderProps) {
  const { isLoading, error } = useAuth();
  const { hasAnyPermission, hasPermission } = usePermissions();

  // Show loading state while auth is loading
  if (isLoading) {
    return (
      <LoadingFallback className={className}>{loadingFallback}</LoadingFallback>
    );
  }

  // Show error state if there's an auth error and showOnError is false
  if (error && !showOnError) {
    return (
      <ErrorFallback
        message={(errorFallback as string) || "Authentication error"}
        className={className}
      />
    );
  }

  // Check if user has required permissions
  const hasRequiredPermissions = requireAllPermissions
    ? requiredPermissions.every(permission => hasPermission(permission))
    : hasAnyPermission(requiredPermissions);

  if (hasRequiredPermissions) {
    return <>{children}</>;
  }

  return <>{fallback}</>;
}

/**
 * CombinedConditionalRender Component
 *
 * Conditionally renders content based on both roles and permissions with flexible logic.
 * Supports complex visibility rules with proper loading and error handling.
 */
export function CombinedConditionalRender({
  children,
  allowedRoles,
  requiredPermissions,
  fallback = null,
  loadingFallback,
  errorFallback,
  showOnError = false,
  className = "",
  requireAllRoles = false,
  requireAllPermissions = false,
  logic = "OR",
}: CombinedConditionalRenderProps) {
  const { isLoading, error } = useAuth();
  const { hasAnyRole, hasRole, hasAnyPermission, hasPermission } =
    usePermissions();

  // Show loading state while auth is loading
  if (isLoading) {
    return (
      <LoadingFallback className={className}>{loadingFallback}</LoadingFallback>
    );
  }

  // Show error state if there's an auth error and showOnError is false
  if (error && !showOnError) {
    return (
      <ErrorFallback
        message={(errorFallback as string) || "Authentication error"}
        className={className}
      />
    );
  }

  let hasRoleAccess = true;
  let hasPermissionAccess = true;

  // Check role access
  if (allowedRoles && allowedRoles.length > 0) {
    hasRoleAccess = requireAllRoles
      ? allowedRoles.every(role => hasRole(role))
      : hasAnyRole(allowedRoles);
  }

  // Check permission access
  if (requiredPermissions && requiredPermissions.length > 0) {
    hasPermissionAccess = requireAllPermissions
      ? requiredPermissions.every(permission => hasPermission(permission))
      : hasAnyPermission(requiredPermissions);
  }

  // Combine checks based on logic
  const hasAccess =
    logic === "AND"
      ? hasRoleAccess && hasPermissionAccess
      : hasRoleAccess || hasPermissionAccess;

  if (hasAccess) {
    return <>{children}</>;
  }

  return <>{fallback}</>;
}

/**
 * CustomConditionalRender Component
 *
 * Conditionally renders content based on a custom condition function.
 * Useful for complex business logic that doesn't fit standard patterns.
 */
export function CustomConditionalRender({
  children,
  condition,
  fallback = null,
  loadingFallback,
  errorFallback,
  showOnError = false,
  className = "",
}: CustomConditionalRenderProps) {
  const { isLoading, error } = useAuth();

  // Show loading state while auth is loading
  if (isLoading) {
    return (
      <LoadingFallback className={className}>{loadingFallback}</LoadingFallback>
    );
  }

  // Show error state if there's an auth error and showOnError is false
  if (error && !showOnError) {
    return (
      <ErrorFallback
        message={(errorFallback as string) || "Authentication error"}
        className={className}
      />
    );
  }

  if (condition()) {
    return <>{children}</>;
  }

  return <>{fallback}</>;
}

/**
 * SuspenseConditionalRender Component
 *
 * Wraps conditional rendering with React Suspense for better loading handling.
 * Useful for components that might need to load additional data.
 */
interface SuspenseConditionalRenderProps extends BaseConditionalRenderProps {
  condition: () => boolean;
  suspenseFallback?: ReactNode;
}

export function SuspenseConditionalRender({
  children,
  condition,
  fallback = null,
  loadingFallback,
  errorFallback,
  showOnError = false,
  className = "",
  suspenseFallback,
}: SuspenseConditionalRenderProps) {
  return (
    <Suspense
      fallback={suspenseFallback || <LoadingFallback className={className} />}
    >
      <CustomConditionalRender
        condition={condition}
        fallback={fallback}
        loadingFallback={loadingFallback}
        errorFallback={errorFallback}
        showOnError={showOnError}
        className={className}
      >
        {children}
      </CustomConditionalRender>
    </Suspense>
  );
}

/**
 * Convenience components for common use cases
 */

/**
 * AdminConditionalRender Component
 *
 * Convenience component that only renders content for admin/owner users.
 */
export function AdminConditionalRender({
  children,
  fallback = null,
  includeTechnicianAdmins = false,
  ...props
}: {
  children: ReactNode;
  fallback?: ReactNode;
  includeTechnicianAdmins?: boolean;
} & Omit<BaseConditionalRenderProps, "children" | "fallback">) {
  const allowedRoles: UserRole[] = ["owner"];

  if (includeTechnicianAdmins) {
    allowedRoles.push("technician");
  }

  return (
    <CombinedConditionalRender
      allowedRoles={allowedRoles}
      requiredPermissions={["organizations.*", "users.*"]}
      logic="OR"
      fallback={fallback}
      {...props}
    >
      {children}
    </CombinedConditionalRender>
  );
}

/**
 * ManagementConditionalRender Component
 *
 * Convenience component for management-level features.
 */
export function ManagementConditionalRender({
  children,
  fallback = null,
  ...props
}: {
  children: ReactNode;
  fallback?: ReactNode;
} & Omit<BaseConditionalRenderProps, "children" | "fallback">) {
  return (
    <PermissionConditionalRender
      requiredPermissions={[
        "organizations.manage",
        "users.manage",
        "technicians.manage",
        "routes.manage",
      ]}
      requireAllPermissions={false}
      fallback={fallback}
      {...props}
    >
      {children}
    </PermissionConditionalRender>
  );
}

/**
 * ReadOnlyConditionalRender Component
 *
 * Convenience component for read-only access.
 */
export function ReadOnlyConditionalRender({
  children,
  fallback = null,
  ...props
}: {
  children: ReactNode;
  fallback?: ReactNode;
} & Omit<BaseConditionalRenderProps, "children" | "fallback">) {
  return (
    <PermissionConditionalRender
      requiredPermissions={[
        "organizations.read",
        "users.read",
        "technicians.read",
        "routes.read",
      ]}
      requireAllPermissions={false}
      fallback={fallback}
      {...props}
    >
      {children}
    </PermissionConditionalRender>
  );
}
