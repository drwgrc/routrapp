"use client";

import { useMemo } from "react";
import { usePermissions } from "./use-permissions";
import { useAuth } from "@/contexts/auth-context";

/**
 * Advanced visibility logic hook
 * Provides utilities for complex visibility scenarios and business logic
 */
export function useVisibility() {
  const { hasPermission, hasAnyPermission, isOwner, isTechnician } =
    usePermissions();
  const { user, isAuthenticated } = useAuth();

  const visibilityUtils = useMemo(() => {
    /**
     * Check if user can perform a specific action on a resource
     */
    const canPerformAction = (action: string, resource: string): boolean => {
      const permission = `${resource}.${action}`;
      return hasPermission(permission);
    };

    /**
     * Check if user can manage a specific resource
     */
    const canManage = (resource: string): boolean => {
      return (
        hasPermission(`${resource}.*`) || hasPermission(`${resource}.manage`)
      );
    };

    /**
     * Check if user can read a specific resource
     */
    const canRead = (resource: string): boolean => {
      return (
        hasPermission(`${resource}.*`) || hasPermission(`${resource}.read`)
      );
    };

    /**
     * Check if user can create a specific resource
     */
    const canCreate = (resource: string): boolean => {
      return (
        hasPermission(`${resource}.*`) || hasPermission(`${resource}.create`)
      );
    };

    /**
     * Check if user can update a specific resource
     */
    const canUpdate = (resource: string): boolean => {
      return (
        hasPermission(`${resource}.*`) || hasPermission(`${resource}.update`)
      );
    };

    /**
     * Check if user can delete a specific resource
     */
    const canDelete = (resource: string): boolean => {
      return (
        hasPermission(`${resource}.*`) || hasPermission(`${resource}.delete`)
      );
    };

    /**
     * Check if user has elevated privileges (owner or admin)
     */
    const hasElevatedPrivileges = (): boolean => {
      return isOwner() || hasPermission("organizations.*");
    };

    /**
     * Check if user can access admin features
     */
    const canAccessAdmin = (): boolean => {
      return isOwner() || hasPermission("organizations.manage");
    };

    /**
     * Check if user can manage users
     */
    const canManageUsers = (): boolean => {
      return hasPermission("users.*") || hasPermission("users.manage");
    };

    /**
     * Check if user can manage technicians
     */
    const canManageTechnicians = (): boolean => {
      return (
        hasPermission("technicians.*") || hasPermission("technicians.manage")
      );
    };

    /**
     * Check if user can manage routes
     */
    const canManageRoutes = (): boolean => {
      return hasPermission("routes.*") || hasPermission("routes.manage");
    };

    /**
     * Check if user can view routes
     */
    const canViewRoutes = (): boolean => {
      return hasPermission("routes.*") || hasPermission("routes.read");
    };

    /**
     * Check if user can create routes
     */
    const canCreateRoutes = (): boolean => {
      return hasPermission("routes.*") || hasPermission("routes.create");
    };

    /**
     * Check if user can update route status
     */
    const canUpdateRouteStatus = (): boolean => {
      return hasPermission("routes.*") || hasPermission("routes.update_status");
    };

    /**
     * Check if user can access their own profile
     */
    const canAccessOwnProfile = (): boolean => {
      return isAuthenticated && user !== null;
    };

    /**
     * Check if user can access other users' profiles
     */
    const canAccessOtherProfiles = (): boolean => {
      return hasPermission("users.*") || hasPermission("users.read");
    };

    /**
     * Check if user can access organization settings
     */
    const canAccessOrgSettings = (): boolean => {
      return (
        hasPermission("organizations.*") ||
        hasPermission("organizations.manage")
      );
    };

    /**
     * Check if user can view analytics/reports
     */
    const canViewAnalytics = (): boolean => {
      return hasPermission("analytics.*") || hasPermission("analytics.read");
    };

    /**
     * Check if user can export data
     */
    const canExportData = (): boolean => {
      return hasPermission("data.*") || hasPermission("data.export");
    };

    /**
     * Check if user can import data
     */
    const canImportData = (): boolean => {
      return hasPermission("data.*") || hasPermission("data.import");
    };

    /**
     * Check if user can access system settings
     */
    const canAccessSystemSettings = (): boolean => {
      return isOwner() || hasPermission("system.*");
    };

    /**
     * Check if user can access billing/payment features
     */
    const canAccessBilling = (): boolean => {
      return isOwner() || hasPermission("billing.*");
    };

    /**
     * Check if user can access audit logs
     */
    const canAccessAuditLogs = (): boolean => {
      return hasPermission("audit.*") || hasPermission("audit.read");
    };

    /**
     * Check if user can perform bulk operations
     */
    const canPerformBulkOperations = (resource: string): boolean => {
      return (
        hasPermission(`${resource}.*`) || hasPermission(`${resource}.bulk`)
      );
    };

    /**
     * Check if user can access advanced features
     */
    const canAccessAdvancedFeatures = (): boolean => {
      return (
        isOwner() ||
        hasAnyPermission(["organizations.manage", "system.*", "advanced.*"])
      );
    };

    /**
     * Check if user can access mobile-specific features
     */
    const canAccessMobileFeatures = (): boolean => {
      return isTechnician() || hasPermission("mobile.*");
    };

    /**
     * Check if user can access desktop-specific features
     */
    const canAccessDesktopFeatures = (): boolean => {
      return isOwner() || hasPermission("desktop.*");
    };

    /**
     * Check if user can access API features
     */
    const canAccessApiFeatures = (): boolean => {
      return isOwner() || hasPermission("api.*");
    };

    /**
     * Check if user can access integration features
     */
    const canAccessIntegrationFeatures = (): boolean => {
      return isOwner() || hasPermission("integrations.*");
    };

    /**
     * Check if user can access notification settings
     */
    const canAccessNotificationSettings = (): boolean => {
      return (
        hasPermission("notifications.*") ||
        hasPermission("notifications.manage")
      );
    };

    /**
     * Check if user can access security settings
     */
    const canAccessSecuritySettings = (): boolean => {
      return isOwner() || hasPermission("security.*");
    };

    /**
     * Check if user can access backup/restore features
     */
    const canAccessBackupRestore = (): boolean => {
      return isOwner() || hasPermission("backup.*");
    };

    /**
     * Check if user can access help/support features
     */
    const canAccessHelpSupport = (): boolean => {
      return isAuthenticated; // All authenticated users can access help
    };

    /**
     * Check if user can access documentation
     */
    const canAccessDocumentation = (): boolean => {
      return isAuthenticated; // All authenticated users can access documentation
    };

    return {
      // Basic permission checks
      canPerformAction,
      canManage,
      canRead,
      canCreate,
      canUpdate,
      canDelete,

      // Role-based checks
      hasElevatedPrivileges,
      canAccessAdmin,

      // Resource-specific checks
      canManageUsers,
      canManageTechnicians,
      canManageRoutes,
      canViewRoutes,
      canCreateRoutes,
      canUpdateRouteStatus,

      // Profile and user management
      canAccessOwnProfile,
      canAccessOtherProfiles,

      // Organization and system
      canAccessOrgSettings,
      canAccessSystemSettings,
      canAccessBilling,

      // Data and analytics
      canViewAnalytics,
      canExportData,
      canImportData,
      canAccessAuditLogs,
      canPerformBulkOperations,

      // Feature access
      canAccessAdvancedFeatures,
      canAccessMobileFeatures,
      canAccessDesktopFeatures,
      canAccessApiFeatures,
      canAccessIntegrationFeatures,

      // Settings and configuration
      canAccessNotificationSettings,
      canAccessSecuritySettings,
      canAccessBackupRestore,
      canAccessHelpSupport,
      canAccessDocumentation,
    };
  }, [
    hasPermission,
    hasAnyPermission,
    isOwner,
    isTechnician,
    user,
    isAuthenticated,
  ]);

  return visibilityUtils;
}
