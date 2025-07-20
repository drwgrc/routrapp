"use client";

import { useEffect, ReactNode } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/contexts/auth-context";
import { Card, CardContent } from "@/components/ui/card";

interface AuthGuardProps {
  children: ReactNode;
  redirectTo?: string;
  requireAuth?: boolean;
}

export function AuthGuard({
  children,
  redirectTo = "/login",
  requireAuth = true,
}: AuthGuardProps) {
  const { isAuthenticated, isLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading) {
      if (requireAuth && !isAuthenticated) {
        router.push(redirectTo);
      } else if (!requireAuth && isAuthenticated) {
        router.push("/");
      }
    }
  }, [isAuthenticated, isLoading, router, redirectTo, requireAuth]);

  // Show loading state while checking authentication
  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <Card className="w-full max-w-md">
          <CardContent className="pt-6">
            <div className="flex items-center justify-center space-x-2">
              <div className="h-4 w-4 rounded-full bg-primary animate-pulse" />
              <span className="text-sm text-muted-foreground">
                Checking authentication...
              </span>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  // Render children if authentication check passes
  if ((requireAuth && isAuthenticated) || (!requireAuth && !isAuthenticated)) {
    return <>{children}</>;
  }

  // Show loading state while redirecting
  return (
    <div className="min-h-screen flex items-center justify-center bg-background">
      <Card className="w-full max-w-md">
        <CardContent className="pt-6">
          <div className="flex items-center justify-center space-x-2">
            <div className="h-4 w-4 rounded-full bg-primary animate-pulse" />
            <span className="text-sm text-muted-foreground">
              Redirecting...
            </span>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
