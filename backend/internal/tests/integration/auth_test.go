package integration_test

import (
	"net/http"
	"testing"
	"time"

	"routrapp-api/internal/models"
	"routrapp-api/internal/tests"
)

func TestAuthHandler_Login(t *testing.T) {
	ctx, err := tests.SetupTestContext()
	if err != nil {
		t.Fatalf("Failed to setup test context: %v", err)
	}
	defer tests.CleanupTestContext(ctx)

	// Create test user
	testUser, err := tests.CreateCompleteTestUser(ctx.DB, "test@example.com", "password123", models.RoleTypeOwner, true)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	testCases := []struct {
		name           string
		email          string
		password       string
		expectedStatus int
		expectedCode   string
		checkSuccess   bool
	}{
		{
			name:           "Valid login",
			email:          "test@example.com",
			password:       "password123",
			expectedStatus: http.StatusOK,
			checkSuccess:   true,
		},
		{
			name:           "Invalid email format",
			email:          "invalid-email",
			password:       "password123",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "VALIDATION_ERROR",
		},
		{
			name:           "Empty password",
			email:          "test@example.com",
			password:       "",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "VALIDATION_ERROR",
		},
		{
			name:           "Non-existent user",
			email:          "nonexistent@example.com",
			password:       "password123",
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "INVALID_CREDENTIALS",
		},
		{
			name:           "Wrong password",
			email:          "test@example.com",
			password:       "wrongpassword",
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "INVALID_CREDENTIALS",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			w := tests.MakeLoginRequest(ctx.Router, tt.email, tt.password)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
				t.Logf("Response body: %s", w.Body.String())
			}

			if tt.checkSuccess {
				if !tests.AssertResponseSuccess(w, tt.expectedStatus) {
					t.Errorf("Expected successful response")
					t.Logf("Response body: %s", w.Body.String())
				}

				// Parse and validate login response
				loginResp, err := tests.ParseLoginResponse(w)
				if err != nil {
					t.Errorf("Failed to parse login response: %v", err)
					return
				}

				// Validate response data
				if loginResp.User.Email != testUser.User.Email {
					t.Errorf("Expected email %s, got %s", testUser.User.Email, loginResp.User.Email)
				}

				if loginResp.AccessToken == "" {
					t.Error("Access token should not be empty")
				}

				if loginResp.RefreshToken == "" {
					t.Error("Refresh token should not be empty")
				}

				if loginResp.TokenType != "Bearer" {
					t.Errorf("Expected token type 'Bearer', got %s", loginResp.TokenType)
				}

				if loginResp.ExpiresIn != 900 { // 15 minutes
					t.Errorf("Expected expires_in 900, got %d", loginResp.ExpiresIn)
				}

				// Verify tokens are valid
				accessClaims, err := ctx.JWTService.ValidateToken(loginResp.AccessToken)
				if err != nil {
					t.Errorf("Access token validation failed: %v", err)
				} else {
					if !accessClaims.IsAccessToken() {
						t.Error("Expected access token type")
					}
					if accessClaims.UserID != testUser.User.ID {
						t.Errorf("Expected user ID %d, got %d", testUser.User.ID, accessClaims.UserID)
					}
				}

				refreshClaims, err := ctx.JWTService.ValidateToken(loginResp.RefreshToken)
				if err != nil {
					t.Errorf("Refresh token validation failed: %v", err)
				} else {
					if !refreshClaims.IsRefreshToken() {
						t.Error("Expected refresh token type")
					}
				}
			} else if tt.expectedCode != "" {
				if !tests.AssertResponseError(w, tt.expectedStatus, tt.expectedCode) {
					t.Errorf("Expected error code %s", tt.expectedCode)
					t.Logf("Response body: %s", w.Body.String())
				}
			}
		})
	}
}

func TestAuthHandler_Login_InactiveUser(t *testing.T) {
	ctx, err := tests.SetupTestContext()
	if err != nil {
		t.Fatalf("Failed to setup test context: %v", err)
	}
	defer tests.CleanupTestContext(ctx)

	// Create inactive test user
	_, err = tests.CreateCompleteTestUser(ctx.DB, "inactive@example.com", "password123", models.RoleTypeOwner, false)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	w := tests.MakeLoginRequest(ctx.Router, "inactive@example.com", "password123")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	if !tests.AssertResponseError(w, http.StatusUnauthorized, "ACCOUNT_DISABLED") {
		t.Error("Expected ACCOUNT_DISABLED error")
		t.Logf("Response body: %s", w.Body.String())
	}
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	ctx, err := tests.SetupTestContext()
	if err != nil {
		t.Fatalf("Failed to setup test context: %v", err)
	}
	defer tests.CleanupTestContext(ctx)

	// Create test user and login to get tokens
	testUser, err := tests.CreateCompleteTestUser(ctx.DB, "test@example.com", "password123", models.RoleTypeOwner, true)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Login to get valid tokens
	loginW := tests.MakeLoginRequest(ctx.Router, "test@example.com", "password123")
	if loginW.Code != http.StatusOK {
		t.Fatalf("Login failed: %s", loginW.Body.String())
	}

	loginResp, err := tests.ParseLoginResponse(loginW)
	if err != nil {
		t.Fatalf("Failed to parse login response: %v", err)
	}

	testCases := []struct {
		name           string
		refreshToken   string
		expectedStatus int
		expectedCode   string
		checkSuccess   bool
	}{
		{
			name:           "Valid refresh token",
			refreshToken:   loginResp.RefreshToken,
			expectedStatus: http.StatusOK,
			checkSuccess:   true,
		},
		{
			name:           "Empty refresh token",
			refreshToken:   "",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "VALIDATION_ERROR",
		},
		{
			name:           "Invalid refresh token",
			refreshToken:   "invalid.token.here",
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "INVALID_REFRESH_TOKEN",
		},
		{
			name:           "Access token used instead of refresh token",
			refreshToken:   loginResp.AccessToken,
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "INVALID_TOKEN_TYPE",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			w := tests.MakeRefreshRequest(ctx.Router, tt.refreshToken)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
				t.Logf("Response body: %s", w.Body.String())
			}

			if tt.checkSuccess {
				if !tests.AssertResponseSuccess(w, tt.expectedStatus) {
					t.Errorf("Expected successful response")
					t.Logf("Response body: %s", w.Body.String())
				}

				// Parse and validate token response
				tokenResp, err := tests.ParseTokenResponse(w)
				if err != nil {
					t.Errorf("Failed to parse token response: %v", err)
					return
				}

				// Validate response data
				if tokenResp.AccessToken == "" {
					t.Error("Access token should not be empty")
				}

				if tokenResp.TokenType != "Bearer" {
					t.Errorf("Expected token type 'Bearer', got %s", tokenResp.TokenType)
				}

				if tokenResp.ExpiresIn != 900 { // 15 minutes
					t.Errorf("Expected expires_in 900, got %d", tokenResp.ExpiresIn)
				}

				// Verify new access token is valid
				accessClaims, err := ctx.JWTService.ValidateToken(tokenResp.AccessToken)
				if err != nil {
					t.Errorf("New access token validation failed: %v", err)
				} else {
					if !accessClaims.IsAccessToken() {
						t.Error("Expected access token type")
					}
					if accessClaims.UserID != testUser.User.ID {
						t.Errorf("Expected user ID %d, got %d", testUser.User.ID, accessClaims.UserID)
					}
				}
			} else if tt.expectedCode != "" {
				if !tests.AssertResponseError(w, tt.expectedStatus, tt.expectedCode) {
					t.Errorf("Expected error code %s", tt.expectedCode)
					t.Logf("Response body: %s", w.Body.String())
				}
			}
		})
	}
}

func TestAuthHandler_RefreshToken_InactiveUser(t *testing.T) {
	ctx, err := tests.SetupTestContext()
	if err != nil {
		t.Fatalf("Failed to setup test context: %v", err)
	}
	defer tests.CleanupTestContext(ctx)

	// Create test user and login
	testUser, err := tests.CreateCompleteTestUser(ctx.DB, "test@example.com", "password123", models.RoleTypeOwner, true)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Login to get tokens
	loginW := tests.MakeLoginRequest(ctx.Router, "test@example.com", "password123")
	loginResp, _ := tests.ParseLoginResponse(loginW)

	// Deactivate user
	testUser.User.Active = false
	ctx.DB.Save(testUser.User)

	// Try to refresh token with deactivated user
	w := tests.MakeRefreshRequest(ctx.Router, loginResp.RefreshToken)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	if !tests.AssertResponseError(w, http.StatusUnauthorized, "ACCOUNT_DISABLED") {
		t.Error("Expected ACCOUNT_DISABLED error")
		t.Logf("Response body: %s", w.Body.String())
	}
}

func TestAuthHandler_RefreshToken_TokenNotInDatabase(t *testing.T) {
	ctx, err := tests.SetupTestContext()
	if err != nil {
		t.Fatalf("Failed to setup test context: %v", err)
	}
	defer tests.CleanupTestContext(ctx)

	// Create test user and generate a refresh token manually
	testUser, err := tests.CreateCompleteTestUser(ctx.DB, "test@example.com", "password123", models.RoleTypeOwner, true)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Generate a refresh token but don't store it in the database
	refreshToken, err := ctx.JWTService.GenerateRefreshToken(
		testUser.User.ID,
		testUser.User.OrganizationID,
		testUser.User.Email,
		testUser.Role.Name.String(),
	)
	if err != nil {
		t.Fatalf("Failed to generate refresh token: %v", err)
	}

	// Try to use the refresh token that's not stored in the database
	w := tests.MakeRefreshRequest(ctx.Router, refreshToken)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	if !tests.AssertResponseError(w, http.StatusUnauthorized, "INVALID_REFRESH_TOKEN") {
		t.Error("Expected INVALID_REFRESH_TOKEN error")
		t.Logf("Response body: %s", w.Body.String())
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	ctx, err := tests.SetupTestContext()
	if err != nil {
		t.Fatalf("Failed to setup test context: %v", err)
	}
	defer tests.CleanupTestContext(ctx)

	// Create test user and login
	_, err = tests.CreateCompleteTestUser(ctx.DB, "test@example.com", "password123", models.RoleTypeOwner, true)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Login to get tokens
	loginW := tests.MakeLoginRequest(ctx.Router, "test@example.com", "password123")
	loginResp, _ := tests.ParseLoginResponse(loginW)

	testCases := []struct {
		name           string
		accessToken    string
		expectedStatus int
		expectedCode   string
		checkSuccess   bool
	}{
		{
			name:           "Valid logout",
			accessToken:    loginResp.AccessToken,
			expectedStatus: http.StatusOK,
			checkSuccess:   true,
		},
		{
			name:           "Missing authorization header",
			accessToken:    "",
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "MISSING_AUTH_HEADER",
		},
		{
			name:           "Invalid token format",
			accessToken:    "invalid.token.here",
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "INVALID_TOKEN",
		},
		{
			name:           "Refresh token used instead of access token",
			accessToken:    loginResp.RefreshToken,
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "INVALID_TOKEN_TYPE",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			w := tests.MakeLogoutRequest(ctx.Router, tt.accessToken)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
				t.Logf("Response body: %s", w.Body.String())
			}

			if tt.checkSuccess {
				if !tests.AssertResponseSuccess(w, tt.expectedStatus) {
					t.Errorf("Expected successful response")
					t.Logf("Response body: %s", w.Body.String())
				}

				// Verify refresh token was cleared from database
				var user models.User
				ctx.DB.First(&user, "email = ?", "test@example.com")
				if user.RefreshToken != "" {
					t.Error("Refresh token should be cleared from database")
				}
			} else if tt.expectedCode != "" {
				if !tests.AssertResponseError(w, tt.expectedStatus, tt.expectedCode) {
					t.Errorf("Expected error code %s", tt.expectedCode)
					t.Logf("Response body: %s", w.Body.String())
				}
			}
		})
	}
}

func TestAuthHandler_FullAuthFlow(t *testing.T) {
	ctx, err := tests.SetupTestContext()
	if err != nil {
		t.Fatalf("Failed to setup test context: %v", err)
	}
	defer tests.CleanupTestContext(ctx)

	// Create test user
	_, err = tests.CreateCompleteTestUser(ctx.DB, "test@example.com", "password123", models.RoleTypeOwner, true)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	t.Run("Complete authentication flow", func(t *testing.T) {
		// Step 1: Login
		loginW := tests.MakeLoginRequest(ctx.Router, "test@example.com", "password123")
		if loginW.Code != http.StatusOK {
			t.Fatalf("Login failed: %s", loginW.Body.String())
		}

		loginResp, err := tests.ParseLoginResponse(loginW)
		if err != nil {
			t.Fatalf("Failed to parse login response: %v", err)
		}

		t.Logf("Login successful - Access token: %s...", loginResp.AccessToken[:20])

		// Step 2: Use refresh token to get new access token
		refreshW := tests.MakeRefreshRequest(ctx.Router, loginResp.RefreshToken)
		if refreshW.Code != http.StatusOK {
			t.Fatalf("Token refresh failed: %s", refreshW.Body.String())
		}

		tokenResp, err := tests.ParseTokenResponse(refreshW)
		if err != nil {
			t.Fatalf("Failed to parse token response: %v", err)
		}

		t.Logf("Token refresh successful - New access token: %s...", tokenResp.AccessToken[:20])

		// Step 3: Logout using new access token
		logoutW := tests.MakeLogoutRequest(ctx.Router, tokenResp.AccessToken)
		if logoutW.Code != http.StatusOK {
			t.Fatalf("Logout failed: %s", logoutW.Body.String())
		}

		t.Log("Logout successful")

		// Step 4: Verify refresh token is cleared and cannot be used again
		finalRefreshW := tests.MakeRefreshRequest(ctx.Router, loginResp.RefreshToken)
		if finalRefreshW.Code != http.StatusUnauthorized {
			t.Errorf("Expected refresh to fail after logout, got status %d", finalRefreshW.Code)
		}

		t.Logf("Final refresh response body: %s", finalRefreshW.Body.String())
		
		if !tests.AssertResponseError(finalRefreshW, http.StatusUnauthorized, "INVALID_REFRESH_TOKEN") {
			t.Error("Expected INVALID_REFRESH_TOKEN error after logout")
			t.Logf("Response body: %s", finalRefreshW.Body.String())
		}

		t.Log("Refresh token properly invalidated after logout")
	})
}

func TestAuthHandler_LoginUpdatesLastLoginTime(t *testing.T) {
	ctx, err := tests.SetupTestContext()
	if err != nil {
		t.Fatalf("Failed to setup test context: %v", err)
	}
	defer tests.CleanupTestContext(ctx)

	// Create test user
	testUser, err := tests.CreateCompleteTestUser(ctx.DB, "test@example.com", "password123", models.RoleTypeOwner, true)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Verify last login is initially nil
	if testUser.User.LastLoginAt != nil {
		t.Error("LastLoginAt should be nil initially")
	}

	// Login
	beforeLogin := time.Now()
	loginW := tests.MakeLoginRequest(ctx.Router, "test@example.com", "password123")
	afterLogin := time.Now()

	if loginW.Code != http.StatusOK {
		t.Fatalf("Login failed: %s", loginW.Body.String())
	}

	// Check that last login time was updated
	var updatedUser models.User
	ctx.DB.First(&updatedUser, testUser.User.ID)

	if updatedUser.LastLoginAt == nil {
		t.Error("LastLoginAt should be set after login")
	} else {
		loginTime := *updatedUser.LastLoginAt
		if loginTime.Before(beforeLogin) || loginTime.After(afterLogin) {
			t.Errorf("LastLoginAt %v should be between %v and %v", loginTime, beforeLogin, afterLogin)
		}
	}
}

func TestAuthHandler_ConcurrentLogins(t *testing.T) {
	// Create separate test contexts for each goroutine to avoid database conflicts
	done := make(chan bool, 3)
	results := make([]bool, 3)

	for i := 0; i < 3; i++ {
		go func(index int) {
			ctx, err := tests.SetupTestContext()
			if err != nil {
				t.Errorf("Failed to setup test context for goroutine %d: %v", index, err)
				done <- true
				return
			}
			defer tests.CleanupTestContext(ctx)

			// Create test user
			_, err = tests.CreateCompleteTestUser(ctx.DB, "test@example.com", "password123", models.RoleTypeOwner, true)
			if err != nil {
				t.Errorf("Failed to create test user for goroutine %d: %v", index, err)
				done <- true
				return
			}

			w := tests.MakeLoginRequest(ctx.Router, "test@example.com", "password123")
			results[index] = w.Code == http.StatusOK
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 3; i++ {
		<-done
	}

	// All logins should succeed
	for i, success := range results {
		if !success {
			t.Errorf("Concurrent login %d failed", i)
		}
	}
} 