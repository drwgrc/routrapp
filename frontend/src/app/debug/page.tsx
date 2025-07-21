"use client";

import { useState, useEffect } from "react";
import { useAuth } from "@/contexts/auth-context";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";
import { Copy, Trash2, RefreshCw, Eye, EyeOff } from "lucide-react";
import Link from "next/link";

export default function DebugPage() {
  const { user, isAuthenticated, refreshUser } = useAuth();
  const [tokens, setTokens] = useState({
    authToken: "",
    refreshToken: "",
  });
  const [showTokens, setShowTokens] = useState(false);

  // Load tokens from localStorage on mount
  useEffect(() => {
    const authToken = localStorage.getItem("auth_token") || "";
    const refreshToken = localStorage.getItem("refresh_token") || "";
    setTokens({ authToken, refreshToken });
  }, []);

  const handleClearTokens = () => {
    localStorage.removeItem("auth_token");
    localStorage.removeItem("refresh_token");
    setTokens({ authToken: "", refreshToken: "" });
    toast.success("Tokens cleared from localStorage");
    setTimeout(() => window.location.reload(), 1000);
  };

  const handleSetTestToken = () => {
    const testAuthToken =
      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo2LCJvcmdhbml6YXRpb25faWQiOjYsImVtYWlsIjoiam9obi5kb2VAZXhhbXBsZS5jb20iLCJyb2xlIjoib3duZXIiLCJ0b2tlbl90eXBlIjoiYWNjZXNzIiwiaXNzIjoicm91dHJhcHAtYXBpIiwic3ViIjoiNiIsImF1ZCI6WyJyb3V0cmFwcC1mcm9udGVuZCJdLCJleHAiOjE3NTMwNjA4MjQsImlhdCI6MTc1MzA1OTkyNH0.i-Dl35BCnIR-6gKAL0V-V5hhFc5R8BZZ4GtDOlj122A";
    const testRefreshToken =
      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo2LCJvcmdhbml6YXRpb25faWQiOjYsImVtYWlsIjoiam9obi5kb2VAZXhhbXBsZS5jb20iLCJyb2xlIjoib3duZXIiLCJ0b2tlbl90eXBlIjoicmVmcmVzaCIsImlzcyI6InJvdXRyYXBwLWFwaSIsInN1YiI6IjYiLCJhdWQiOlsicm91dHJhcHAtZnJvbnRlbmQiXSwiZXhwIjoxNzUzNjY0NzI0LCJpYXQiOjE3NTMwNTk5MjR9.IB5gl62gScXvwcR9zUUiK1YhEkx8TDowDbKQtzLR3K0";

    localStorage.setItem("auth_token", testAuthToken);
    localStorage.setItem("refresh_token", testRefreshToken);
    setTokens({ authToken: testAuthToken, refreshToken: testRefreshToken });
    toast.success("Test tokens set in localStorage");
    setTimeout(() => window.location.reload(), 1000);
  };

  const handleCopyToken = (token: string, type: string) => {
    navigator.clipboard.writeText(token);
    toast.success(`${type} copied to clipboard`);
  };

  const handleRefreshUser = async () => {
    try {
      await refreshUser();
      toast.success("User data refreshed");
    } catch {
      toast.error("Failed to refresh user data");
    }
  };

  const formatToken = (token: string) => {
    if (!token) return "No token";
    if (token.length > 50) {
      return `${token.substring(0, 20)}...${token.substring(token.length - 20)}`;
    }
    return token;
  };

  const parseJWT = (token: string) => {
    try {
      const payload = JSON.parse(atob(token.split(".")[1]));
      return payload;
    } catch {
      return null;
    }
  };

  const getTokenStatus = (token: string) => {
    if (!token) return { status: "missing", color: "text-red-600" };

    const parts = token.split(".");
    if (parts.length !== 3)
      return { status: "malformed", color: "text-red-600" };

    const payload = parseJWT(token);
    if (!payload) return { status: "invalid", color: "text-red-600" };

    const now = Math.floor(Date.now() / 1000);
    if (payload.exp && payload.exp < now)
      return { status: "expired", color: "text-orange-600" };

    return { status: "valid", color: "text-green-600" };
  };

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="mb-6">
          <h1 className="text-3xl font-bold text-gray-900">
            Authentication Debug
          </h1>
          <p className="mt-2 text-gray-600">
            Debug authentication tokens and user state
          </p>
          <div className="mt-4">
            <Link href="/" className="text-blue-600 hover:text-blue-800">
              ← Back to Dashboard
            </Link>
          </div>
        </div>

        {/* User Status */}
        <Card className="mb-6">
          <CardHeader>
            <CardTitle>Authentication Status</CardTitle>
            <CardDescription>Current user authentication state</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <Label>Is Authenticated</Label>
                <p
                  className={`font-medium ${isAuthenticated ? "text-green-600" : "text-red-600"}`}
                >
                  {isAuthenticated ? "✅ Yes" : "❌ No"}
                </p>
              </div>
              <div>
                <Label>User Loaded</Label>
                <p
                  className={`font-medium ${user ? "text-green-600" : "text-red-600"}`}
                >
                  {user ? "✅ Yes" : "❌ No"}
                </p>
              </div>
            </div>

            {user && (
              <div className="mt-4 p-4 bg-green-50 rounded-lg">
                <h4 className="font-medium text-green-800">User Information</h4>
                <div className="mt-2 text-sm text-green-700">
                  <p>Email: {user.email}</p>
                  <p>
                    Name: {user.first_name} {user.last_name}
                  </p>
                  <p>Role: {user.role}</p>
                  <p>ID: {user.id}</p>
                </div>
              </div>
            )}

            <div className="flex gap-2">
              <Button onClick={handleRefreshUser} variant="outline" size="sm">
                <RefreshCw className="mr-2 h-4 w-4" />
                Refresh User
              </Button>
              <Link href="/profile">
                <Button size="sm">Go to Profile</Button>
              </Link>
            </div>
          </CardContent>
        </Card>

        {/* Token Management */}
        <Card className="mb-6">
          <CardHeader>
            <CardTitle className="flex items-center justify-between">
              Token Management
              <Button
                variant="outline"
                size="sm"
                onClick={() => setShowTokens(!showTokens)}
              >
                {showTokens ? (
                  <EyeOff className="h-4 w-4" />
                ) : (
                  <Eye className="h-4 w-4" />
                )}
                {showTokens ? "Hide" : "Show"} Tokens
              </Button>
            </CardTitle>
            <CardDescription>
              Manage authentication tokens stored in localStorage
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            {/* Auth Token */}
            <div>
              <div className="flex items-center justify-between mb-2">
                <Label>Access Token</Label>
                <span
                  className={`text-sm font-medium ${getTokenStatus(tokens.authToken).color}`}
                >
                  {getTokenStatus(tokens.authToken).status}
                </span>
              </div>
              <div className="flex gap-2">
                <Input
                  value={
                    showTokens
                      ? tokens.authToken
                      : formatToken(tokens.authToken)
                  }
                  readOnly
                  className="font-mono text-xs"
                />
                {tokens.authToken && (
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() =>
                      handleCopyToken(tokens.authToken, "Access token")
                    }
                  >
                    <Copy className="h-4 w-4" />
                  </Button>
                )}
              </div>
            </div>

            {/* Refresh Token */}
            <div>
              <div className="flex items-center justify-between mb-2">
                <Label>Refresh Token</Label>
                <span
                  className={`text-sm font-medium ${getTokenStatus(tokens.refreshToken).color}`}
                >
                  {getTokenStatus(tokens.refreshToken).status}
                </span>
              </div>
              <div className="flex gap-2">
                <Input
                  value={
                    showTokens
                      ? tokens.refreshToken
                      : formatToken(tokens.refreshToken)
                  }
                  readOnly
                  className="font-mono text-xs"
                />
                {tokens.refreshToken && (
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() =>
                      handleCopyToken(tokens.refreshToken, "Refresh token")
                    }
                  >
                    <Copy className="h-4 w-4" />
                  </Button>
                )}
              </div>
            </div>

            {/* Actions */}
            <div className="flex gap-2 pt-4">
              <Button onClick={handleSetTestToken} variant="default" size="sm">
                Set Test Tokens
              </Button>
              <Button
                onClick={handleClearTokens}
                variant="destructive"
                size="sm"
              >
                <Trash2 className="mr-2 h-4 w-4" />
                Clear Tokens
              </Button>
            </div>

            <div className="text-sm text-gray-600">
              <p className="font-medium">Test User Credentials:</p>
              <p>Email: john.doe@example.com</p>
              <p>Password: password123</p>
            </div>
          </CardContent>
        </Card>

        {/* Token Details */}
        {showTokens && tokens.authToken && (
          <Card>
            <CardHeader>
              <CardTitle>Token Details</CardTitle>
              <CardDescription>Decoded JWT token information</CardDescription>
            </CardHeader>
            <CardContent>
              {(() => {
                const payload = parseJWT(tokens.authToken);
                return payload ? (
                  <pre className="bg-gray-100 p-4 rounded-lg text-sm overflow-auto">
                    {JSON.stringify(payload, null, 2)}
                  </pre>
                ) : (
                  <p className="text-red-600">Unable to parse token</p>
                );
              })()}
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}
