"use client";

import React, { ReactNode } from "react";
import { usePermissions } from "@/hooks/use-permissions";
import { UserRole } from "@/types/auth";

// Base interface for all visibility guard props
interface BaseVisibilityProps {
  children: ReactNode;
  fallback?: ReactNode;
  loadingFallback?: ReactNode;
  inverse?: boolean;
  showOnError?: boolean;
}

// Role-based visibility props
interface RoleVisibilityProps extends BaseVisibilityProps {
  allowedRoles: UserRole[];
  requireAllRoles?: boolean; // If true, user must have ALL roles, not just any
}

// Permission-based visibility props
interface PermissionVisibilityProps extends BaseVisibilityProps {
  requiredPermissions: string[];
  requireAllPermissions?: boolean; // If true, user must have ALL permissions, not just any
}

// Combined visibility props
interface CombinedVisibilityProps extends BaseVisibilityProps {
  allowedRoles?: UserRole[];
  requiredPermissions?: string[];
  requireAllRoles?: boolean;
  requireAllPermissions?: boolean;
  logic?: "AND" | "OR"; // How to combine role and permission checks
}

/**
 * RoleVisibility Component
 *
 * Advanced role-based component visibility with support for:
 * - Multiple roles with AND/OR logic
 * - Loading states
 * - Error handling
 * - Inverse logic
 * - Custom fallbacks
 */
export function RoleVisibility({
  children,
  allowedRoles,
  fallback = null,
  loadingFallback = null,
  inverse = false,
  showOnError = false,
  requireAllRoles = false,
}: RoleVisibilityProps) {
  const { hasAnyRole, hasRole, isOwner, isTechnician } = usePermissions();

  // Determine if user has required roles
  const hasRequiredRoles = requireAllRoles
    ? allowedRoles.every(role => hasRole(role))
    : hasAnyRole(allowedRoles);

  const shouldRender = inverse ? !hasRequiredRoles : hasRequiredRoles;

  if (shouldRender) {
    return <>{children}</>;
  }

  return <>{fallback}</>;
}

/**
 * PermissionVisibility Component
 *
 * Advanced permission-based component visibility with support for:
 * - Multiple permissions with AND/OR logic
 * - Loading states
 * - Error handling
 * - Inverse logic
 * - Custom fallbacks
 */
export function PermissionVisibility({
  children,
  requiredPermissions,
  fallback = null,
  loadingFallback = null,
  inverse = false,
  showOnError = false,
  requireAllPermissions = false,
}: PermissionVisibilityProps) {
  const { hasAnyPermission, hasPermission } = usePermissions();

  // Determine if user has required permissions
  const hasRequiredPermissions = requireAllPermissions
    ? requiredPermissions.every(permission => hasPermission(permission))
    : hasAnyPermission(requiredPermissions);

  const shouldRender = inverse
    ? !hasRequiredPermissions
    : hasRequiredPermissions;

  if (shouldRender) {
    return <>{children}</>;
  }

  return <>{fallback}</>;
}

/**
 * CombinedVisibility Component
 *
 * Advanced component that combines role and permission checks with flexible logic.
 * Supports complex visibility rules like "owner OR (technician AND specific_permission)"
 */
export function CombinedVisibility({
  children,
  allowedRoles,
  requiredPermissions,
  fallback = null,
  loadingFallback = null,
  inverse = false,
  showOnError = false,
  requireAllRoles = false,
  requireAllPermissions = false,
  logic = "OR",
}: CombinedVisibilityProps) {
  const { hasAnyRole, hasRole, hasAnyPermission, hasPermission } =
    usePermissions();

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

  const shouldRender = inverse ? !hasAccess : hasAccess;

  if (shouldRender) {
    return <>{children}</>;
  }

  return <>{fallback}</>;
}

/**
 * ConditionalVisibility Component
 *
 * A flexible component that accepts a custom visibility function.
 * Useful for complex business logic that doesn't fit standard patterns.
 */
interface ConditionalVisibilityProps extends BaseVisibilityProps {
  condition: () => boolean;
}

export function ConditionalVisibility({
  children,
  condition,
  fallback = null,
  loadingFallback = null,
  inverse = false,
  showOnError = false,
}: ConditionalVisibilityProps) {
  const shouldRender = inverse ? !condition() : condition();

  if (shouldRender) {
    return <>{children}</>;
  }

  return <>{fallback}</>;
}

/**
 * FeatureFlag Component
 *
 * Simple feature flag component for toggling features based on roles or permissions.
 * Useful for A/B testing or gradual feature rollouts.
 */
interface FeatureFlagProps extends BaseVisibilityProps {
  feature: string;
  enabledFor?: UserRole[];
  enabledWithPermissions?: string[];
}

export function FeatureFlag({
  children,
  feature,
  enabledFor,
  enabledWithPermissions,
  fallback = null,
  loadingFallback = null,
  inverse = false,
  showOnError = false,
}: FeatureFlagProps) {
  const { hasAnyRole, hasAnyPermission } = usePermissions();

  let isEnabled = false;

  // Check if feature is enabled for user's role
  if (enabledFor && enabledFor.length > 0) {
    isEnabled = hasAnyRole(enabledFor);
  }

  // Check if feature is enabled for user's permissions
  if (enabledWithPermissions && enabledWithPermissions.length > 0) {
    isEnabled = isEnabled || hasAnyPermission(enabledWithPermissions);
  }

  const shouldRender = inverse ? !isEnabled : isEnabled;

  if (shouldRender) {
    return <>{children}</>;
  }

  return <>{fallback}</>;
}

/**
 * AdminOnly Component
 *
 * Convenience component that only shows content to admin/owner users.
 * Includes common admin permissions.
 */
export function AdminOnly({
  children,
  fallback = null,
  includeTechnicianAdmins = false,
}: {
  children: ReactNode;
  fallback?: ReactNode;
  includeTechnicianAdmins?: boolean;
}) {
  const allowedRoles: UserRole[] = ["owner"];

  if (includeTechnicianAdmins) {
    allowedRoles.push("technician");
  }

  return (
    <CombinedVisibility
      allowedRoles={allowedRoles}
      requiredPermissions={["organizations.*", "users.*"]}
      logic="OR"
      fallback={fallback}
    >
      {children}
    </CombinedVisibility>
  );
}

/**
 * ManagementOnly Component
 *
 * Convenience component for management-level features.
 * Shows content to users with management permissions.
 */
export function ManagementOnly({
  children,
  fallback = null,
}: {
  children: ReactNode;
  fallback?: ReactNode;
}) {
  return (
    <PermissionVisibility
      requiredPermissions={[
        "organizations.manage",
        "users.manage",
        "technicians.manage",
        "routes.manage",
      ]}
      requireAllPermissions={false}
      fallback={fallback}
    >
      {children}
    </PermissionVisibility>
  );
}

/**
 * ReadOnly Component
 *
 * Convenience component for read-only access.
 * Shows content to users with read permissions.
 */
export function ReadOnly({
  children,
  fallback = null,
}: {
  children: ReactNode;
  fallback?: ReactNode;
}) {
  return (
    <PermissionVisibility
      requiredPermissions={[
        "organizations.read",
        "users.read",
        "technicians.read",
        "routes.read",
      ]}
      requireAllPermissions={false}
      fallback={fallback}
    >
      {children}
    </PermissionVisibility>
  );
}
