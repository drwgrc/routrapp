import { TechnicianPage, PermissionGuard } from "@/components/auth";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { MapPin, Clock, CheckCircle, AlertCircle } from "lucide-react";

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

      <div className="grid gap-6 md:grid-cols-2 mt-6">
        <Card>
          <CardHeader>
            <CardTitle>Current Route</CardTitle>
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
            <CardTitle>Today&apos;s Schedule</CardTitle>
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
