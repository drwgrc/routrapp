import {
  TechnicianPage,
  PermissionGuard,
  RoleVisibility,
  PermissionVisibility,
  ReadOnly,
} from "@/components/auth";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  MapPin,
  Clock,
  CheckCircle,
  AlertCircle,
  Eye,
  User,
  Route,
  Settings,
} from "lucide-react";

export default function TechnicianDashboard() {
  return (
    <TechnicianPage
      title="Technician Dashboard"
      description="View your routes and update job status"
    >
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
        <StatsCard
          title="Today's Routes"
          value="8"
          icon={<MapPin className="h-4 w-4" />}
          variant="default"
        />
        <StatsCard
          title="Completed"
          value="5"
          icon={<CheckCircle className="h-4 w-4" />}
          variant="success"
        />
        <StatsCard
          title="In Progress"
          value="2"
          icon={<Clock className="h-4 w-4" />}
          variant="warning"
        />
        <StatsCard
          title="Pending"
          value="1"
          icon={<AlertCircle className="h-4 w-4" />}
          variant="destructive"
        />
      </div>

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3 mt-6">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Route className="h-4 w-4" />
              Current Route
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <p className="font-medium">Route #TEC-2024-001</p>
              <p className="text-sm text-muted-foreground">
                5 stops remaining • ETA: 2:30 PM
              </p>
            </div>

            <PermissionGuard requiredPermissions={["routes.update_status"]}>
              <Button className="w-full">Update Status</Button>
            </PermissionGuard>

            <PermissionGuard
              requiredPermissions={["routes.create"]}
              fallback={
                <p className="text-xs text-muted-foreground">
                  Only route creators can see this button
                </p>
              }
            >
              <Button variant="outline" className="w-full">
                Create New Route
              </Button>
            </PermissionGuard>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Clock className="h-4 w-4" />
              Today&apos;s Schedule
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              <ScheduleItem
                time="8:00 AM"
                location="123 Main St"
                status="completed"
              />
              <ScheduleItem
                time="9:30 AM"
                location="456 Oak Ave"
                status="completed"
              />
              <ScheduleItem
                time="11:00 AM"
                location="789 Pine Rd"
                status="in-progress"
              />
              <ScheduleItem
                time="2:30 PM"
                location="321 Elm St"
                status="pending"
              />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Eye className="h-4 w-4" />
              Access Control
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <RoleVisibility allowedRoles={["technician"]}>
              <div className="p-3 border rounded-lg bg-blue-50">
                <p className="text-sm font-medium text-blue-800">
                  ✅ Technician Access
                </p>
                <p className="text-xs text-blue-600">
                  You have technician privileges
                </p>
              </div>
            </RoleVisibility>

            <PermissionVisibility requiredPermissions={["routes.read"]}>
              <div className="p-3 border rounded-lg bg-green-50">
                <p className="text-sm font-medium text-green-800">
                  ✅ Route Access
                </p>
                <p className="text-xs text-green-600">
                  You can view and manage routes
                </p>
              </div>
            </PermissionVisibility>

            <ReadOnly>
              <div className="p-3 border rounded-lg bg-gray-50">
                <p className="text-sm font-medium text-gray-800">
                  ✅ Read Access
                </p>
                <p className="text-xs text-gray-600">
                  You have read-only access to data
                </p>
              </div>
            </ReadOnly>

            <RoleVisibility allowedRoles={["owner"]}>
              <div className="p-3 border rounded-lg bg-red-50">
                <p className="text-sm font-medium text-red-800">
                  ❌ Owner Content
                </p>
                <p className="text-xs text-red-600">
                  This won&apos;t show for technicians
                </p>
              </div>
            </RoleVisibility>
          </CardContent>
        </Card>
      </div>
    </TechnicianPage>
  );
}

interface StatsCardProps {
  title: string;
  value: string;
  icon: React.ReactNode;
  variant?: "default" | "success" | "warning" | "destructive";
}

function StatsCard({
  title,
  value,
  icon,
  variant = "default",
}: StatsCardProps) {
  const variantClasses = {
    default: "text-foreground",
    success: "text-green-600",
    warning: "text-yellow-600",
    destructive: "text-red-600",
  };

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">{title}</CardTitle>
        <div className={variantClasses[variant]}>{icon}</div>
      </CardHeader>
      <CardContent>
        <div className={`text-2xl font-bold ${variantClasses[variant]}`}>
          {value}
        </div>
      </CardContent>
    </Card>
  );
}

interface ScheduleItemProps {
  time: string;
  location: string;
  status: "completed" | "in-progress" | "pending";
}

function ScheduleItem({ time, location, status }: ScheduleItemProps) {
  const statusConfig = {
    completed: { icon: CheckCircle, color: "text-green-600" },
    "in-progress": { icon: Clock, color: "text-yellow-600" },
    pending: { icon: AlertCircle, color: "text-gray-400" },
  };

  const { icon: StatusIcon, color } = statusConfig[status];

  return (
    <div className="flex items-center space-x-3">
      <StatusIcon className={`h-4 w-4 ${color}`} />
      <div className="flex-1 min-w-0">
        <p className="text-sm font-medium">{time}</p>
        <p className="text-sm text-muted-foreground truncate">{location}</p>
      </div>
    </div>
  );
}
