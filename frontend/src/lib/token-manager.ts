/**
 * Token Manager Service
 *
 * Comprehensive JWT token management with automatic refresh, expiry tracking,
 * request queuing, and robust error handling.
 */

import { TokenStorage, defaultTokenStorage } from "./token-storage";

// Token-related constants
export const TOKEN_KEYS = {
  ACCESS_TOKEN: "auth_token",
  REFRESH_TOKEN: "refresh_token",
  TOKEN_EXPIRY: "token_expiry",
} as const;

// JWT token structure (partial, for client-side parsing)
export interface JWTPayload {
  sub: string; // user_id
  email: string;
  role: string;
  organization_id: string;
  exp: number; // expiry timestamp
  iat: number; // issued at timestamp
  type: "access" | "refresh";
}

// Token refresh response structure
export interface TokenRefreshResponse {
  access_token: string;
  token_type: string;
  expires_in: number;
}

// Token manager configuration
export interface TokenManagerConfig {
  storage?: TokenStorage;
  refreshEndpoint?: string;
  refreshThresholdMs?: number; // How early to refresh before expiry
  maxRetries?: number;
  retryDelayMs?: number;
  onTokenRefreshed?: (tokens: { accessToken: string; expiresAt: Date }) => void;
  onRefreshFailed?: (error: Error) => void;
  onTokenExpired?: () => void;
}

// Error types for better error handling
export class TokenError extends Error {
  constructor(
    message: string,
    public readonly code: string,
    public readonly originalError?: unknown
  ) {
    super(message);
    this.name = "TokenError";
  }
}

export class TokenExpiredError extends TokenError {
  constructor(message = "Token has expired") {
    super(message, "TOKEN_EXPIRED");
  }
}

export class TokenRefreshError extends TokenError {
  constructor(message: string, originalError?: unknown) {
    super(message, "REFRESH_FAILED", originalError);
  }
}

/**
 * Utility functions for JWT parsing
 */
export class JWTUtils {
  /**
   * Decode JWT token payload without verification
   * Note: This is for client-side expiry checking only, not for security
   */
  static decodePayload(token: string): JWTPayload | null {
    try {
      const parts = token.split(".");
      if (parts.length !== 3) {
        return null;
      }

      const payload = parts[1];
      // Add padding if needed
      const paddedPayload =
        payload + "=".repeat((4 - (payload.length % 4)) % 4);
      const decoded = atob(paddedPayload);
      return JSON.parse(decoded) as JWTPayload;
    } catch (error) {
      console.warn("Failed to decode JWT payload:", error);
      return null;
    }
  }

  /**
   * Check if token is expired (with optional buffer)
   */
  static isExpired(token: string, bufferMs = 0): boolean {
    const payload = this.decodePayload(token);
    if (!payload?.exp) {
      return true; // Treat invalid tokens as expired
    }

    const expiryTime = payload.exp * 1000; // Convert to milliseconds
    const now = Date.now();
    return now + bufferMs >= expiryTime;
  }

  /**
   * Get expiry date from token
   */
  static getExpiryDate(token: string): Date | null {
    const payload = this.decodePayload(token);
    if (!payload?.exp) {
      return null;
    }
    return new Date(payload.exp * 1000);
  }

  /**
   * Get time until expiry in milliseconds
   */
  static getTimeUntilExpiry(token: string): number | null {
    const expiryDate = this.getExpiryDate(token);
    if (!expiryDate) {
      return null;
    }
    return expiryDate.getTime() - Date.now();
  }
}

/**
 * Request queue for handling concurrent requests during token refresh
 */
class RequestQueue {
  private queue: Array<{
    resolve: (token: string) => void;
    reject: (error: Error) => void;
  }> = [];
  private isRefreshing = false;

  async enqueue(): Promise<string> {
    return new Promise((resolve, reject) => {
      this.queue.push({ resolve, reject });
    });
  }

  setRefreshing(refreshing: boolean): void {
    this.isRefreshing = refreshing;
  }

  isCurrentlyRefreshing(): boolean {
    return this.isRefreshing;
  }

  resolveAll(token: string): void {
    const queue = [...this.queue];
    this.queue = [];
    this.isRefreshing = false;
    queue.forEach(({ resolve }) => resolve(token));
  }

  rejectAll(error: Error): void {
    const queue = [...this.queue];
    this.queue = [];
    this.isRefreshing = false;
    queue.forEach(({ reject }) => reject(error));
  }

  clear(): void {
    this.queue = [];
    this.isRefreshing = false;
  }
}

/**
 * Main Token Manager class
 */
export class TokenManager {
  private storage: TokenStorage;
  private refreshEndpoint: string;
  private refreshThresholdMs: number;
  private maxRetries: number;
  private retryDelayMs: number;
  private onTokenRefreshed?: (tokens: {
    accessToken: string;
    expiresAt: Date;
  }) => void;
  private onRefreshFailed?: (error: Error) => void;
  private onTokenExpired?: () => void;

  private refreshQueue = new RequestQueue();
  private refreshTimer?: NodeJS.Timeout;
  private isDestroyed = false;

  constructor(config: TokenManagerConfig = {}) {
    this.storage = config.storage || defaultTokenStorage;
    this.refreshEndpoint = config.refreshEndpoint || "/api/v1/auth/refresh";
    this.refreshThresholdMs = config.refreshThresholdMs || 5 * 60 * 1000; // 5 minutes
    this.maxRetries = config.maxRetries || 3;
    this.retryDelayMs = config.retryDelayMs || 1000;
    this.onTokenRefreshed = config.onTokenRefreshed;
    this.onRefreshFailed = config.onRefreshFailed;
    this.onTokenExpired = config.onTokenExpired;
  }

  /**
   * Initialize token manager and set up automatic refresh
   */
  async initialize(): Promise<void> {
    await this.scheduleNextRefresh();
  }

  /**
   * Store tokens securely
   */
  async setTokens(accessToken: string, refreshToken: string): Promise<void> {
    if (this.isDestroyed) return;

    try {
      await Promise.all([
        this.storage.setToken(TOKEN_KEYS.ACCESS_TOKEN, accessToken),
        this.storage.setToken(TOKEN_KEYS.REFRESH_TOKEN, refreshToken),
      ]);

      // Store expiry time for quick checks
      const expiryDate = JWTUtils.getExpiryDate(accessToken);
      if (expiryDate) {
        await this.storage.setToken(
          TOKEN_KEYS.TOKEN_EXPIRY,
          expiryDate.toISOString()
        );
      }

      // Schedule next refresh
      await this.scheduleNextRefresh();
    } catch (error) {
      throw new TokenError("Failed to store tokens", "STORAGE_ERROR", error);
    }
  }

  /**
   * Get access token (with automatic refresh if needed)
   */
  async getAccessToken(): Promise<string | null> {
    if (this.isDestroyed) return null;

    try {
      const accessToken = await this.storage.getToken(TOKEN_KEYS.ACCESS_TOKEN);
      if (!accessToken) {
        return null;
      }

      // Check if token needs refresh
      if (JWTUtils.isExpired(accessToken, this.refreshThresholdMs)) {
        try {
          const refreshedToken = await this.refreshTokenIfNeeded();
          return refreshedToken || accessToken;
        } catch (error) {
          // If refresh fails, return the original token if it's still technically valid
          if (!JWTUtils.isExpired(accessToken)) {
            return accessToken;
          }
          // Token is expired and refresh failed
          this.onTokenExpired?.();
          return null;
        }
      }

      return accessToken;
    } catch (error) {
      console.warn("Failed to get access token:", error);
      return null;
    }
  }

  /**
   * Get refresh token
   */
  async getRefreshToken(): Promise<string | null> {
    if (this.isDestroyed) return null;

    try {
      return await this.storage.getToken(TOKEN_KEYS.REFRESH_TOKEN);
    } catch (error) {
      console.warn("Failed to get refresh token:", error);
      return null;
    }
  }

  /**
   * Check if user is authenticated
   */
  async isAuthenticated(): Promise<boolean> {
    const accessToken = await this.storage.getToken(TOKEN_KEYS.ACCESS_TOKEN);
    const refreshToken = await this.storage.getToken(TOKEN_KEYS.REFRESH_TOKEN);

    // User is authenticated if they have either a valid access token or a refresh token
    if (!accessToken && !refreshToken) {
      return false;
    }

    // If we have an access token, check if it's valid or can be refreshed
    if (accessToken && !JWTUtils.isExpired(accessToken)) {
      return true;
    }

    // If access token is expired but we have a refresh token, we're still authenticated
    return !!refreshToken;
  }

  /**
   * Refresh token if needed (with queue management for concurrent requests)
   */
  async refreshTokenIfNeeded(): Promise<string | null> {
    if (this.isDestroyed) return null;

    // If already refreshing, wait in queue
    if (this.refreshQueue.isCurrentlyRefreshing()) {
      try {
        return await this.refreshQueue.enqueue();
      } catch (error) {
        throw new TokenRefreshError("Token refresh failed", error);
      }
    }

    // Start refresh process
    this.refreshQueue.setRefreshing(true);

    try {
      const newToken = await this.performTokenRefresh();
      this.refreshQueue.resolveAll(newToken);
      return newToken;
    } catch (error) {
      this.refreshQueue.rejectAll(error as Error);
      throw error;
    }
  }

  /**
   * Perform the actual token refresh
   */
  private async performTokenRefresh(): Promise<string> {
    const refreshToken = await this.getRefreshToken();
    if (!refreshToken) {
      throw new TokenRefreshError("No refresh token available");
    }

    let lastError: Error | null = null;

    for (let attempt = 1; attempt <= this.maxRetries; attempt++) {
      try {
        const response = await fetch(this.refreshEndpoint, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ refresh_token: refreshToken }),
        });

        if (!response.ok) {
          const errorData = await response.json().catch(() => ({}));
          throw new TokenRefreshError(
            errorData.error?.message ||
              `HTTP ${response.status}: ${response.statusText}`
          );
        }

        const data = await response.json();
        const newAccessToken = data.data?.access_token;

        if (!newAccessToken) {
          throw new TokenRefreshError(
            "Invalid refresh response: missing access token"
          );
        }

        // Store the new access token
        await this.storage.setToken(TOKEN_KEYS.ACCESS_TOKEN, newAccessToken);

        // Update expiry time
        const expiryDate = JWTUtils.getExpiryDate(newAccessToken);
        if (expiryDate) {
          await this.storage.setToken(
            TOKEN_KEYS.TOKEN_EXPIRY,
            expiryDate.toISOString()
          );
          this.onTokenRefreshed?.({
            accessToken: newAccessToken,
            expiresAt: expiryDate,
          });
        }

        // Schedule next refresh
        await this.scheduleNextRefresh();

        return newAccessToken;
      } catch (error) {
        lastError = error as Error;
        console.warn(`Token refresh attempt ${attempt} failed:`, error);

        // If this was the last attempt, don't wait
        if (attempt < this.maxRetries) {
          await this.delay(this.retryDelayMs * attempt); // Exponential backoff
        }
      }
    }

    // All attempts failed
    const refreshError = new TokenRefreshError(
      `Token refresh failed after ${this.maxRetries} attempts`,
      lastError
    );

    this.onRefreshFailed?.(refreshError);
    throw refreshError;
  }

  /**
   * Clear all tokens and reset state
   */
  async clearTokens(): Promise<void> {
    try {
      await Promise.all([
        this.storage.removeToken(TOKEN_KEYS.ACCESS_TOKEN),
        this.storage.removeToken(TOKEN_KEYS.REFRESH_TOKEN),
        this.storage.removeToken(TOKEN_KEYS.TOKEN_EXPIRY),
      ]);

      this.clearRefreshTimer();
      this.refreshQueue.clear();
    } catch (error) {
      console.warn("Failed to clear tokens:", error);
    }
  }

  /**
   * Schedule the next automatic token refresh
   */
  private async scheduleNextRefresh(): Promise<void> {
    if (this.isDestroyed) return;

    this.clearRefreshTimer();

    try {
      const accessToken = await this.storage.getToken(TOKEN_KEYS.ACCESS_TOKEN);
      if (!accessToken) {
        return;
      }

      const timeUntilExpiry = JWTUtils.getTimeUntilExpiry(accessToken);
      if (!timeUntilExpiry || timeUntilExpiry <= 0) {
        return;
      }

      // Schedule refresh before the threshold
      const refreshIn = Math.max(0, timeUntilExpiry - this.refreshThresholdMs);

      this.refreshTimer = setTimeout(async () => {
        try {
          await this.refreshTokenIfNeeded();
        } catch (error) {
          console.warn("Automatic token refresh failed:", error);
        }
      }, refreshIn);
    } catch (error) {
      console.warn("Failed to schedule token refresh:", error);
    }
  }

  /**
   * Clear the refresh timer
   */
  private clearRefreshTimer(): void {
    if (this.refreshTimer) {
      clearTimeout(this.refreshTimer);
      this.refreshTimer = undefined;
    }
  }

  /**
   * Utility delay function
   */
  private delay(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  /**
   * Destroy the token manager and clean up resources
   */
  destroy(): void {
    this.isDestroyed = true;
    this.clearRefreshTimer();
    this.refreshQueue.clear();
  }

  /**
   * Get token expiry information
   */
  async getTokenInfo(): Promise<{
    accessToken: string | null;
    isExpired: boolean;
    expiresAt: Date | null;
    timeUntilExpiry: number | null;
  }> {
    const accessToken = await this.storage.getToken(TOKEN_KEYS.ACCESS_TOKEN);

    if (!accessToken) {
      return {
        accessToken: null,
        isExpired: true,
        expiresAt: null,
        timeUntilExpiry: null,
      };
    }

    const isExpired = JWTUtils.isExpired(accessToken);
    const expiresAt = JWTUtils.getExpiryDate(accessToken);
    const timeUntilExpiry = JWTUtils.getTimeUntilExpiry(accessToken);

    return {
      accessToken,
      isExpired,
      expiresAt,
      timeUntilExpiry,
    };
  }
}

// Create default token manager instance
export const defaultTokenManager = new TokenManager();

// Initialize on module load (client-side only)
if (typeof window !== "undefined") {
  defaultTokenManager.initialize().catch(console.warn);
}
