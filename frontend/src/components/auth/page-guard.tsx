"use client";

import React from "react";
import { ProtectedRoute } from "./protected-route";
import { RouteProtectionOptions } from "@/types/auth";
import { MainLayout } from "@/components/layout";

interface PageGuardProps extends RouteProtectionOptions {
  children: React.ReactNode;
  title?: string;
  description?: string;
}

/**
 * PageGuard Component
 *
 * A wrapper component for protecting entire pages with authentication and authorization.
 * Provides a clean API for page-level protection with optional metadata.
 */
export function PageGuard({
  children,
  title,
  description,
  ...protectionOptions
}: PageGuardProps) {
  return (
    <ProtectedRoute {...protectionOptions}>
      <div className="page-container">
        {(title || description) && (
          <div className="page-header mb-6">
            {title && (
              <h1 className="text-3xl font-bold tracking-tight">{title}</h1>
            )}
            {description && (
              <p className="text-muted-foreground mt-2">{description}</p>
            )}
          </div>
        )}
        {children}
      </div>
    </ProtectedRoute>
  );
}

/**
 * OwnerPage Component
 *
 * Convenience component for pages that should only be accessible to owners
 */
export function OwnerPage({
  children,
  title,
  description,
  redirectTo = "/",
}: {
  children: React.ReactNode;
  title?: string;
  description?: string;
  redirectTo?: string;
}) {
  return (
    <PageGuard
      allowedRoles={["owner"]}
      redirectTo={redirectTo}
      title={title}
      description={description}
    >
      {children}
    </PageGuard>
  );
}

/**
 * TechnicianPage Component
 *
 * Convenience component for pages that should only be accessible to technicians
 */
export function TechnicianPage({
  children,
  title,
  description,
  redirectTo = "/",
}: {
  children: React.ReactNode;
  title?: string;
  description?: string;
  redirectTo?: string;
}) {
  return (
    <PageGuard
      allowedRoles={["technician"]}
      redirectTo={redirectTo}
      title={title}
      description={description}
    >
      {children}
    </PageGuard>
  );
}

/**
 * HierarchicalTechnicianPage Component
 *
 * Page component that allows both technicians and owners to access technician content.
 * Understands role hierarchy - owners can access technician functionality.
 */
export function HierarchicalTechnicianPage({
  children,
  title,
  description,
  redirectTo = "/",
}: {
  children: React.ReactNode;
  title?: string;
  description?: string;
  redirectTo?: string;
}) {
  return (
    <PageGuard
      allowedRoles={["technician", "owner"]}
      redirectTo={redirectTo}
      title={title}
      description={description}
    >
      {children}
    </PageGuard>
  );
}

/**
 * TechnicianPageWithLayout Component
 *
 * Technician page with proper layout (header, sidebar, footer) and role hierarchy.
 * Owners can access technician content due to role hierarchy.
 */
export function TechnicianPageWithLayout({
  children,
  title,
  description,
  redirectTo = "/",
}: {
  children: React.ReactNode;
  title?: string;
  description?: string;
  redirectTo?: string;
}) {
  return (
    <HierarchicalTechnicianPage redirectTo={redirectTo}>
      <MainLayout>
        <div className="page-container">
          {(title || description) && (
            <div className="page-header mb-6">
              {title && (
                <h1 className="text-3xl font-bold tracking-tight">{title}</h1>
              )}
              {description && (
                <p className="text-muted-foreground mt-2">{description}</p>
              )}
            </div>
          )}
          {children}
        </div>
      </MainLayout>
    </HierarchicalTechnicianPage>
  );
}

/**
 * OwnerPageWithLayout Component
 *
 * Owner/admin page with proper layout (header, sidebar, footer).
 */
export function OwnerPageWithLayout({
  children,
  title,
  description,
  redirectTo = "/",
}: {
  children: React.ReactNode;
  title?: string;
  description?: string;
  redirectTo?: string;
}) {
  return (
    <OwnerPage redirectTo={redirectTo}>
      <MainLayout>
        <div className="page-container">
          {(title || description) && (
            <div className="page-header mb-6">
              {title && (
                <h1 className="text-3xl font-bold tracking-tight">{title}</h1>
              )}
              {description && (
                <p className="text-muted-foreground mt-2">{description}</p>
              )}
            </div>
          )}
          {children}
        </div>
      </MainLayout>
    </OwnerPage>
  );
}

/**
 * AuthenticatedPage Component
 *
 * Convenience component for pages that require authentication but no specific role
 */
export function AuthenticatedPage({
  children,
  title,
  description,
  redirectTo = "/login",
}: {
  children: React.ReactNode;
  title?: string;
  description?: string;
  redirectTo?: string;
}) {
  return (
    <PageGuard
      requireAuth={true}
      redirectTo={redirectTo}
      title={title}
      description={description}
    >
      {children}
    </PageGuard>
  );
}
