/**
 * Token Storage Service
 *
 * Provides secure, flexible token storage with multiple storage strategies
 * and fallbacks for SSR compatibility and enhanced security.
 */

export interface TokenStorage {
  setToken: (key: string, value: string) => Promise<void>;
  getToken: (key: string) => Promise<string | null>;
  removeToken: (key: string) => Promise<void>;
  clear: () => Promise<void>;
  isAvailable: () => boolean;
}

/**
 * In-memory storage implementation
 * Used as fallback when other storage methods aren't available
 */
class MemoryStorage implements TokenStorage {
  private storage = new Map<string, string>();

  async setToken(key: string, value: string): Promise<void> {
    this.storage.set(key, value);
  }

  async getToken(key: string): Promise<string | null> {
    return this.storage.get(key) || null;
  }

  async removeToken(key: string): Promise<void> {
    this.storage.delete(key);
  }

  async clear(): Promise<void> {
    this.storage.clear();
  }

  isAvailable(): boolean {
    return true; // Memory storage is always available
  }
}

/**
 * localStorage implementation with error handling
 */
class LocalStorage implements TokenStorage {
  async setToken(key: string, value: string): Promise<void> {
    if (!this.isAvailable()) {
      throw new Error("localStorage is not available");
    }
    try {
      localStorage.setItem(key, value);
    } catch (error) {
      throw new Error(`Failed to store token: ${error}`);
    }
  }

  async getToken(key: string): Promise<string | null> {
    if (!this.isAvailable()) {
      return null;
    }
    try {
      return localStorage.getItem(key);
    } catch (error) {
      console.warn("Failed to retrieve token from localStorage:", error);
      return null;
    }
  }

  async removeToken(key: string): Promise<void> {
    if (!this.isAvailable()) {
      return;
    }
    try {
      localStorage.removeItem(key);
    } catch (error) {
      console.warn("Failed to remove token from localStorage:", error);
    }
  }

  async clear(): Promise<void> {
    if (!this.isAvailable()) {
      return;
    }
    try {
      const keysToRemove: string[] = [];
      for (let i = 0; i < localStorage.length; i++) {
        const key = localStorage.key(i);
        if (key && this.isTokenKey(key)) {
          keysToRemove.push(key);
        }
      }
      keysToRemove.forEach(key => localStorage.removeItem(key));
    } catch (error) {
      console.warn("Failed to clear tokens from localStorage:", error);
    }
  }

  private isTokenKey(key: string): boolean {
    // More specific matching for actual token keys
    const tokenPatterns = [
      /^access_?token$/i,
      /^refresh_?token$/i,
      /^id_?token$/i,
      /^auth_?token$/i,
      /^bearer_?token$/i,
      /^jwt_?token$/i,
      /^session_?token$/i,
      /^api_?token$/i,
      /^oauth_?token$/i,
      /^token$/i,
      /^auth$/i,
    ];

    return tokenPatterns.some(pattern => pattern.test(key));
  }

  isAvailable(): boolean {
    if (typeof window === "undefined") {
      return false;
    }
    try {
      const testKey = "__test_storage__";
      localStorage.setItem(testKey, "test");
      localStorage.removeItem(testKey);
      return true;
    } catch {
      return false;
    }
  }
}

/**
 * sessionStorage implementation with error handling
 */
class SessionStorage implements TokenStorage {
  async setToken(key: string, value: string): Promise<void> {
    if (!this.isAvailable()) {
      throw new Error("sessionStorage is not available");
    }
    try {
      sessionStorage.setItem(key, value);
    } catch (error) {
      throw new Error(`Failed to store token: ${error}`);
    }
  }

  async getToken(key: string): Promise<string | null> {
    if (!this.isAvailable()) {
      return null;
    }
    try {
      return sessionStorage.getItem(key);
    } catch (error) {
      console.warn("Failed to retrieve token from sessionStorage:", error);
      return null;
    }
  }

  async removeToken(key: string): Promise<void> {
    if (!this.isAvailable()) {
      return;
    }
    try {
      sessionStorage.removeItem(key);
    } catch (error) {
      console.warn("Failed to remove token from sessionStorage:", error);
    }
  }

  async clear(): Promise<void> {
    if (!this.isAvailable()) {
      return;
    }
    try {
      const keysToRemove: string[] = [];
      for (let i = 0; i < sessionStorage.length; i++) {
        const key = sessionStorage.key(i);
        if (key && this.isTokenKey(key)) {
          keysToRemove.push(key);
        }
      }
      keysToRemove.forEach(key => sessionStorage.removeItem(key));
    } catch (error) {
      console.warn("Failed to clear tokens from sessionStorage:", error);
    }
  }

  private isTokenKey(key: string): boolean {
    // More specific matching for actual token keys
    const tokenPatterns = [
      /^access_?token$/i,
      /^refresh_?token$/i,
      /^id_?token$/i,
      /^auth_?token$/i,
      /^bearer_?token$/i,
      /^jwt_?token$/i,
      /^session_?token$/i,
      /^api_?token$/i,
      /^oauth_?token$/i,
      /^token$/i,
      /^auth$/i,
    ];

    return tokenPatterns.some(pattern => pattern.test(key));
  }

  isAvailable(): boolean {
    if (typeof window === "undefined") {
      return false;
    }
    try {
      const testKey = "__test_session_storage__";
      sessionStorage.setItem(testKey, "test");
      sessionStorage.removeItem(testKey);
      return true;
    } catch {
      return false;
    }
  }
}

/**
 * Cookie-based storage for enhanced security
 * Useful for httpOnly cookies when implemented server-side
 */
class CookieStorage implements TokenStorage {
  private readonly secure: boolean;
  private readonly sameSite: "strict" | "lax" | "none";

  constructor(secure = true, sameSite: "strict" | "lax" | "none" = "strict") {
    this.secure = secure;
    this.sameSite = sameSite;
  }

  async setToken(key: string, value: string): Promise<void> {
    if (!this.isAvailable()) {
      throw new Error("Cookies are not available");
    }

    const cookieOptions = [
      `${key}=${value}`,
      "path=/",
      `samesite=${this.sameSite}`,
      ...(this.secure ? ["secure"] : []),
      // Note: httpOnly cookies cannot be set from client-side JavaScript
      // This would need to be implemented server-side for maximum security
    ];

    document.cookie = cookieOptions.join("; ");
  }

  async getToken(key: string): Promise<string | null> {
    if (!this.isAvailable()) {
      return null;
    }

    try {
      const cookies = document.cookie.split(";");
      for (const cookie of cookies) {
        const [cookieKey, cookieValue] = cookie.trim().split("=");
        if (cookieKey === key) {
          return decodeURIComponent(cookieValue);
        }
      }
      return null;
    } catch (error) {
      console.warn("Failed to retrieve token from cookies:", error);
      return null;
    }
  }

  async removeToken(key: string): Promise<void> {
    if (!this.isAvailable()) {
      return;
    }

    try {
      document.cookie = `${key}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`;
    } catch (error) {
      console.warn("Failed to remove token from cookies:", error);
    }
  }

  async clear(): Promise<void> {
    if (!this.isAvailable()) {
      return;
    }

    try {
      const cookies = document.cookie.split(";");
      for (const cookie of cookies) {
        const [key] = cookie.trim().split("=");
        if (this.isTokenKey(key)) {
          await this.removeToken(key);
        }
      }
    } catch (error) {
      console.warn("Failed to clear tokens from cookies:", error);
    }
  }

  private isTokenKey(key: string): boolean {
    // More specific matching for actual token keys
    const tokenPatterns = [
      /^access_?token$/i,
      /^refresh_?token$/i,
      /^id_?token$/i,
      /^auth_?token$/i,
      /^bearer_?token$/i,
      /^jwt_?token$/i,
      /^session_?token$/i,
      /^api_?token$/i,
      /^oauth_?token$/i,
      /^token$/i,
      /^auth$/i,
    ];

    return tokenPatterns.some(pattern => pattern.test(key));
  }

  isAvailable(): boolean {
    return (
      typeof document !== "undefined" && typeof document.cookie === "string"
    );
  }
}

/**
 * Multi-strategy storage with automatic fallbacks
 */
class MultiStrategyStorage implements TokenStorage {
  private primaryStorage: TokenStorage;
  private fallbackStorage: TokenStorage;

  constructor(primary: TokenStorage, fallback: TokenStorage) {
    this.primaryStorage = primary;
    this.fallbackStorage = fallback;
  }

  private async tryPrimary<T>(
    operation: (storage: TokenStorage) => Promise<T>,
    fallbackValue: T
  ): Promise<T> {
    if (this.primaryStorage.isAvailable()) {
      try {
        return await operation(this.primaryStorage);
      } catch (error) {
        console.warn("Primary storage failed, falling back:", error);
      }
    }

    try {
      return await operation(this.fallbackStorage);
    } catch (error) {
      console.warn("Fallback storage also failed:", error);
      return fallbackValue;
    }
  }

  async setToken(key: string, value: string): Promise<void> {
    await this.tryPrimary(
      async storage => await storage.setToken(key, value),
      undefined
    );
  }

  async getToken(key: string): Promise<string | null> {
    return await this.tryPrimary(
      async storage => await storage.getToken(key),
      null
    );
  }

  async removeToken(key: string): Promise<void> {
    // Try to remove from both storages to ensure cleanup
    const promises = [this.primaryStorage, this.fallbackStorage]
      .filter(storage => storage.isAvailable())
      .map(storage => storage.removeToken(key).catch(console.warn));

    await Promise.allSettled(promises);
  }

  async clear(): Promise<void> {
    // Try to clear from both storages
    const promises = [this.primaryStorage, this.fallbackStorage]
      .filter(storage => storage.isAvailable())
      .map(storage => storage.clear().catch(console.warn));

    await Promise.allSettled(promises);
  }

  isAvailable(): boolean {
    return (
      this.primaryStorage.isAvailable() || this.fallbackStorage.isAvailable()
    );
  }
}

// Storage strategy configuration
export type StorageStrategy =
  | "localStorage"
  | "sessionStorage"
  | "cookie"
  | "memory";

export interface StorageConfig {
  strategy: StorageStrategy;
  fallback?: StorageStrategy;
  cookieOptions?: {
    secure?: boolean;
    sameSite?: "strict" | "lax" | "none";
  };
}

/**
 * Factory function to create appropriate storage instance
 */
export function createTokenStorage(config: StorageConfig): TokenStorage {
  const createStorage = (strategy: StorageStrategy): TokenStorage => {
    switch (strategy) {
      case "localStorage":
        return new LocalStorage();
      case "sessionStorage":
        return new SessionStorage();
      case "cookie":
        return new CookieStorage(
          config.cookieOptions?.secure,
          config.cookieOptions?.sameSite
        );
      case "memory":
        return new MemoryStorage();
      default:
        throw new Error(`Unsupported storage strategy: ${strategy}`);
    }
  };

  const primaryStorage = createStorage(config.strategy);

  if (config.fallback) {
    const fallbackStorage = createStorage(config.fallback);
    return new MultiStrategyStorage(primaryStorage, fallbackStorage);
  }

  return primaryStorage;
}

// Default storage instance with localStorage primary and memory fallback
export const defaultTokenStorage = createTokenStorage({
  strategy: "localStorage",
  fallback: "memory",
});

// Secure storage instance with cookie primary and localStorage fallback
export const secureTokenStorage = createTokenStorage({
  strategy: "cookie",
  fallback: "localStorage",
  cookieOptions: {
    secure: process.env.NODE_ENV === "production",
    sameSite: "strict",
  },
});
