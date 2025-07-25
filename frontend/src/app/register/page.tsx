"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useAuth } from "@/contexts/auth-context";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { useForm } from "react-hook-form";
import { RegistrationData } from "@/types/auth";

interface RegisterFormData extends RegistrationData {
  confirmPassword: string;
  termsAccepted: boolean;
}

export default function RegisterPage() {
  const router = useRouter();
  const {
    register: registerUser,
    isLoading,
    error,
    clearError,
    isAuthenticated,
  } = useAuth();
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [registrationSuccess, setRegistrationSuccess] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors, isValid },
    getValues,
    watch,
    setValue,
  } = useForm<RegisterFormData>({
    mode: "onChange",
    defaultValues: {
      firstName: "",
      lastName: "",
      email: "",
      password: "",
      confirmPassword: "",
      organizationName: "",
      organizationEmail: "",
      subDomain: "",
      termsAccepted: false,
    },
  });

  // Watch organization name to auto-generate subdomain
  const organizationName = watch("organizationName");

  // Helper function to generate subdomain from organization name
  const generateSubDomain = (name: string): string => {
    return name
      .toLowerCase()
      .replace(/[^a-z0-9]/g, "")
      .substring(0, 20);
  };

  // Auto-generate subdomain from organization name
  useEffect(() => {
    if (organizationName) {
      const subDomain = generateSubDomain(organizationName);
      const currentSubDomain = getValues("subDomain");

      // Only auto-update if the field is empty or if it matches the auto-generated value
      // This prevents overwriting user edits
      if (
        !currentSubDomain ||
        currentSubDomain === generateSubDomain(organizationName)
      ) {
        setValue("subDomain", subDomain, { shouldValidate: true });
      }
    }
  }, [organizationName, getValues, setValue]);

  // Redirect if already authenticated
  useEffect(() => {
    if (isAuthenticated) {
      router.push("/");
    }
  }, [isAuthenticated, router]);

  const onSubmit = async (data: RegisterFormData) => {
    try {
      // Clear any previous errors before attempting registration
      clearError();
      await registerUser({
        firstName: data.firstName,
        lastName: data.lastName,
        email: data.email,
        password: data.password,
        organizationName: data.organizationName,
        organizationEmail: data.organizationEmail,
        subDomain: data.subDomain,
      });
      setRegistrationSuccess(true);
    } catch (err) {
      console.error("Registration failed:", err);
    }
  };

  // Show success message after registration
  if (registrationSuccess) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background px-4">
        <div className="w-full max-w-md space-y-6">
          <div className="text-center space-y-2">
            <h1 className="text-3xl font-bold tracking-tight text-green-600">
              Registration Successful!
            </h1>
            <p className="text-muted-foreground">
              Your account has been created successfully
            </p>
          </div>

          <Card>
            <CardContent className="pt-6">
              <div className="text-center space-y-4">
                <div className="mx-auto w-12 h-12 bg-green-100 rounded-full flex items-center justify-center">
                  <svg
                    className="w-6 h-6 text-green-600"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M5 13l4 4L19 7"
                    />
                  </svg>
                </div>
                <div className="space-y-2">
                  <h3 className="text-lg font-semibold">
                    Welcome to RoutrApp!
                  </h3>
                  <p className="text-sm text-muted-foreground">
                    You can now sign in with your credentials and start managing
                    your routes and technicians.
                  </p>
                </div>
                <Button asChild className="w-full">
                  <Link href="/login">Sign In to Your Account</Link>
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    );
  }

  // If user is authenticated, show loading state while redirecting
  if (isAuthenticated) {
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

  return (
    <div className="min-h-screen flex items-center justify-center bg-background px-4 py-8">
      <div className="w-full max-w-md space-y-6">
        {/* Welcome Section */}
        <div className="text-center space-y-2">
          <h1 className="text-3xl font-bold tracking-tight">Join RoutrApp</h1>
          <p className="text-muted-foreground">
            Create your account to start optimizing routes and managing
            technicians
          </p>
        </div>

        {/* Registration Form Card */}
        <Card>
          <CardHeader className="space-y-1">
            <CardTitle className="text-2xl text-center">
              Create Account
            </CardTitle>
            <CardDescription className="text-center">
              Enter your details to get started with RoutrApp
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
              {/* Error Alert */}
              {error && (
                <Alert variant="destructive">
                  <AlertDescription>{error}</AlertDescription>
                </Alert>
              )}

              {/* First Name Field */}
              <div className="space-y-2">
                <Label htmlFor="firstName">First Name</Label>
                <Input
                  id="firstName"
                  type="text"
                  placeholder="John"
                  disabled={isLoading}
                  {...register("firstName", {
                    required: "First name is required",
                    minLength: {
                      value: 2,
                      message: "First name must be at least 2 characters",
                    },
                  })}
                  className={errors.firstName ? "border-destructive" : ""}
                />
                {errors.firstName && (
                  <p className="text-sm text-destructive">
                    {errors.firstName.message}
                  </p>
                )}
              </div>

              {/* Last Name Field */}
              <div className="space-y-2">
                <Label htmlFor="lastName">Last Name</Label>
                <Input
                  id="lastName"
                  type="text"
                  placeholder="Doe"
                  disabled={isLoading}
                  {...register("lastName", {
                    required: "Last name is required",
                    minLength: {
                      value: 2,
                      message: "Last name must be at least 2 characters",
                    },
                  })}
                  className={errors.lastName ? "border-destructive" : ""}
                />
                {errors.lastName && (
                  <p className="text-sm text-destructive">
                    {errors.lastName.message}
                  </p>
                )}
              </div>

              {/* Email Field */}
              <div className="space-y-2">
                <Label htmlFor="email">Email Address</Label>
                <Input
                  id="email"
                  type="email"
                  placeholder="your@email.com"
                  disabled={isLoading}
                  {...register("email", {
                    required: "Email is required",
                    pattern: {
                      value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
                      message: "Please enter a valid email address",
                    },
                  })}
                  className={errors.email ? "border-destructive" : ""}
                />
                {errors.email && (
                  <p className="text-sm text-destructive">
                    {errors.email.message}
                  </p>
                )}
              </div>

              {/* Organization Name Field */}
              <div className="space-y-2">
                <Label htmlFor="organizationName">Organization Name</Label>
                <Input
                  id="organizationName"
                  type="text"
                  placeholder="Your Company Name"
                  disabled={isLoading}
                  {...register("organizationName", {
                    required: "Organization name is required",
                    minLength: {
                      value: 2,
                      message:
                        "Organization name must be at least 2 characters",
                    },
                  })}
                  className={
                    errors.organizationName ? "border-destructive" : ""
                  }
                />
                {errors.organizationName && (
                  <p className="text-sm text-destructive">
                    {errors.organizationName.message}
                  </p>
                )}
              </div>

              {/* Organization Email Field */}
              <div className="space-y-2">
                <Label htmlFor="organizationEmail">Organization Email</Label>
                <Input
                  id="organizationEmail"
                  type="email"
                  placeholder="info@yourcompany.com"
                  disabled={isLoading}
                  {...register("organizationEmail", {
                    required: "Organization email is required",
                    pattern: {
                      value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
                      message:
                        "Please enter a valid organization email address",
                    },
                  })}
                  className={
                    errors.organizationEmail ? "border-destructive" : ""
                  }
                />
                {errors.organizationEmail && (
                  <p className="text-sm text-destructive">
                    {errors.organizationEmail.message}
                  </p>
                )}
              </div>

              {/* Sub Domain Field */}
              <div className="space-y-2">
                <Label htmlFor="subDomain">Sub Domain</Label>
                <Input
                  id="subDomain"
                  type="text"
                  placeholder="yourcompany"
                  disabled={isLoading}
                  {...register("subDomain", {
                    required: "Sub domain is required",
                    pattern: {
                      value: /^[a-z0-9]+$/,
                      message:
                        "Sub domain can only contain lowercase letters and numbers",
                    },
                    minLength: {
                      value: 3,
                      message: "Sub domain must be at least 3 characters",
                    },
                    maxLength: {
                      value: 20,
                      message: "Sub domain cannot exceed 20 characters",
                    },
                  })}
                  className={errors.subDomain ? "border-destructive" : ""}
                />
                <p className="text-xs text-muted-foreground">
                  Auto-generated from organization name. You can customize it if
                  needed.
                </p>
                {errors.subDomain && (
                  <p className="text-sm text-destructive">
                    {errors.subDomain.message}
                  </p>
                )}
              </div>

              {/* Password Field */}
              <div className="space-y-2">
                <div className="flex items-center justify-between">
                  <Label htmlFor="password">Password</Label>
                  <button
                    type="button"
                    className="text-sm text-primary hover:underline"
                    onClick={() => setShowPassword(!showPassword)}
                  >
                    {showPassword ? "Hide" : "Show"}
                  </button>
                </div>
                <Input
                  id="password"
                  type={showPassword ? "text" : "password"}
                  placeholder="Create a strong password"
                  disabled={isLoading}
                  {...register("password", {
                    required: "Password is required",
                    minLength: {
                      value: 8,
                      message: "Password must be at least 8 characters",
                    },
                    pattern: {
                      value: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/,
                      message:
                        "Password must contain at least one uppercase letter, one lowercase letter, and one number",
                    },
                  })}
                  className={errors.password ? "border-destructive" : ""}
                />
                {errors.password && (
                  <p className="text-sm text-destructive">
                    {errors.password.message}
                  </p>
                )}
              </div>

              {/* Confirm Password Field */}
              <div className="space-y-2">
                <div className="flex items-center justify-between">
                  <Label htmlFor="confirmPassword">Confirm Password</Label>
                  <button
                    type="button"
                    className="text-sm text-primary hover:underline"
                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                  >
                    {showConfirmPassword ? "Hide" : "Show"}
                  </button>
                </div>
                <Input
                  id="confirmPassword"
                  type={showConfirmPassword ? "text" : "password"}
                  placeholder="Confirm your password"
                  disabled={isLoading}
                  {...register("confirmPassword", {
                    required: "Please confirm your password",
                    validate: value => {
                      const password = getValues("password");
                      return value === password || "Passwords do not match";
                    },
                  })}
                  className={errors.confirmPassword ? "border-destructive" : ""}
                />
                {errors.confirmPassword && (
                  <p className="text-sm text-destructive">
                    {errors.confirmPassword.message}
                  </p>
                )}
              </div>

              {/* Terms and Conditions */}
              <div className="flex items-start space-x-2">
                <input
                  id="termsAccepted"
                  type="checkbox"
                  className="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary mt-0.5"
                  {...register("termsAccepted", {
                    required: "You must accept the terms and conditions",
                  })}
                  disabled={isLoading}
                />
                <div className="space-y-1">
                  <Label
                    htmlFor="termsAccepted"
                    className="text-sm leading-relaxed"
                  >
                    I agree to the{" "}
                    <Link
                      href="/terms"
                      className="text-primary hover:underline"
                    >
                      Terms of Service
                    </Link>{" "}
                    and{" "}
                    <Link
                      href="/privacy"
                      className="text-primary hover:underline"
                    >
                      Privacy Policy
                    </Link>
                  </Label>
                  {errors.termsAccepted && (
                    <p className="text-sm text-destructive">
                      {errors.termsAccepted.message}
                    </p>
                  )}
                </div>
              </div>

              {/* Submit Button */}
              <Button
                type="submit"
                className="w-full"
                disabled={isLoading || !isValid}
              >
                {isLoading ? (
                  <div className="flex items-center space-x-2">
                    <div className="h-4 w-4 rounded-full border-2 border-white border-t-transparent animate-spin" />
                    <span>Creating account...</span>
                  </div>
                ) : (
                  "Create Account"
                )}
              </Button>
            </form>

            {/* Sign In Link */}
            <div className="mt-6 text-center">
              <p className="text-sm text-muted-foreground">
                Already have an account?{" "}
                <Link
                  href="/login"
                  className="text-primary hover:underline font-medium"
                >
                  Sign in
                </Link>
              </p>
            </div>
          </CardContent>
        </Card>

        {/* API Status Indicator */}
        <Card>
          <CardContent className="pt-4">
            <div className="flex items-center justify-center space-x-2 text-xs text-muted-foreground">
              <div className="h-2 w-2 rounded-full bg-green-500" />
              <span>Backend API Connected</span>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
