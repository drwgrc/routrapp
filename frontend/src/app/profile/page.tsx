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
import {
  ArrowLeft,
  User,
  Mail,
  Shield,
  Calendar,
  Edit3,
  Save,
  X,
} from "lucide-react";
import Link from "next/link";

export default function ProfilePage() {
  const { user, updateProfile, isLoading } = useAuth();
  const [isEditing, setIsEditing] = useState(false);
  const [formData, setFormData] = useState({
    firstName: user?.first_name || "",
    lastName: user?.last_name || "",
  });
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [isSubmitting, setIsSubmitting] = useState(false);

  // Update form data when user data changes
  useEffect(() => {
    if (user) {
      setFormData({
        firstName: user.first_name || "",
        lastName: user.last_name || "",
      });
    }
  }, [user]);

  // Validation function
  const validateForm = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.firstName.trim()) {
      newErrors.firstName = "First name is required";
    } else if (formData.firstName.trim().length < 1) {
      newErrors.firstName = "First name must be at least 1 character";
    } else if (formData.firstName.trim().length > 100) {
      newErrors.firstName = "First name must be less than 100 characters";
    }

    if (!formData.lastName.trim()) {
      newErrors.lastName = "Last name is required";
    } else if (formData.lastName.trim().length < 1) {
      newErrors.lastName = "Last name must be at least 1 character";
    } else if (formData.lastName.trim().length > 100) {
      newErrors.lastName = "Last name must be less than 100 characters";
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  // Handle form submission
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) {
      return;
    }

    setIsSubmitting(true);
    try {
      await updateProfile({
        firstName: formData.firstName.trim(),
        lastName: formData.lastName.trim(),
      });

      toast.success("Profile updated successfully!");
      setIsEditing(false);
    } catch (error) {
      console.error("Profile update error:", error);
      toast.error("Failed to update profile. Please try again.");
    } finally {
      setIsSubmitting(false);
    }
  };

  // Handle edit mode toggle
  const handleEditToggle = () => {
    if (isEditing) {
      // Reset form data to current user data when canceling
      setFormData({
        firstName: user?.first_name || "",
        lastName: user?.last_name || "",
      });
      setErrors({});
    }
    setIsEditing(!isEditing);
  };

  // Handle input changes
  const handleInputChange = (field: string, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }));
    // Clear error for this field when user starts typing
    if (errors[field]) {
      setErrors(prev => ({ ...prev, [field]: "" }));
    }
  };

  // Format date for display
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  };

  if (isLoading || !user) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center p-4">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-4 sm:py-8">
      <div className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="mb-6">
          <Link
            href="/"
            className="inline-flex items-center text-sm text-gray-500 hover:text-gray-700 mb-4"
          >
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Dashboard
          </Link>
          <h1 className="text-2xl sm:text-3xl font-bold text-gray-900">
            User Profile
          </h1>
          <p className="mt-2 text-gray-600">
            Manage your account information and settings
          </p>
        </div>

        {/* Profile Information Card */}
        <Card className="mb-6">
          <CardHeader className="pb-4">
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="flex items-center text-lg">
                  <User className="mr-2 h-5 w-5" />
                  Personal Information
                </CardTitle>
                <CardDescription>
                  Update your personal details and profile information
                </CardDescription>
              </div>
              <Button
                variant={isEditing ? "outline" : "default"}
                size="sm"
                onClick={handleEditToggle}
                disabled={isSubmitting}
                className="flex items-center"
              >
                {isEditing ? (
                  <>
                    <X className="mr-2 h-4 w-4" />
                    Cancel
                  </>
                ) : (
                  <>
                    <Edit3 className="mr-2 h-4 w-4" />
                    Edit
                  </>
                )}
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                {/* First Name */}
                <div className="space-y-2">
                  <Label htmlFor="firstName">First Name</Label>
                  {isEditing ? (
                    <div>
                      <Input
                        id="firstName"
                        type="text"
                        value={formData.firstName}
                        onChange={e =>
                          handleInputChange("firstName", e.target.value)
                        }
                        className={errors.firstName ? "border-red-500" : ""}
                        placeholder="Enter your first name"
                        disabled={isSubmitting}
                      />
                      {errors.firstName && (
                        <p className="text-sm text-red-600 mt-1">
                          {errors.firstName}
                        </p>
                      )}
                    </div>
                  ) : (
                    <p className="text-gray-900 font-medium">
                      {user.first_name}
                    </p>
                  )}
                </div>

                {/* Last Name */}
                <div className="space-y-2">
                  <Label htmlFor="lastName">Last Name</Label>
                  {isEditing ? (
                    <div>
                      <Input
                        id="lastName"
                        type="text"
                        value={formData.lastName}
                        onChange={e =>
                          handleInputChange("lastName", e.target.value)
                        }
                        className={errors.lastName ? "border-red-500" : ""}
                        placeholder="Enter your last name"
                        disabled={isSubmitting}
                      />
                      {errors.lastName && (
                        <p className="text-sm text-red-600 mt-1">
                          {errors.lastName}
                        </p>
                      )}
                    </div>
                  ) : (
                    <p className="text-gray-900 font-medium">
                      {user.last_name}
                    </p>
                  )}
                </div>
              </div>

              {/* Save Button - only show when editing */}
              {isEditing && (
                <div className="flex justify-end pt-4">
                  <Button
                    type="submit"
                    disabled={isSubmitting}
                    className="flex items-center"
                  >
                    {isSubmitting ? (
                      <>
                        <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                        Saving...
                      </>
                    ) : (
                      <>
                        <Save className="mr-2 h-4 w-4" />
                        Save Changes
                      </>
                    )}
                  </Button>
                </div>
              )}
            </form>
          </CardContent>
        </Card>

        {/* Account Information Card */}
        <Card className="mb-6">
          <CardHeader>
            <CardTitle className="flex items-center text-lg">
              <Shield className="mr-2 h-5 w-5" />
              Account Information
            </CardTitle>
            <CardDescription>
              Read-only account details and system information
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              {/* Email */}
              <div className="space-y-2">
                <Label className="flex items-center">
                  <Mail className="mr-2 h-4 w-4" />
                  Email Address
                </Label>
                <p className="text-gray-900 font-medium">{user.email}</p>
                <p className="text-xs text-gray-500">
                  Contact support to change your email
                </p>
              </div>

              {/* Role */}
              <div className="space-y-2">
                <Label className="flex items-center">
                  <Shield className="mr-2 h-4 w-4" />
                  Role
                </Label>
                <p className="text-gray-900 font-medium capitalize">
                  {user.role}
                </p>
                <p className="text-xs text-gray-500">
                  Your access level in the system
                </p>
              </div>

              {/* Account Status */}
              <div className="space-y-2">
                <Label>Account Status</Label>
                <div className="flex items-center">
                  <div
                    className={`w-2 h-2 rounded-full mr-2 ${user.active ? "bg-green-500" : "bg-red-500"}`}
                  ></div>
                  <p className="text-gray-900 font-medium">
                    {user.active ? "Active" : "Inactive"}
                  </p>
                </div>
                {!user.active && (
                  <p className="text-xs text-red-600">
                    Contact support to reactivate your account
                  </p>
                )}
              </div>

              {/* Member Since */}
              <div className="space-y-2">
                <Label className="flex items-center">
                  <Calendar className="mr-2 h-4 w-4" />
                  Member Since
                </Label>
                <p className="text-gray-900 font-medium">
                  {formatDate(user.created_at)}
                </p>
                <p className="text-xs text-gray-500">
                  When your account was created
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Security Section */}
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Security</CardTitle>
            <CardDescription>
              Manage your account security settings
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="flex items-center justify-between p-4 border rounded-lg">
                <div>
                  <h4 className="font-medium">Password</h4>
                  <p className="text-sm text-gray-600">
                    Update your password to keep your account secure
                  </p>
                </div>
                <Link href="/change-password">
                  <Button variant="outline" size="sm">
                    Change Password
                  </Button>
                </Link>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
