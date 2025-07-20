"use client";

import React, { ReactNode } from "react";
import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "@/lib/query-client";
import { AuthProvider } from "@/contexts/auth-context";

interface AppProvidersProps {
  children: ReactNode;
}

/**
 * Main app providers wrapper that combines all context providers
 * This should wrap the entire application to provide global state management
 */
export function AppProviders({ children }: AppProvidersProps) {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>{children}</AuthProvider>
    </QueryClientProvider>
  );
}
