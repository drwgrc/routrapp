package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"routrapp-api/internal/models"
	"routrapp-api/internal/tests"
	"routrapp-api/internal/validation"
)

func TestAuthHandler_Register(t *testing.T) {
	ctx, err := tests.SetupTestContext()
	if err != nil {
		t.Fatalf("Failed to setup test context: %v", err)
	}
	defer tests.CleanupTestContext(ctx)

	// Create test organization and role for registration
	org, err := tests.CreateTestOrganization(ctx.DB)
	if err != nil {
		t.Fatalf("Failed to create test organization: %v", err)
	}

	_, err = tests.CreateTestRole(ctx.DB, org.ID, models.RoleTypeOwner)
	if err != nil {
		t.Fatalf("Failed to create test role: %v", err)
	}

	testCases := []struct {
		name           string
		request        validation.UserRegistrationRequest
		expectedStatus int
		expectedCode   string
		checkSuccess   bool
	}{
		{
			name: "Valid registration with strong password",
			request: validation.UserRegistrationRequest{
				Email:     "newuser@example.com",
				Password:  "StrongPass123!",
				FirstName: "John",
				LastName:  "Doe",
				Role:      models.RoleTypeOwner,
				TenantID:  org.ID,
			},
			expectedStatus: http.StatusCreated,
			checkSuccess:   true,
		},
		{
			name: "Registration with weak password (no uppercase)",
			request: validation.UserRegistrationRequest{
				Email:     "newuser2@example.com",
				Password:  "weakpass123!",
				FirstName: "Jane",
				LastName:  "Doe",
				Role:      models.RoleTypeOwner,
				TenantID:  org.ID,
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "WEAK_PASSWORD",
		},
		{
			name: "Registration with weak password (no special chars)",
			request: validation.UserRegistrationRequest{
				Email:     "newuser3@example.com",
				Password:  "WeakPass123",
				FirstName: "Bob",
				LastName:  "Smith",
				Role:      models.RoleTypeOwner,
				TenantID:  org.ID,
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "WEAK_PASSWORD",
		},
		{
			name: "Registration with common password",
			request: validation.UserRegistrationRequest{
				Email:     "newuser4@example.com",
				Password:  "TestPass123!", // This is in the common list and meets all validation requirements
				FirstName: "Alice",
				LastName:  "Johnson",
				Role:      models.RoleTypeOwner,
				TenantID:  org.ID,
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "COMMON_PASSWORD",
		},
		{
			name: "Registration with invalid email",
			request: validation.UserRegistrationRequest{
				Email:     "invalid-email",
				Password:  "StrongPass123!",
				FirstName: "Test",
				LastName:  "User",
				Role:      models.RoleTypeOwner,
				TenantID:  org.ID,
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "VALIDATION_ERROR",
		},
		{
			name: "Registration with non-existent organization",
			request: validation.UserRegistrationRequest{
				Email:     "user@example.com",
				Password:  "StrongPass123!",
				FirstName: "Test",
				LastName:  "User",
				Role:      models.RoleTypeOwner,
				TenantID:  99999,
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "INVALID_ORGANIZATION",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tc.request)
			req := httptest.NewRequest("POST", "/api/v1/auth/register-user", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			ctx.Router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d. Response: %s", tc.expectedStatus, w.Code, w.Body.String())
			}

			if tc.checkSuccess {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}

				if success, ok := response["success"]; !ok || success != true {
					t.Errorf("Expected success to be true, got %v", success)
				}

				if data, ok := response["data"]; !ok {
					t.Errorf("Expected data field in response")
				} else if userResp, ok := data.(map[string]interface{}); ok {
					if email := userResp["email"]; email != tc.request.Email {
						t.Errorf("Expected email %s, got %v", tc.request.Email, email)
					}
				}
			} else if tc.expectedCode != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("Failed to parse error response: %v", err)
				}

				if errorField, ok := response["error"]; ok {
					if errorObj, ok := errorField.(map[string]interface{}); ok {
						if details, ok := errorObj["details"]; ok {
							if detailsObj, ok := details.(map[string]interface{}); ok {
								if code := detailsObj["code"]; code != tc.expectedCode {
									t.Errorf("Expected error code %s, got %v", tc.expectedCode, code)
								}
							} else {
								t.Errorf("Expected error details object, got %v", details)
							}
						} else {
							t.Errorf("Expected error details field in response")
						}
					} else {
						t.Errorf("Expected error object, got %v", errorField)
					}
				} else {
					t.Errorf("Expected error field in response")
				}
			}
		})
	}
}

func TestAuthHandler_ChangePassword(t *testing.T) {
	ctx, err := tests.SetupTestContext()
	if err != nil {
		t.Fatalf("Failed to setup test context: %v", err)
	}
	defer tests.CleanupTestContext(ctx)

	// Create test user
	testUser, err := tests.CreateCompleteTestUser(ctx.DB, "test@example.com", "CurrentPass123!", models.RoleTypeOwner, true)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Generate access token for authentication
	accessToken, err := ctx.JWTService.GenerateAccessToken(testUser.User.ID, testUser.User.OrganizationID, testUser.User.Email, testUser.User.Role.Name.String())
	if err != nil {
		t.Fatalf("Failed to generate access token: %v", err)
	}

	testCases := []struct {
		name           string
		request        validation.ChangePasswordRequest
		accessToken    string
		expectedStatus int
		expectedCode   string
		checkSuccess   bool
	}{
		{
			name: "Valid password change",
			request: validation.ChangePasswordRequest{
				CurrentPassword: "CurrentPass123!",
				NewPassword:     "NewStrongPass456@",
			},
			accessToken:    accessToken,
			expectedStatus: http.StatusOK,
			checkSuccess:   true,
		},
		{
			name: "Password change with weak new password",
			request: validation.ChangePasswordRequest{
				CurrentPassword: "CurrentPass123!",
				NewPassword:     "weakpass",
			},
			accessToken:    accessToken,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "WEAK_PASSWORD",
		},
		{
			name: "Password change with common new password",
			request: validation.ChangePasswordRequest{
				CurrentPassword: "CurrentPass123!",
				NewPassword:     "TestPass123!", // This password meets all validation requirements but is in the common list
			},
			accessToken:    accessToken,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "COMMON_PASSWORD",
		},
		{
			name: "Password change with same password",
			request: validation.ChangePasswordRequest{
				CurrentPassword: "CurrentPass123!",
				NewPassword:     "CurrentPass123!",
			},
			accessToken:    accessToken,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "SAME_PASSWORD",
		},
		{
			name: "Password change with wrong current password",
			request: validation.ChangePasswordRequest{
				CurrentPassword: "WrongPassword123!",
				NewPassword:     "NewStrongPass456@",
			},
			accessToken:    accessToken,
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "INVALID_CURRENT_PASSWORD",
		},
		{
			name: "Password change without authentication",
			request: validation.ChangePasswordRequest{
				CurrentPassword: "CurrentPass123!",
				NewPassword:     "NewStrongPass456@",
			},
			accessToken:    "",
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   "MISSING_AUTH_HEADER",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tc.request)
			req := httptest.NewRequest("POST", "/api/v1/auth/change-password", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			if tc.accessToken != "" {
				req.Header.Set("Authorization", "Bearer "+tc.accessToken)
			}

			w := httptest.NewRecorder()
			ctx.Router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d. Response: %s", tc.expectedStatus, w.Code, w.Body.String())
			}

			if tc.checkSuccess {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}

				if success, ok := response["success"]; !ok || success != true {
					t.Errorf("Expected success to be true, got %v", success)
				}
			} else if tc.expectedCode != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("Failed to parse error response: %v", err)
				}

				if errorField, ok := response["error"]; ok {
					if errorObj, ok := errorField.(map[string]interface{}); ok {
						if details, ok := errorObj["details"]; ok {
							if detailsObj, ok := details.(map[string]interface{}); ok {
								if code := detailsObj["code"]; code != tc.expectedCode {
									t.Errorf("Expected error code %s, got %v", tc.expectedCode, code)
								}
							} else {
								t.Errorf("Expected error details object, got %v", details)
							}
						} else {
							t.Errorf("Expected error details field in response")
						}
					} else {
						t.Errorf("Expected error object, got %v", errorField)
					}
				} else {
					t.Errorf("Expected error field in response")
				}
			}
		})
	}
}

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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := tests.MakeLoginRequest(ctx.Router, tc.email, tc.password)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d. Response: %s", tc.expectedStatus, w.Code, w.Body.String())
			}

			if tc.checkSuccess {
				loginResp, err := tests.ParseLoginResponse(w)
				if err != nil {
					t.Fatalf("Failed to parse login response: %v", err)
				}

				if loginResp.User.Email != tc.email {
					t.Errorf("Expected email %s, got %s", tc.email, loginResp.User.Email)
				}

				if loginResp.AccessToken == "" {
					t.Error("Expected access token to be present")
				}

				if loginResp.RefreshToken == "" {
					t.Error("Expected refresh token to be present")
				}
			} else if tc.expectedCode != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("Failed to parse error response: %v", err)
				}

				if errorField, ok := response["error"]; ok {
					if errorObj, ok := errorField.(map[string]interface{}); ok {
						if details, ok := errorObj["details"]; ok {
							if detailsObj, ok := details.(map[string]interface{}); ok {
								if code := detailsObj["code"]; code != tc.expectedCode {
									t.Errorf("Expected error code %s, got %v", tc.expectedCode, code)
								}
							} else {
								t.Errorf("Expected error details object, got %v", details)
							}
						} else {
							t.Errorf("Expected error details field in response")
						}
					} else {
						t.Errorf("Expected error object, got %v", errorField)
					}
				} else {
					t.Errorf("Expected error field in response")
				}
			}
		})
	}

	t.Run("Inactive user login", func(t *testing.T) {
		// Create technician role in the same organization
		techRole, err := tests.CreateTestRole(ctx.DB, testUser.Organization.ID, models.RoleTypeTechnician)
		if err != nil {
			t.Fatalf("Failed to create technician role: %v", err)
		}

		// Create inactive user in the same organization
		_, err = tests.CreateTestUser(ctx.DB, testUser.Organization.ID, techRole.ID, "inactive@example.com", "password123", false)
		if err != nil {
			t.Fatalf("Failed to create inactive user: %v", err)
		}

		w := tests.MakeLoginRequest(ctx.Router, "inactive@example.com", "password123")

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to parse error response: %v", err)
		}

		if errorField, ok := response["error"]; ok {
			if errorObj, ok := errorField.(map[string]interface{}); ok {
				if details, ok := errorObj["details"]; ok {
					if detailsObj, ok := details.(map[string]interface{}); ok {
						if code := detailsObj["code"]; code != "ACCOUNT_DISABLED" {
							t.Errorf("Expected error code ACCOUNT_DISABLED, got %v", code)
						}
					} else {
						t.Errorf("Expected error details object, got %v", details)
					}
				} else {
					t.Errorf("Expected error details field in response")
				}
			} else {
				t.Errorf("Expected error object, got %v", errorField)
			}
		} else {
			t.Errorf("Expected error field in response")
		}
	})
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

func TestAuthHandler_GetCurrentUser(t *testing.T) {
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

	// Login to get access token
	loginW := tests.MakeLoginRequest(ctx.Router, "test@example.com", "password123")
	if loginW.Code != http.StatusOK {
		t.Fatalf("Login failed: %s", loginW.Body.String())
	}

	var loginResp map[string]interface{}
	if err := json.Unmarshal(loginW.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("Failed to parse login response: %v", err)
	}

	accessToken := loginResp["data"].(map[string]interface{})["access_token"].(string)

	// Test GetCurrentUser with valid token
	t.Run("Valid token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/auth/me", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		w := httptest.NewRecorder()
		ctx.Router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if !response["success"].(bool) {
			t.Error("Expected success to be true")
		}

		userData := response["data"].(map[string]interface{})
		if userData["email"] != "test@example.com" {
			t.Errorf("Expected email 'test@example.com', got '%s'", userData["email"])
		}
		if userData["first_name"] != testUser.User.FirstName {
			t.Errorf("Expected first_name '%s', got '%s'", testUser.User.FirstName, userData["first_name"])
		}
		if userData["last_name"] != testUser.User.LastName {
			t.Errorf("Expected last_name '%s', got '%s'", testUser.User.LastName, userData["last_name"])
		}
		if userData["role"] != "owner" {
			t.Errorf("Expected role 'owner', got '%s'", userData["role"])
		}
	})

	// Test GetCurrentUser without token
	t.Run("No token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/auth/me", nil)
		w := httptest.NewRecorder()
		ctx.Router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d: %s", w.Code, w.Body.String())
		}
	})

	// Test GetCurrentUser with invalid token
	t.Run("Invalid token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/auth/me", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		w := httptest.NewRecorder()
		ctx.Router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d: %s", w.Code, w.Body.String())
		}
	})
} 