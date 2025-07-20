"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";
import {
  Package2,
  Home,
  Route,
  Users,
  Settings,
  BarChart3,
  Calendar,
} from "lucide-react";

const navigation = [
  {
    name: "Dashboard",
    href: "/dashboard",
    icon: Home,
  },
  {
    name: "Routes",
    href: "/routes",
    icon: Route,
  },
  {
    name: "Technicians",
    href: "/technicians",
    icon: Users,
  },
  {
    name: "Schedules",
    href: "/schedules",
    icon: Calendar,
  },
  {
    name: "Analytics",
    href: "/analytics",
    icon: BarChart3,
  },
  {
    name: "Settings",
    href: "/settings",
    icon: Settings,
  },
];

interface SidebarProps {
  className?: string;
}

export function Sidebar({ className }: SidebarProps) {
  const pathname = usePathname();

  return (
    <div className={cn("pb-12 w-64", className)}>
      <div className="space-y-4 py-4">
        {/* Logo */}
        <div className="px-3 py-2">
          <Link href="/" className="flex items-center space-x-2">
            <Package2 className="h-6 w-6" />
            <span className="text-lg font-semibold">RoutrApp</span>
          </Link>
        </div>

        {/* Navigation */}
        <div className="px-3 py-2">
          <div className="space-y-1">
            {navigation.map(item => {
              const isActive = pathname === item.href;
              return (
                <Link
                  key={item.name}
                  href={item.href}
                  className={cn(
                    "flex items-center rounded-lg px-3 py-2 text-sm font-medium transition-colors",
                    isActive
                      ? "bg-sidebar-accent text-sidebar-accent-foreground"
                      : "text-sidebar-foreground hover:bg-sidebar-accent/50 hover:text-sidebar-accent-foreground"
                  )}
                >
                  <item.icon className="mr-3 h-4 w-4" />
                  {item.name}
                </Link>
              );
            })}
          </div>
        </div>

        {/* Bottom section */}
        <div className="px-3 py-2 mt-auto">
          <div className="rounded-lg bg-sidebar-accent/20 p-3">
            <div className="flex items-center">
              <div className="flex-1">
                <p className="text-sm font-medium text-sidebar-foreground">
                  Need help?
                </p>
                <p className="text-xs text-sidebar-foreground/70">
                  Check our documentation
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
