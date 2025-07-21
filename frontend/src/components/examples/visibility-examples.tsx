"use client";

import React from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  RoleVisibility,
  PermissionVisibility,
  CombinedVisibility,
  ConditionalVisibility,
  FeatureFlag,
  AdminOnly,
  ManagementOnly,
  ReadOnly,
  RoleConditionalRender,
  PermissionConditionalRender,
  AdminConditionalRender,
} from "@/components/auth";
import { useVisibility } from "@/hooks/use-visibility";
import { useAuth } from "@/contexts/auth-context";
import {
  Users,
  UserCog,
  Route,
  Settings,
  BarChart3,
  Shield,
  Eye,
  EyeOff,
  Plus,
  Download,
  Bell,
  HelpCircle,
} from "lucide-react";

/**
 * Comprehensive example component demonstrating all visibility features
 */
export function VisibilityExamples() {
  const { user } = useAuth();
  const {
    canManageUsers,
    canManageRoutes,
    canViewAnalytics,
    canExportData,
    canAccessAdmin,
    hasElevatedPrivileges,
  } = useVisibility();

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Shield className="h-5 w-5" />
            Role-Based Visibility Examples
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {/* Basic Role Visibility */}
          <div className="grid gap-4 md:grid-cols-2">
            <RoleVisibility allowedRoles={["owner"]}>
              <div className="p-4 border rounded-lg bg-green-50">
                <h4 className="font-medium text-green-800">
                  Owner Only Content
                </h4>
                <p className="text-sm text-green-600">
                  This content is only visible to organization owners.
                </p>
              </div>
            </RoleVisibility>

            <RoleVisibility allowedRoles={["technician"]}>
              <div className="p-4 border rounded-lg bg-blue-50">
                <h4 className="font-medium text-blue-800">
                  Technician Only Content
                </h4>
                <p className="text-sm text-blue-600">
                  This content is only visible to technicians.
                </p>
              </div>
            </RoleVisibility>
          </div>

          {/* Role Visibility with Fallback */}
          <RoleVisibility
            allowedRoles={["owner"]}
            fallback={
              <div className="p-4 border rounded-lg bg-gray-50">
                <h4 className="font-medium text-gray-800">Access Restricted</h4>
                <p className="text-sm text-gray-600">
                  You need owner privileges to view this content.
                </p>
              </div>
            }
          >
            <div className="p-4 border rounded-lg bg-purple-50">
              <h4 className="font-medium text-purple-800">
                Premium Owner Content
              </h4>
              <p className="text-sm text-purple-600">
                This is premium content with a custom fallback message.
              </p>
            </div>
          </RoleVisibility>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Eye className="h-5 w-5" />
            Permission-Based Visibility Examples
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {/* Permission Visibility */}
          <div className="grid gap-4 md:grid-cols-2">
            <PermissionVisibility requiredPermissions={["users.manage"]}>
              <div className="p-4 border rounded-lg bg-orange-50">
                <h4 className="font-medium text-orange-800">User Management</h4>
                <p className="text-sm text-orange-600">
                  You can manage users in the system.
                </p>
                <Button size="sm" className="mt-2">
                  <Users className="h-4 w-4 mr-2" />
                  Manage Users
                </Button>
              </div>
            </PermissionVisibility>

            <PermissionVisibility requiredPermissions={["routes.create"]}>
              <div className="p-4 border rounded-lg bg-indigo-50">
                <h4 className="font-medium text-indigo-800">Route Creation</h4>
                <p className="text-sm text-indigo-600">
                  You can create new routes.
                </p>
                <Button size="sm" className="mt-2">
                  <Plus className="h-4 w-4 mr-2" />
                  Create Route
                </Button>
              </div>
            </PermissionVisibility>
          </div>

          {/* Multiple Permissions */}
          <PermissionVisibility
            requiredPermissions={["users.read", "routes.read"]}
            requireAllPermissions={true}
          >
            <div className="p-4 border rounded-lg bg-teal-50">
              <h4 className="font-medium text-teal-800">
                Multi-Permission Access
              </h4>
              <p className="text-sm text-teal-600">
                You have both user and route read permissions.
              </p>
            </div>
          </PermissionVisibility>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Settings className="h-5 w-5" />
            Combined Visibility Examples
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {/* Combined Role and Permission */}
          <CombinedVisibility
            allowedRoles={["owner"]}
            requiredPermissions={["organizations.manage"]}
            logic="AND"
          >
            <div className="p-4 border rounded-lg bg-red-50">
              <h4 className="font-medium text-red-800">
                Owner + Organization Management
              </h4>
              <p className="text-sm text-red-600">
                You are an owner AND have organization management permissions.
              </p>
            </div>
          </CombinedVisibility>

          {/* OR Logic */}
          <CombinedVisibility
            allowedRoles={["owner"]}
            requiredPermissions={["users.manage"]}
            logic="OR"
          >
            <div className="p-4 border rounded-lg bg-yellow-50">
              <h4 className="font-medium text-yellow-800">
                Owner OR User Management
              </h4>
              <p className="text-sm text-yellow-600">
                You are an owner OR have user management permissions.
              </p>
            </div>
          </CombinedVisibility>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <BarChart3 className="h-5 w-5" />
            Conditional Visibility Examples
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {/* Custom Condition */}
          <ConditionalVisibility
            condition={() => canManageUsers() && canManageRoutes()}
          >
            <div className="p-4 border rounded-lg bg-pink-50">
              <h4 className="font-medium text-pink-800">Custom Condition</h4>
              <p className="text-sm text-pink-600">
                You can manage both users and routes.
              </p>
            </div>
          </ConditionalVisibility>

          {/* Feature Flag */}
          <FeatureFlag
            feature="advanced-analytics"
            enabledFor={["owner"]}
            enabledWithPermissions={["analytics.*"]}
          >
            <div className="p-4 border rounded-lg bg-violet-50">
              <h4 className="font-medium text-violet-800">
                Advanced Analytics
              </h4>
              <p className="text-sm text-violet-600">
                This feature is enabled for owners or users with analytics
                permissions.
              </p>
              <Button size="sm" className="mt-2">
                <BarChart3 className="h-4 w-4 mr-2" />
                View Analytics
              </Button>
            </div>
          </FeatureFlag>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <UserCog className="h-5 w-5" />
            Convenience Components Examples
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {/* Admin Only */}
          <AdminOnly>
            <div className="p-4 border rounded-lg bg-emerald-50">
              <h4 className="font-medium text-emerald-800">
                Admin Only Section
              </h4>
              <p className="text-sm text-emerald-600">
                This section is only visible to administrators.
              </p>
              <div className="flex gap-2 mt-2">
                <Button size="sm">
                  <Settings className="h-4 w-4 mr-2" />
                  System Settings
                </Button>
                <Button size="sm" variant="outline">
                  <Shield className="h-4 w-4 mr-2" />
                  Security
                </Button>
              </div>
            </div>
          </AdminOnly>

          {/* Management Only */}
          <ManagementOnly>
            <div className="p-4 border rounded-lg bg-cyan-50">
              <h4 className="font-medium text-cyan-800">Management Section</h4>
              <p className="text-sm text-cyan-600">
                This section is for management-level users.
              </p>
              <div className="flex gap-2 mt-2">
                <Button size="sm">
                  <Users className="h-4 w-4 mr-2" />
                  Manage Team
                </Button>
                <Button size="sm" variant="outline">
                  <Route className="h-4 w-4 mr-2" />
                  Manage Routes
                </Button>
              </div>
            </div>
          </ManagementOnly>

          {/* Read Only */}
          <ReadOnly>
            <div className="p-4 border rounded-lg bg-slate-50">
              <h4 className="font-medium text-slate-800">Read-Only Access</h4>
              <p className="text-sm text-slate-600">
                You have read-only access to system data.
              </p>
              <div className="flex gap-2 mt-2">
                <Button size="sm" variant="outline">
                  <Eye className="h-4 w-4 mr-2" />
                  View Data
                </Button>
                <Button size="sm" variant="outline">
                  <Download className="h-4 w-4 mr-2" />
                  Export
                </Button>
              </div>
            </div>
          </ReadOnly>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <EyeOff className="h-5 w-5" />
            Conditional Rendering with Loading States
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {/* Role Conditional Render */}
          <RoleConditionalRender
            allowedRoles={["owner"]}
            loadingFallback={
              <div className="p-4 border rounded-lg bg-gray-100 animate-pulse">
                <div className="h-4 bg-gray-200 rounded w-1/2 mb-2"></div>
                <div className="h-3 bg-gray-200 rounded w-3/4"></div>
              </div>
            }
            fallback={
              <div className="p-4 border rounded-lg bg-gray-50">
                <h4 className="font-medium text-gray-800">Access Denied</h4>
                <p className="text-sm text-gray-600">
                  Owner privileges required.
                </p>
              </div>
            }
          >
            <div className="p-4 border rounded-lg bg-green-50">
              <h4 className="font-medium text-green-800">Owner Dashboard</h4>
              <p className="text-sm text-green-600">
                Welcome to your owner dashboard with loading states.
              </p>
            </div>
          </RoleConditionalRender>

          {/* Permission Conditional Render */}
          <PermissionConditionalRender
            requiredPermissions={["data.export"]}
            loadingFallback={
              <div className="p-4 border rounded-lg bg-gray-100 animate-pulse">
                <div className="h-4 bg-gray-200 rounded w-1/3 mb-2"></div>
                <div className="h-3 bg-gray-200 rounded w-1/2"></div>
              </div>
            }
          >
            <div className="p-4 border rounded-lg bg-blue-50">
              <h4 className="font-medium text-blue-800">Data Export</h4>
              <p className="text-sm text-blue-600">
                You have permission to export data.
              </p>
              <Button size="sm" className="mt-2">
                <Download className="h-4 w-4 mr-2" />
                Export Data
              </Button>
            </div>
          </PermissionConditionalRender>

          {/* Admin Conditional Render */}
          <AdminConditionalRender
            loadingFallback={
              <div className="p-4 border rounded-lg bg-gray-100 animate-pulse">
                <div className="h-4 bg-gray-200 rounded w-2/3 mb-2"></div>
                <div className="h-3 bg-gray-200 rounded w-1/2"></div>
              </div>
            }
          >
            <div className="p-4 border rounded-lg bg-purple-50">
              <h4 className="font-medium text-purple-800">Admin Panel</h4>
              <p className="text-sm text-purple-600">
                Advanced admin features with proper loading states.
              </p>
              <div className="flex gap-2 mt-2">
                <Button size="sm">
                  <Settings className="h-4 w-4 mr-2" />
                  Settings
                </Button>
                <Button size="sm" variant="outline">
                  <Bell className="h-4 w-4 mr-2" />
                  Notifications
                </Button>
              </div>
            </div>
          </AdminConditionalRender>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <HelpCircle className="h-5 w-5" />
            Current User Information
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            <p>
              <strong>User:</strong> {user?.first_name} {user?.last_name}
            </p>
            <p>
              <strong>Role:</strong> {user?.role}
            </p>
            <p>
              <strong>Email:</strong> {user?.email}
            </p>
            <p>
              <strong>Can Manage Users:</strong>{" "}
              {canManageUsers() ? "Yes" : "No"}
            </p>
            <p>
              <strong>Can Manage Routes:</strong>{" "}
              {canManageRoutes() ? "Yes" : "No"}
            </p>
            <p>
              <strong>Can View Analytics:</strong>{" "}
              {canViewAnalytics() ? "Yes" : "No"}
            </p>
            <p>
              <strong>Can Export Data:</strong> {canExportData() ? "Yes" : "No"}
            </p>
            <p>
              <strong>Has Elevated Privileges:</strong>{" "}
              {hasElevatedPrivileges() ? "Yes" : "No"}
            </p>
            <p>
              <strong>Can Access Admin:</strong>{" "}
              {canAccessAdmin() ? "Yes" : "No"}
            </p>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
