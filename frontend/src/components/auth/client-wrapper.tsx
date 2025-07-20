"use client";

import { ReactNode } from "react";
import { AuthErrorBoundary, CombinedMiddleware } from "./index";

interface ClientAuthWrapperProps {
  children: ReactNode;
  routePermissions?: Record<string, string[]>;
}

export function ClientAuthWrapper({
  children,
  routePermissions,
}: ClientAuthWrapperProps) {
  return (
    <AuthErrorBoundary>
      <CombinedMiddleware
        routePermissions={routePermissions}
        enableSessionTimeout={true}
        enableRoleRedirect={true}
      >
        {children}
      </CombinedMiddleware>
    </AuthErrorBoundary>
  );
}
