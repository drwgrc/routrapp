import React from "react";
import { render, screen } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { AuthProvider } from "@/contexts/auth-context";
import {
  RoleVisibility,
  PermissionVisibility,
  CombinedVisibility,
  ConditionalVisibility,
  FeatureFlag,
  AdminOnly,
  ManagementOnly,
  ReadOnly,
  RoleConditionalRender,
  PermissionConditionalRender,
  CombinedConditionalRender,
  CustomConditionalRender,
  AdminConditionalRender,
  ManagementConditionalRender,
  ReadOnlyConditionalRender,
} from "@/components/auth";
import { UserRole } from "@/types/auth";

// Mock auth service
jest.mock("@/services/auth-service", () => ({
  __esModule: true,
  default: {
    isAuthenticated: jest.fn(),
    getCurrentUser: jest.fn(),
    login: jest.fn(),
    logout: jest.fn(),
    register: jest.fn(),
  },
}));

// Test wrapper component
const TestWrapper: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });

  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>{children}</AuthProvider>
    </QueryClientProvider>
  );
};

// Mock user data
const mockOwnerUser = {
  id: 1,
  email: "owner@example.com",
  first_name: "John",
  last_name: "Owner",
  active: true,
  role: "owner" as UserRole,
  created_at: "2024-01-01T00:00:00Z",
  updated_at: "2024-01-01T00:00:00Z",
};

const mockTechnicianUser = {
  id: 2,
  email: "tech@example.com",
  first_name: "Jane",
  last_name: "Technician",
  active: true,
  role: "technician" as UserRole,
  created_at: "2024-01-01T00:00:00Z",
  updated_at: "2024-01-01T00:00:00Z",
};

describe("Role-Based Visibility System", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe("RoleVisibility", () => {
    it("should show content for owner when allowed", () => {
      render(
        <TestWrapper>
          <RoleVisibility allowedRoles={["owner"]}>
            <div>Owner content</div>
          </RoleVisibility>
        </TestWrapper>
      );

      expect(screen.getByText("Owner content")).toBeInTheDocument();
    });

    it("should hide content for technician when only owner allowed", () => {
      render(
        <TestWrapper>
          <RoleVisibility allowedRoles={["owner"]}>
            <div>Owner content</div>
          </RoleVisibility>
        </TestWrapper>
      );

      expect(screen.queryByText("Owner content")).not.toBeInTheDocument();
    });

    it("should show fallback when access denied", () => {
      render(
        <TestWrapper>
          <RoleVisibility
            allowedRoles={["owner"]}
            fallback={<div>Access denied</div>}
          >
            <div>Owner content</div>
          </RoleVisibility>
        </TestWrapper>
      );

      expect(screen.getByText("Access denied")).toBeInTheDocument();
      expect(screen.queryByText("Owner content")).not.toBeInTheDocument();
    });

    it("should support multiple roles", () => {
      render(
        <TestWrapper>
          <RoleVisibility allowedRoles={["owner", "technician"]}>
            <div>Multi-role content</div>
          </RoleVisibility>
        </TestWrapper>
      );

      expect(screen.getByText("Multi-role content")).toBeInTheDocument();
    });

    it("should support inverse logic", () => {
      render(
        <TestWrapper>
          <RoleVisibility allowedRoles={["owner"]} inverse={true}>
            <div>Non-owner content</div>
          </RoleVisibility>
        </TestWrapper>
      );

      expect(screen.queryByText("Non-owner content")).not.toBeInTheDocument();
    });
  });

  describe("PermissionVisibility", () => {
    it("should show content when user has required permission", () => {
      render(
        <TestWrapper>
          <PermissionVisibility requiredPermissions={["users.manage"]}>
            <div>User management content</div>
          </PermissionVisibility>
        </TestWrapper>
      );

      expect(screen.getByText("User management content")).toBeInTheDocument();
    });

    it("should hide content when user lacks permission", () => {
      render(
        <TestWrapper>
          <PermissionVisibility requiredPermissions={["system.*"]}>
            <div>System content</div>
          </PermissionVisibility>
        </TestWrapper>
      );

      expect(screen.queryByText("System content")).not.toBeInTheDocument();
    });

    it("should support multiple permissions with OR logic", () => {
      render(
        <TestWrapper>
          <PermissionVisibility
            requiredPermissions={["users.manage", "routes.manage"]}
            requireAllPermissions={false}
          >
            <div>Management content</div>
          </PermissionVisibility>
        </TestWrapper>
      );

      expect(screen.getByText("Management content")).toBeInTheDocument();
    });

    it("should support multiple permissions with AND logic", () => {
      render(
        <TestWrapper>
          <PermissionVisibility
            requiredPermissions={["users.read", "routes.read"]}
            requireAllPermissions={true}
          >
            <div>Read access content</div>
          </PermissionVisibility>
        </TestWrapper>
      );

      expect(screen.getByText("Read access content")).toBeInTheDocument();
    });
  });

  describe("CombinedVisibility", () => {
    it("should show content with AND logic when both conditions met", () => {
      render(
        <TestWrapper>
          <CombinedVisibility
            allowedRoles={["owner"]}
            requiredPermissions={["organizations.manage"]}
            logic="AND"
          >
            <div>Owner with org management</div>
          </CombinedVisibility>
        </TestWrapper>
      );

      expect(screen.getByText("Owner with org management")).toBeInTheDocument();
    });

    it("should show content with OR logic when either condition met", () => {
      render(
        <TestWrapper>
          <CombinedVisibility
            allowedRoles={["owner"]}
            requiredPermissions={["users.manage"]}
            logic="OR"
          >
            <div>Owner or user manager</div>
          </CombinedVisibility>
        </TestWrapper>
      );

      expect(screen.getByText("Owner or user manager")).toBeInTheDocument();
    });

    it("should hide content when neither condition met", () => {
      render(
        <TestWrapper>
          <CombinedVisibility
            allowedRoles={["technician"]}
            requiredPermissions={["system.*"]}
            logic="AND"
          >
            <div>Technician with system access</div>
          </CombinedVisibility>
        </TestWrapper>
      );

      expect(
        screen.queryByText("Technician with system access")
      ).not.toBeInTheDocument();
    });
  });

  describe("ConditionalVisibility", () => {
    it("should show content when condition is true", () => {
      render(
        <TestWrapper>
          <ConditionalVisibility condition={() => true}>
            <div>Conditional content</div>
          </ConditionalVisibility>
        </TestWrapper>
      );

      expect(screen.getByText("Conditional content")).toBeInTheDocument();
    });

    it("should hide content when condition is false", () => {
      render(
        <TestWrapper>
          <ConditionalVisibility condition={() => false}>
            <div>Conditional content</div>
          </ConditionalVisibility>
        </TestWrapper>
      );

      expect(screen.queryByText("Conditional content")).not.toBeInTheDocument();
    });

    it("should show fallback when condition is false", () => {
      render(
        <TestWrapper>
          <ConditionalVisibility
            condition={() => false}
            fallback={<div>Condition not met</div>}
          >
            <div>Conditional content</div>
          </ConditionalVisibility>
        </TestWrapper>
      );

      expect(screen.getByText("Condition not met")).toBeInTheDocument();
    });
  });

  describe("FeatureFlag", () => {
    it("should show content when feature is enabled for user role", () => {
      render(
        <TestWrapper>
          <FeatureFlag feature="advanced-analytics" enabledFor={["owner"]}>
            <div>Advanced analytics</div>
          </FeatureFlag>
        </TestWrapper>
      );

      expect(screen.getByText("Advanced analytics")).toBeInTheDocument();
    });

    it("should show content when feature is enabled for user permissions", () => {
      render(
        <TestWrapper>
          <FeatureFlag
            feature="user-management"
            enabledWithPermissions={["users.manage"]}
          >
            <div>User management feature</div>
          </FeatureFlag>
        </TestWrapper>
      );

      expect(screen.getByText("User management feature")).toBeInTheDocument();
    });

    it("should hide content when feature is not enabled", () => {
      render(
        <TestWrapper>
          <FeatureFlag feature="beta-feature" enabledFor={["technician"]}>
            <div>Beta feature</div>
          </FeatureFlag>
        </TestWrapper>
      );

      expect(screen.queryByText("Beta feature")).not.toBeInTheDocument();
    });
  });

  describe("Convenience Components", () => {
    describe("AdminOnly", () => {
      it("should show content for owner", () => {
        render(
          <TestWrapper>
            <AdminOnly>
              <div>Admin content</div>
            </AdminOnly>
          </TestWrapper>
        );

        expect(screen.getByText("Admin content")).toBeInTheDocument();
      });

      it("should hide content for technician", () => {
        render(
          <TestWrapper>
            <AdminOnly>
              <div>Admin content</div>
            </AdminOnly>
          </TestWrapper>
        );

        expect(screen.queryByText("Admin content")).not.toBeInTheDocument();
      });

      it("should show fallback when access denied", () => {
        render(
          <TestWrapper>
            <AdminOnly fallback={<div>Admin access required</div>}>
              <div>Admin content</div>
            </AdminOnly>
          </TestWrapper>
        );

        expect(screen.getByText("Admin access required")).toBeInTheDocument();
      });
    });

    describe("ManagementOnly", () => {
      it("should show content for users with management permissions", () => {
        render(
          <TestWrapper>
            <ManagementOnly>
              <div>Management content</div>
            </ManagementOnly>
          </TestWrapper>
        );

        expect(screen.getByText("Management content")).toBeInTheDocument();
      });

      it("should show fallback when access denied", () => {
        render(
          <TestWrapper>
            <ManagementOnly fallback={<div>Management access required</div>}>
              <div>Management content</div>
            </ManagementOnly>
          </TestWrapper>
        );

        expect(
          screen.getByText("Management access required")
        ).toBeInTheDocument();
      });
    });

    describe("ReadOnly", () => {
      it("should show content for users with read permissions", () => {
        render(
          <TestWrapper>
            <ReadOnly>
              <div>Read-only content</div>
            </ReadOnly>
          </TestWrapper>
        );

        expect(screen.getByText("Read-only content")).toBeInTheDocument();
      });

      it("should show fallback when access denied", () => {
        render(
          <TestWrapper>
            <ReadOnly fallback={<div>Read access required</div>}>
              <div>Read-only content</div>
            </ReadOnly>
          </TestWrapper>
        );

        expect(screen.getByText("Read access required")).toBeInTheDocument();
      });
    });
  });

  describe("Conditional Rendering Components", () => {
    describe("RoleConditionalRender", () => {
      it("should show content for allowed role", () => {
        render(
          <TestWrapper>
            <RoleConditionalRender allowedRoles={["owner"]}>
              <div>Owner content</div>
            </RoleConditionalRender>
          </TestWrapper>
        );

        expect(screen.getByText("Owner content")).toBeInTheDocument();
      });

      it("should show loading fallback while loading", () => {
        render(
          <TestWrapper>
            <RoleConditionalRender
              allowedRoles={["owner"]}
              loadingFallback={<div>Loading...</div>}
            >
              <div>Owner content</div>
            </RoleConditionalRender>
          </TestWrapper>
        );

        expect(screen.getByText("Loading...")).toBeInTheDocument();
      });
    });

    describe("PermissionConditionalRender", () => {
      it("should show content for allowed permission", () => {
        render(
          <TestWrapper>
            <PermissionConditionalRender requiredPermissions={["users.manage"]}>
              <div>User management content</div>
            </PermissionConditionalRender>
          </TestWrapper>
        );

        expect(screen.getByText("User management content")).toBeInTheDocument();
      });

      it("should show error fallback on error", () => {
        render(
          <TestWrapper>
            <PermissionConditionalRender
              requiredPermissions={["users.manage"]}
              errorFallback="Permission check failed"
            >
              <div>User management content</div>
            </PermissionConditionalRender>
          </TestWrapper>
        );

        expect(screen.getByText("Permission check failed")).toBeInTheDocument();
      });
    });

    describe("CombinedConditionalRender", () => {
      it("should show content when both conditions met with AND logic", () => {
        render(
          <TestWrapper>
            <CombinedConditionalRender
              allowedRoles={["owner"]}
              requiredPermissions={["organizations.manage"]}
              logic="AND"
            >
              <div>Owner with org management</div>
            </CombinedConditionalRender>
          </TestWrapper>
        );

        expect(
          screen.getByText("Owner with org management")
        ).toBeInTheDocument();
      });

      it("should show content when either condition met with OR logic", () => {
        render(
          <TestWrapper>
            <CombinedConditionalRender
              allowedRoles={["owner"]}
              requiredPermissions={["users.manage"]}
              logic="OR"
            >
              <div>Owner or user manager</div>
            </CombinedConditionalRender>
          </TestWrapper>
        );

        expect(screen.getByText("Owner or user manager")).toBeInTheDocument();
      });
    });

    describe("CustomConditionalRender", () => {
      it("should show content when custom condition is true", () => {
        render(
          <TestWrapper>
            <CustomConditionalRender condition={() => true}>
              <div>Custom content</div>
            </CustomConditionalRender>
          </TestWrapper>
        );

        expect(screen.getByText("Custom content")).toBeInTheDocument();
      });

      it("should hide content when custom condition is false", () => {
        render(
          <TestWrapper>
            <CustomConditionalRender condition={() => false}>
              <div>Custom content</div>
            </CustomConditionalRender>
          </TestWrapper>
        );

        expect(screen.queryByText("Custom content")).not.toBeInTheDocument();
      });
    });

    describe("AdminConditionalRender", () => {
      it("should show content for admin users", () => {
        render(
          <TestWrapper>
            <AdminConditionalRender>
              <div>Admin content</div>
            </AdminConditionalRender>
          </TestWrapper>
        );

        expect(screen.getByText("Admin content")).toBeInTheDocument();
      });

      it("should show fallback for non-admin users", () => {
        render(
          <TestWrapper>
            <AdminConditionalRender fallback={<div>Admin access required</div>}>
              <div>Admin content</div>
            </AdminConditionalRender>
          </TestWrapper>
        );

        expect(screen.getByText("Admin access required")).toBeInTheDocument();
      });
    });

    describe("ManagementConditionalRender", () => {
      it("should show content for management users", () => {
        render(
          <TestWrapper>
            <ManagementConditionalRender>
              <div>Management content</div>
            </ManagementConditionalRender>
          </TestWrapper>
        );

        expect(screen.getByText("Management content")).toBeInTheDocument();
      });

      it("should show fallback for non-management users", () => {
        render(
          <TestWrapper>
            <ManagementConditionalRender
              fallback={<div>Management access required</div>}
            >
              <div>Management content</div>
            </ManagementConditionalRender>
          </TestWrapper>
        );

        expect(
          screen.getByText("Management access required")
        ).toBeInTheDocument();
      });
    });

    describe("ReadOnlyConditionalRender", () => {
      it("should show content for read-only users", () => {
        render(
          <TestWrapper>
            <ReadOnlyConditionalRender>
              <div>Read-only content</div>
            </ReadOnlyConditionalRender>
          </TestWrapper>
        );

        expect(screen.getByText("Read-only content")).toBeInTheDocument();
      });

      it("should show fallback for users without read access", () => {
        render(
          <TestWrapper>
            <ReadOnlyConditionalRender
              fallback={<div>Read access required</div>}
            >
              <div>Read-only content</div>
            </ReadOnlyConditionalRender>
          </TestWrapper>
        );

        expect(screen.getByText("Read access required")).toBeInTheDocument();
      });
    });
  });

  describe("Edge Cases", () => {
    it("should handle empty roles array", () => {
      render(
        <TestWrapper>
          <RoleVisibility allowedRoles={[]}>
            <div>Content</div>
          </RoleVisibility>
        </TestWrapper>
      );

      expect(screen.queryByText("Content")).not.toBeInTheDocument();
    });

    it("should handle empty permissions array", () => {
      render(
        <TestWrapper>
          <PermissionVisibility requiredPermissions={[]}>
            <div>Content</div>
          </PermissionVisibility>
        </TestWrapper>
      );

      expect(screen.queryByText("Content")).not.toBeInTheDocument();
    });

    it("should handle null children", () => {
      render(
        <TestWrapper>
          <RoleVisibility allowedRoles={["owner"]}>{null}</RoleVisibility>
        </TestWrapper>
      );

      // Should not throw error
      expect(screen.queryByText("Owner content")).not.toBeInTheDocument();
    });

    it("should handle undefined children", () => {
      render(
        <TestWrapper>
          <RoleVisibility allowedRoles={["owner"]}>{undefined}</RoleVisibility>
        </TestWrapper>
      );

      // Should not throw error
      expect(screen.queryByText("Owner content")).not.toBeInTheDocument();
    });
  });
});
