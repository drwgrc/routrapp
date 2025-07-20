"use client";

import { useMemo } from "react";
import { useAuth } from "@/contexts/auth-context";
import { UserRole, PermissionCheck } from "@/types/auth";

/**
 * Default permissions for each role type
 * This mirrors the backend permission system
 */
const DEFAULT_PERMISSIONS: Record<UserRole, string[]> = {
  owner: ["organizations.*", "users.*", "technicians.*", "routes.*", "roles.*"],
  technician: [
    "routes.read",
    "routes.update_status",
    "technicians.read_own",
    "technicians.update_own",
  ],
};

/**
 * Custom hook for checking user permissions and roles
 * Provides utilities for role-based and permission-based access control
 */
export function usePermissions(): PermissionCheck {
  const { user, isAuthenticated } = useAuth();

  return useMemo(() => {
    const userRole = user?.role as UserRole | undefined;

    /**
     * Check if user has a specific role
     */
    const hasRole = (role: UserRole): boolean => {
      if (!isAuthenticated || !userRole) return false;
      return userRole === role;
    };

    /**
     * Check if user has any of the specified roles
     */
    const hasAnyRole = (roles: UserRole[]): boolean => {
      if (!isAuthenticated || !userRole) return false;
      return roles.includes(userRole);
    };

    /**
     * Check if a permission matches the requested permission
     * Supports wildcard permissions (e.g., "routes.*" matches "routes.read")
     */
    const permissionMatches = (
      storedPerm: string,
      requestedPerm: string
    ): boolean => {
      // Exact match
      if (storedPerm === requestedPerm) return true;

      // Wildcard match (e.g., "routes.*" matches "routes.read")
      if (storedPerm.endsWith(".*")) {
        const prefix = storedPerm.slice(0, -2);
        return requestedPerm.startsWith(prefix + ".");
      }

      return false;
    };

    /**
     * Check if user has a specific permission
     */
    const hasPermission = (permission: string): boolean => {
      if (!isAuthenticated || !userRole) return false;

      const rolePermissions = DEFAULT_PERMISSIONS[userRole] || [];

      return rolePermissions.some(perm => permissionMatches(perm, permission));
    };

    /**
     * Check if user has any of the specified permissions
     */
    const hasAnyPermission = (permissions: string[]): boolean => {
      if (!isAuthenticated || !userRole) return false;

      return permissions.some(permission => hasPermission(permission));
    };

    /**
     * Check if user is an owner
     */
    const isOwner = (): boolean => hasRole("owner");

    /**
     * Check if user is a technician
     */
    const isTechnician = (): boolean => hasRole("technician");

    return {
      hasRole,
      hasAnyRole,
      hasPermission,
      hasAnyPermission,
      isOwner,
      isTechnician,
    };
  }, [user, isAuthenticated]);
}
