"use client";

import { useState } from "react";
import { Header } from "./header";
import { Sidebar } from "./sidebar";
import { Footer } from "./footer";
import { cn } from "@/lib/utils";

interface MainLayoutProps {
  children: React.ReactNode;
  showSidebar?: boolean;
}

export function MainLayout({ children, showSidebar = true }: MainLayoutProps) {
  const [sidebarOpen, setSidebarOpen] = useState(false);

  return (
    <div className="min-h-screen bg-background">
      <Header setSidebarOpen={setSidebarOpen} />

      <div className="flex">
        {showSidebar && (
          <>
            {/* Desktop Sidebar */}
            <aside className="hidden md:flex md:w-64 md:flex-col md:fixed md:inset-y-0 md:top-16">
              <div className="flex-1 flex flex-col bg-sidebar border-r border-sidebar-border">
                <Sidebar />
              </div>
            </aside>

            {/* Mobile Sidebar */}
            {sidebarOpen && (
              <div className="fixed inset-0 z-50 md:hidden">
                <div
                  className="fixed inset-0 bg-background/80 backdrop-blur-sm"
                  onClick={() => setSidebarOpen(false)}
                />
                <aside className="fixed inset-y-0 left-0 z-50 w-64 bg-sidebar border-r border-sidebar-border">
                  <Sidebar />
                </aside>
              </div>
            )}
          </>
        )}

        {/* Main Content */}
        <main
          className={cn(
            "flex-1 min-h-[calc(100vh-4rem)]",
            showSidebar && "md:ml-64"
          )}
        >
          <div className="container mx-auto px-4 py-6">{children}</div>
        </main>
      </div>

      <Footer />
    </div>
  );
}
