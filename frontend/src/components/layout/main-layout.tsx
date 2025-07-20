"use client";

import { ReactNode } from "react";
import { Header } from "./header";
import { Footer } from "./footer";
import { useAuth } from "@/contexts/auth-context";
import { Card, CardContent } from "@/components/ui/card";

interface MainLayoutProps {
  children: ReactNode;
}

export function MainLayout({ children }: MainLayoutProps) {
  const { isLoading } = useAuth();

  // Show loading state while checking authentication
  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <Card className="w-full max-w-md">
          <CardContent className="pt-6">
            <div className="flex items-center justify-center space-x-2">
              <div className="h-4 w-4 rounded-full bg-primary animate-pulse" />
              <span className="text-sm text-muted-foreground">
                Loading RoutrApp...
              </span>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex flex-col bg-background">
      <Header />
      <main className="flex-1 container mx-auto px-4 py-6">{children}</main>
      <Footer />
    </div>
  );
}
