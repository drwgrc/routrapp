"use client";

import React, { Component, ReactNode } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { AlertTriangle, RefreshCw, LogOut } from "lucide-react";
import { AuthError } from "@/types/auth";

interface AuthErrorBoundaryState {
  hasError: boolean;
  error: AuthError | null;
}

interface AuthErrorBoundaryProps {
  children: ReactNode;
  fallback?: (error: AuthError, retry: () => void) => ReactNode;
  onError?: (error: AuthError) => void;
  onLogout?: () => void;
}

/**
 * AuthErrorBoundary Component
 *
 * Error boundary specifically designed for handling authentication-related errors.
 * Provides appropriate UI for different types of auth errors and recovery options.
 */
export class AuthErrorBoundary extends Component<
  AuthErrorBoundaryProps,
  AuthErrorBoundaryState
> {
  constructor(props: AuthErrorBoundaryProps) {
    super(props);
    this.state = {
      hasError: false,
      error: null,
    };
  }

  static getDerivedStateFromError(error: Error): AuthErrorBoundaryState {
    // Convert generic error to AuthError
    const authError: AuthError = {
      type: "UNKNOWN",
      message: error.message,
      code: "UNKNOWN_ERROR",
    };

    // Try to determine error type from error message or properties
    if (
      error.message.includes("401") ||
      error.message.includes("Unauthorized")
    ) {
      authError.type = "UNAUTHORIZED";
      authError.statusCode = 401;
    } else if (
      error.message.includes("403") ||
      error.message.includes("Forbidden")
    ) {
      authError.type = "FORBIDDEN";
      authError.statusCode = 403;
    } else if (
      error.message.includes("token") &&
      error.message.includes("expired")
    ) {
      authError.type = "TOKEN_EXPIRED";
      authError.statusCode = 401;
    }

    return {
      hasError: true,
      error: authError,
    };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error("AuthErrorBoundary caught an error:", error, errorInfo);

    if (this.props.onError && this.state.error) {
      this.props.onError(this.state.error);
    }
  }

  handleRetry = () => {
    this.setState({
      hasError: false,
      error: null,
    });
  };

  handleLogout = () => {
    if (this.props.onLogout) {
      this.props.onLogout();
    }
    this.handleRetry();
  };

  render() {
    if (this.state.hasError && this.state.error) {
      // Use custom fallback if provided
      if (this.props.fallback) {
        return this.props.fallback(this.state.error, this.handleRetry);
      }

      // Default error UI
      return (
        <DefaultAuthErrorUI
          error={this.state.error}
          onRetry={this.handleRetry}
          onLogout={this.handleLogout}
        />
      );
    }

    return this.props.children;
  }
}

/**
 * Default authentication error UI component
 */
interface DefaultAuthErrorUIProps {
  error: AuthError;
  onRetry: () => void;
  onLogout: () => void;
}

function DefaultAuthErrorUI({
  error,
  onRetry,
  onLogout,
}: DefaultAuthErrorUIProps) {
  const getErrorContent = () => {
    switch (error.type) {
      case "UNAUTHORIZED":
        return {
          title: "Authentication Required",
          message:
            "Your session has expired or you're not authorized to access this resource.",
          icon: <LogOut className="h-12 w-12 text-destructive" />,
          actions: (
            <>
              <Button onClick={onLogout} className="w-full">
                <LogOut className="h-4 w-4 mr-2" />
                Login Again
              </Button>
              <Button variant="outline" onClick={onRetry} className="w-full">
                <RefreshCw className="h-4 w-4 mr-2" />
                Retry
              </Button>
            </>
          ),
        };

      case "FORBIDDEN":
        return {
          title: "Access Forbidden",
          message: "You don't have permission to access this resource.",
          icon: <AlertTriangle className="h-12 w-12 text-destructive" />,
          actions: (
            <>
              <Button variant="outline" onClick={onRetry} className="w-full">
                <RefreshCw className="h-4 w-4 mr-2" />
                Retry
              </Button>
              <Button variant="outline" onClick={onLogout} className="w-full">
                <LogOut className="h-4 w-4 mr-2" />
                Switch Account
              </Button>
            </>
          ),
        };

      case "TOKEN_EXPIRED":
        return {
          title: "Session Expired",
          message: "Your session has expired. Please log in again to continue.",
          icon: <LogOut className="h-12 w-12 text-destructive" />,
          actions: (
            <Button onClick={onLogout} className="w-full">
              <LogOut className="h-4 w-4 mr-2" />
              Login Again
            </Button>
          ),
        };

      default:
        return {
          title: "Authentication Error",
          message: error.message || "An unknown authentication error occurred.",
          icon: <AlertTriangle className="h-12 w-12 text-destructive" />,
          actions: (
            <>
              <Button variant="outline" onClick={onRetry} className="w-full">
                <RefreshCw className="h-4 w-4 mr-2" />
                Retry
              </Button>
              <Button variant="outline" onClick={onLogout} className="w-full">
                <LogOut className="h-4 w-4 mr-2" />
                Logout
              </Button>
            </>
          ),
        };
    }
  };

  const { title, message, icon, actions } = getErrorContent();

  return (
    <div className="min-h-screen flex items-center justify-center bg-background p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <div className="flex justify-center mb-4">{icon}</div>
          <CardTitle className="text-xl text-destructive">{title}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <p className="text-sm text-muted-foreground text-center">{message}</p>
          <div className="space-y-2">{actions}</div>
          {error.statusCode && (
            <p className="text-xs text-muted-foreground text-center">
              Error Code: {error.statusCode}
            </p>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

/**
 * Hook for using the AuthErrorBoundary in functional components
 */
export function useAuthErrorBoundary() {
  const throwAuthError = (error: AuthError) => {
    const err = new Error(error.message);
    err.name = `AUTH_${error.type}`;
    throw err;
  };

  return { throwAuthError };
}
