"use client";

import React from "react";
import { usePermissions } from "@/hooks/use-permissions";
import { UserRole } from "@/types/auth";

interface RoleGuardProps {
  children: React.ReactNode;
  allowedRoles: UserRole[];
  fallback?: React.ReactNode;
  inverse?: boolean; // If true, shows children when user DOESN'T have the roles
}

/**
 * RoleGuard Component
 *
 * A component-level guard for conditionally rendering content based on user roles.
 * Useful for hiding/showing UI elements based on permissions within a page.
 */
export function RoleGuard({
  children,
  allowedRoles,
  fallback = null,
  inverse = false,
}: RoleGuardProps) {
  const { hasAnyRole } = usePermissions();

  const hasRequiredRole = hasAnyRole(allowedRoles);
  const shouldRender = inverse ? !hasRequiredRole : hasRequiredRole;

  if (shouldRender) {
    return <>{children}</>;
  }

  return <>{fallback}</>;
}

interface PermissionGuardProps {
  children: React.ReactNode;
  requiredPermissions: string[];
  fallback?: React.ReactNode;
  inverse?: boolean; // If true, shows children when user DOESN'T have the permissions
}

/**
 * PermissionGuard Component
 *
 * A component-level guard for conditionally rendering content based on user permissions.
 * Useful for hiding/showing UI elements based on specific permissions within a page.
 */
export function PermissionGuard({
  children,
  requiredPermissions,
  fallback = null,
  inverse = false,
}: PermissionGuardProps) {
  const { hasAnyPermission } = usePermissions();

  const hasRequiredPermission = hasAnyPermission(requiredPermissions);
  const shouldRender = inverse ? !hasRequiredPermission : hasRequiredPermission;

  if (shouldRender) {
    return <>{children}</>;
  }

  return <>{fallback}</>;
}

/**
 * OwnerOnly Component
 *
 * Convenience component that only shows children to users with "owner" role
 */
export function OwnerOnly({
  children,
  fallback = null,
}: {
  children: React.ReactNode;
  fallback?: React.ReactNode;
}) {
  return (
    <RoleGuard allowedRoles={["owner"]} fallback={fallback}>
      {children}
    </RoleGuard>
  );
}

/**
 * TechnicianOnly Component
 *
 * Convenience component that only shows children to users with "technician" role
 */
export function TechnicianOnly({
  children,
  fallback = null,
}: {
  children: React.ReactNode;
  fallback?: React.ReactNode;
}) {
  return (
    <RoleGuard allowedRoles={["technician"]} fallback={fallback}>
      {children}
    </RoleGuard>
  );
}
