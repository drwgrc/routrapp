// Auth components exports
export { AuthGuard } from "./auth-guard";
export { ProtectedRoute } from "./protected-route";
export {
  RoleGuard,
  PermissionGuard,
  OwnerOnly,
  TechnicianOnly,
} from "./role-guard";
export {
  PageGuard,
  OwnerPage,
  TechnicianPage,
  AuthenticatedPage,
  HierarchicalTechnicianPage,
  TechnicianPageWithLayout,
  OwnerPageWithLayout,
  AuthenticatedPageWithLayout,
} from "./page-guard";
export {
  AuthMiddleware,
  RoleRedirectMiddleware,
  PermissionMiddleware,
  SessionTimeoutMiddleware,
  CombinedMiddleware,
} from "./route-middleware";
export { AuthErrorBoundary, useAuthErrorBoundary } from "./auth-error-boundary";
export { ClientOnly } from "./client-only";
export { ClientAuthWrapper } from "./client-wrapper";

// Advanced visibility components
export {
  RoleVisibility,
  PermissionVisibility,
  CombinedVisibility,
  ConditionalVisibility,
  FeatureFlag,
  AdminOnly,
  ManagementOnly,
  ReadOnly,
} from "./visibility-guard";

// Conditional rendering components with loading states
export {
  RoleConditionalRender,
  PermissionConditionalRender,
  CombinedConditionalRender,
  CustomConditionalRender,
  SuspenseConditionalRender,
  AdminConditionalRender,
  ManagementConditionalRender,
  ReadOnlyConditionalRender,
} from "./conditional-render";
