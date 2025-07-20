import {
  OwnerPage,
  RoleGuard,
  AdminOnly,
  ManagementOnly,
  RoleVisibility,
  PermissionVisibility,
} from "@/components/auth";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Users,
  UserCog,
  Route,
  BarChart3,
  Settings,
  Shield,
  Eye,
} from "lucide-react";

export default function AdminDashboard() {
  return (
    <OwnerPage
      title="Admin Dashboard"
      description="Manage your organization, users, and routes"
    >
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
        <StatsCard
          title="Total Users"
          value="24"
          icon={<Users className="h-4 w-4" />}
        />
        <StatsCard
          title="Technicians"
          value="18"
          icon={<UserCog className="h-4 w-4" />}
        />
        <StatsCard
          title="Active Routes"
          value="156"
          icon={<Route className="h-4 w-4" />}
        />
        <StatsCard
          title="Efficiency"
          value="94.2%"
          icon={<BarChart3 className="h-4 w-4" />}
        />
      </div>

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3 mt-6">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Shield className="h-4 w-4" />
              Admin Features
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <AdminOnly>
              <div className="p-3 border rounded-lg bg-green-50">
                <p className="text-sm font-medium text-green-800">
                  ✅ Admin Panel Access
                </p>
                <p className="text-xs text-green-600">
                  You have full administrative privileges
                </p>
              </div>
            </AdminOnly>

            <ManagementOnly>
              <div className="p-3 border rounded-lg bg-blue-50">
                <p className="text-sm font-medium text-blue-800">
                  ✅ Management Access
                </p>
                <p className="text-xs text-blue-600">
                  You can manage users and routes
                </p>
              </div>
            </ManagementOnly>

            <RoleVisibility allowedRoles={["technician"]}>
              <div className="p-3 border rounded-lg bg-gray-50">
                <p className="text-sm font-medium text-gray-800">
                  ❌ Technician Content
                </p>
                <p className="text-xs text-gray-600">
                  This won&apos;t show for owners
                </p>
              </div>
            </RoleVisibility>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Settings className="h-4 w-4" />
              System Controls
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <PermissionVisibility
              requiredPermissions={["organizations.manage"]}
            >
              <div className="p-3 border rounded-lg bg-purple-50">
                <p className="text-sm font-medium text-purple-800">
                  ✅ Organization Management
                </p>
                <p className="text-xs text-purple-600">
                  You can manage organization settings
                </p>
              </div>
            </PermissionVisibility>

            <PermissionVisibility requiredPermissions={["users.manage"]}>
              <div className="p-3 border rounded-lg bg-orange-50">
                <p className="text-sm font-medium text-orange-800">
                  ✅ User Management
                </p>
                <p className="text-xs text-orange-600">
                  You can manage user accounts
                </p>
              </div>
            </PermissionVisibility>

            <PermissionVisibility requiredPermissions={["system.*"]}>
              <div className="p-3 border rounded-lg bg-red-50">
                <p className="text-sm font-medium text-red-800">
                  ✅ System Access
                </p>
                <p className="text-xs text-red-600">
                  You have system-level access
                </p>
              </div>
            </PermissionVisibility>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Eye className="h-4 w-4" />
              Recent Activity
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">
              Recent organization activity will appear here.
            </p>
          </CardContent>
        </Card>
      </div>
    </OwnerPage>
  );
}

interface StatsCardProps {
  title: string;
  value: string;
  icon: React.ReactNode;
}

function StatsCard({ title, value, icon }: StatsCardProps) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">{title}</CardTitle>
        {icon}
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold">{value}</div>
      </CardContent>
    </Card>
  );
}
