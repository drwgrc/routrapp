import authService from "./auth-service";
import { authMigrationUtils } from "./auth-service";

export { authService, authMigrationUtils };

const services = {
  auth: authService,
  authMigration: authMigrationUtils,
};

export default services;
