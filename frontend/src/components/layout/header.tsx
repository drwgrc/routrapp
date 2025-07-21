"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useAuth } from "@/contexts/auth-context";
import { Button } from "@/components/ui/button";
import { useTheme } from "@/components/theme-provider";
import { User, Settings, LogOut, ChevronDown } from "lucide-react";

export function Header() {
  const { user, isAuthenticated } = useAuth();
  const { setTheme, actualTheme } = useTheme();
  const router = useRouter();
  const [isUserMenuOpen, setIsUserMenuOpen] = useState(false);

  const handleLogout = () => {
    setIsUserMenuOpen(false);
    router.push("/logout");
  };

  const toggleTheme = () => {
    setTheme(actualTheme === "dark" ? "light" : "dark");
  };

  const toggleUserMenu = () => {
    setIsUserMenuOpen(!isUserMenuOpen);
  };

  // Close menu when clicking outside
  const closeUserMenu = () => {
    setIsUserMenuOpen(false);
  };

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container flex h-14 items-center">
        {/* Logo and Title */}
        <div className="mr-4 flex">
          <Link href="/" className="mr-6 flex items-center space-x-2">
            <div className="h-6 w-6 bg-primary rounded-sm flex items-center justify-center">
              <span className="text-primary-foreground font-bold text-sm">
                R
              </span>
            </div>
            <span className="font-bold inline-block">RoutrApp</span>
          </Link>
        </div>

        {/* Navigation */}
        {isAuthenticated && (
          <nav className="flex items-center space-x-6 text-sm font-medium">
            <Link
              href="/"
              className="transition-colors hover:text-foreground/80 text-foreground/60"
            >
              Dashboard
            </Link>
            <Link
              href="/routes"
              className="transition-colors hover:text-foreground/80 text-foreground/60"
            >
              Routes
            </Link>
            <Link
              href="/technicians"
              className="transition-colors hover:text-foreground/80 text-foreground/60"
            >
              Technicians
            </Link>
            <Link
              href="/analytics"
              className="transition-colors hover:text-foreground/80 text-foreground/60"
            >
              Analytics
            </Link>
          </nav>
        )}

        {/* Spacer */}
        <div className="flex-1" />

        {/* Right side items */}
        <div className="flex items-center space-x-4">
          {/* Theme Toggle */}
          <Button
            variant="ghost"
            size="sm"
            onClick={toggleTheme}
            className="h-8 w-8 px-0"
          >
            {actualTheme === "dark" ? (
              <svg
                className="h-4 w-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"
                />
              </svg>
            ) : (
              <svg
                className="h-4 w-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"
                />
              </svg>
            )}
          </Button>

          {/* Authentication Section */}
          {isAuthenticated ? (
            <div className="relative">
              {/* User Menu Button */}
              <button
                onClick={toggleUserMenu}
                className="flex items-center space-x-2 px-3 py-2 rounded-md hover:bg-accent transition-colors focus:outline-none focus:ring-2 focus:ring-ring"
              >
                {/* User Avatar */}
                <div className="w-8 h-8 bg-primary rounded-full flex items-center justify-center text-primary-foreground font-semibold text-sm">
                  {(user?.first_name || user?.last_name || user?.email || "U")
                    .charAt(0)
                    .toUpperCase()}
                </div>

                {/* User Info - Hidden on small screens */}
                {user && (
                  <div className="hidden sm:block text-left">
                    <p className="text-sm font-medium truncate max-w-32">
                      {user.first_name && user.last_name
                        ? `${user.first_name} ${user.last_name}`
                        : user.email}
                    </p>
                    <p className="text-xs text-muted-foreground capitalize">
                      {user.role || "User"}
                    </p>
                  </div>
                )}

                {/* Dropdown Arrow */}
                <ChevronDown
                  className={`h-4 w-4 transition-transform ${isUserMenuOpen ? "rotate-180" : ""}`}
                />
              </button>

              {/* Dropdown Menu */}
              {isUserMenuOpen && (
                <>
                  {/* Backdrop to close menu */}
                  <div className="fixed inset-0 z-40" onClick={closeUserMenu} />

                  {/* Menu Content */}
                  <div className="absolute right-0 mt-2 w-56 bg-popover border rounded-md shadow-lg z-50">
                    <div className="py-1">
                      {/* User Info - Always visible in dropdown */}
                      {user && (
                        <div className="px-4 py-2 border-b">
                          <p className="text-sm font-medium">{user.email}</p>
                          <p className="text-xs text-muted-foreground capitalize">
                            {user.role || "User"}
                          </p>
                        </div>
                      )}

                      {/* Menu Items */}
                      <Link
                        href="/profile"
                        className="flex items-center px-4 py-2 text-sm hover:bg-accent transition-colors"
                        onClick={closeUserMenu}
                      >
                        <User className="mr-3 h-4 w-4" />
                        View Profile
                      </Link>

                      <Link
                        href="/settings"
                        className="flex items-center px-4 py-2 text-sm hover:bg-accent transition-colors"
                        onClick={closeUserMenu}
                      >
                        <Settings className="mr-3 h-4 w-4" />
                        Settings
                      </Link>

                      <div className="border-t my-1" />

                      <button
                        onClick={handleLogout}
                        className="flex items-center w-full px-4 py-2 text-sm hover:bg-accent transition-colors text-left"
                      >
                        <LogOut className="mr-3 h-4 w-4" />
                        Sign Out
                      </button>
                    </div>
                  </div>
                </>
              )}
            </div>
          ) : (
            <div className="flex items-center space-x-2">
              <Button variant="ghost" size="sm" asChild>
                <Link href="/login">Sign In</Link>
              </Button>
              <Button size="sm" asChild>
                <Link href="/register">Sign Up</Link>
              </Button>
            </div>
          )}
        </div>
      </div>
    </header>
  );
}
