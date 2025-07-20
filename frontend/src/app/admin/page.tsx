import { OwnerPage, RoleGuard } from "@/components/auth";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Users, UserCog, Route, BarChart3 } from "lucide-react";

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

      <div className="grid gap-6 md:grid-cols-2 mt-6">
        <Card>
          <CardHeader>
            <CardTitle>Quick Actions</CardTitle>
          </CardHeader>
          <CardContent className="space-y-2">
            <RoleGuard allowedRoles={["owner"]}>
              <p className="text-sm text-muted-foreground">
                ✅ You can see this because you&apos;re an owner
              </p>
            </RoleGuard>
            <RoleGuard allowedRoles={["technician"]}>
              <p className="text-sm text-muted-foreground">
                ❌ This won&apos;t show for owners
              </p>
            </RoleGuard>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Recent Activity</CardTitle>
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
