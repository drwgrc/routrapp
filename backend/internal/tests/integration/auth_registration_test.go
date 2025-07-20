package integration

import (
	"net/http"
	"testing"

	"routrapp-api/internal/models"
	"routrapp-api/internal/tests"
)

func TestAuthHandler_Register_Success(t *testing.T) {
	ctx, err := tests.SetupTestContext()
	if err != nil {
		t.Fatalf("Failed to setup test context: %v", err)
	}
	defer tests.CleanupTestContext(ctx)

	t.Run("Successful registration creates organization and user", func(t *testing.T) {
		w := tests.MakeRegistrationRequest(
			ctx.Router,
			"john@example.com",
			"password123",
			"John",
			"Doe",
			"Acme Corp",
			"contact@acme.com",
			"acme",
		)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusCreated, w.Code, w.Body.String())
			return
		}

		resp, err := tests.ParseRegistrationResponse(w)
		if err != nil {
			t.Fatalf("Failed to parse registration response: %v", err)
		}

		// Verify user data
		if resp.User.Email != "john@example.com" {
			t.Errorf("Expected user email 'john@example.com', got '%s'", resp.User.Email)
		}
		if resp.User.FirstName != "John" {
			t.Errorf("Expected user first name 'John', got '%s'", resp.User.FirstName)
		}
		if resp.User.LastName != "Doe" {
			t.Errorf("Expected user last name 'Doe', got '%s'", resp.User.LastName)
		}
		if resp.User.Role != "owner" {
			t.Errorf("Expected user role 'owner', got '%s'", resp.User.Role)
		}
		if !resp.User.Active {
			t.Error("Expected user to be active")
		}

		// Verify organization data
		if resp.Organization.Name != "Acme Corp" {
			t.Errorf("Expected organization name 'Acme Corp', got '%s'", resp.Organization.Name)
		}
		if resp.Organization.SubDomain != "acme" {
			t.Errorf("Expected organization subdomain 'acme', got '%s'", resp.Organization.SubDomain)
		}
		if resp.Organization.ContactEmail != "contact@acme.com" {
			t.Errorf("Expected organization email 'contact@acme.com', got '%s'", resp.Organization.ContactEmail)
		}
		if !resp.Organization.Active {
			t.Error("Expected organization to be active")
		}
		if resp.Organization.PlanType != "basic" {
			t.Errorf("Expected organization plan type 'basic', got '%s'", resp.Organization.PlanType)
		}

		// Verify tokens
		if resp.AccessToken == "" {
			t.Error("Expected access token to be provided")
		}
		if resp.RefreshToken == "" {
			t.Error("Expected refresh token to be provided")
		}
		if resp.TokenType != "Bearer" {
			t.Errorf("Expected token type 'Bearer', got '%s'", resp.TokenType)
		}
		if resp.ExpiresIn != 900 { // 15 minutes
			t.Errorf("Expected expires in 900 seconds, got %d", resp.ExpiresIn)
		}

		// Verify database state - organization was created
		var org models.Organization
		if err := ctx.DB.Where("sub_domain = ?", "acme").First(&org).Error; err != nil {
			t.Errorf("Organization was not created in database: %v", err)
		}

		// Verify database state - roles were created
		var roles []models.Role
		if err := ctx.DB.Where("organization_id = ?", org.ID).Find(&roles).Error; err != nil {
			t.Errorf("Failed to fetch roles: %v", err)
		}
		if len(roles) != 2 {
			t.Errorf("Expected 2 roles to be created, got %d", len(roles))
		}

		// Verify owner and technician roles exist
		var ownerRole, techRole models.Role
		for _, role := range roles {
			if role.Name == models.RoleTypeOwner {
				ownerRole = role
			} else if role.Name == models.RoleTypeTechnician {
				techRole = role
			}
		}

		if ownerRole.ID == 0 {
			t.Error("Owner role was not created")
		}
		if techRole.ID == 0 {
			t.Error("Technician role was not created")
		}

		// Verify user was created with correct role
		var user models.User
		if err := ctx.DB.Preload("Role").Where("email = ? AND organization_id = ?", "john@example.com", org.ID).First(&user).Error; err != nil {
			t.Errorf("User was not created in database: %v", err)
		}
		if user.RoleID != ownerRole.ID {
			t.Errorf("User was not assigned owner role. Expected role ID %d, got %d", ownerRole.ID, user.RoleID)
		}
		if user.RefreshToken == "" {
			t.Error("User refresh token was not set")
		}
		if user.LastLoginAt == nil {
			t.Error("User last login time was not set")
		}
	})
}

func TestAuthHandler_Register_ValidationErrors(t *testing.T) {
	ctx, err := tests.SetupTestContext()
	if err != nil {
		t.Fatalf("Failed to setup test context: %v", err)
	}
	defer tests.CleanupTestContext(ctx)

	testCases := []struct {
		name     string
		email    string
		password string
		expected int
	}{
		{
			name:     "Invalid email format",
			email:    "invalid-email",
			password: "password123",
			expected: http.StatusBadRequest,
		},
		{
			name:     "Empty email",
			email:    "",
			password: "password123",
			expected: http.StatusBadRequest,
		},
		{
			name:     "Password too short",
			email:    "test@example.com",
			password: "123",
			expected: http.StatusBadRequest,
		},
		{
			name:     "Empty password",
			email:    "test@example.com",
			password: "",
			expected: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := tests.MakeRegistrationRequest(
				ctx.Router,
				tc.email,
				tc.password,
				"Test",
				"User",
				"Test Org",
				"contact@test.com",
				"test",
			)

			if w.Code != tc.expected {
				t.Errorf("Expected status %d, got %d. Body: %s", tc.expected, w.Code, w.Body.String())
			}

			errorResp, err := tests.ParseErrorResponse(w)
			if err != nil {
				t.Errorf("Failed to parse error response: %v", err)
				return
			}

			// Verify error structure
			if errorResp["error"] == nil {
				t.Error("Expected error field in response")
			}
		})
	}
}

func TestAuthHandler_Register_SubdomainConflicts(t *testing.T) {
	ctx, err := tests.SetupTestContext()
	if err != nil {
		t.Fatalf("Failed to setup test context: %v", err)
	}
	defer tests.CleanupTestContext(ctx)

	// Create first organization
	firstOrg := &models.Organization{
		Name:         "First Organization",
		SubDomain:    "testorg",
		ContactEmail: "first@test.com",
		Active:       true,
		PlanType:     "basic",
	}
	if err := ctx.DB.Create(firstOrg).Error; err != nil {
		t.Fatalf("Failed to create first organization: %v", err)
	}

	t.Run("Subdomain already taken", func(t *testing.T) {
		w := tests.MakeRegistrationRequest(
			ctx.Router,
			"second@example.com",
			"password123",
			"Second",
			"User",
			"Second Organization",
			"second@test.com",
			"testorg", // Same subdomain
		)

		if w.Code != http.StatusConflict {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusConflict, w.Code, w.Body.String())
		}

		errorResp, err := tests.ParseErrorResponse(w)
		if err != nil {
			t.Errorf("Failed to parse error response: %v", err)
			return
		}

		// Verify error code
		errorData := errorResp["error"].(map[string]interface{})
		if errorData["details"].(map[string]interface{})["code"] != "SUBDOMAIN_TAKEN" {
			t.Error("Expected SUBDOMAIN_TAKEN error code")
		}
	})
}

func TestAuthHandler_Register_EmailConflicts(t *testing.T) {
	ctx, err := tests.SetupTestContext()
	if err != nil {
		t.Fatalf("Failed to setup test context: %v", err)
	}
	defer tests.CleanupTestContext(ctx)

	// Create first user successfully
	w1 := tests.MakeRegistrationRequest(
		ctx.Router,
		"test@example.com",
		"password123",
		"First",
		"User",
		"First Organization",
		"contact1@test.com",
		"firstorg",
	)

	if w1.Code != http.StatusCreated {
		t.Fatalf("Failed to create first user: %s", w1.Body.String())
	}

	t.Run("Same email in different organization should be allowed", func(t *testing.T) {
		w2 := tests.MakeRegistrationRequest(
			ctx.Router,
			"test@example.com", // Same email
			"password456",
			"Second",
			"User",
			"Second Organization",
			"contact2@test.com",
			"secondorg", // Different subdomain
		)

		if w2.Code != http.StatusCreated {
			t.Errorf("Expected to allow same email in different organization. Status: %d, Body: %s", w2.Code, w2.Body.String())
		}
	})
}

func TestAuthHandler_Register_TransactionRollback(t *testing.T) {
	ctx, err := tests.SetupTestContext()
	if err != nil {
		t.Fatalf("Failed to setup test context: %v", err)
	}
	defer tests.CleanupTestContext(ctx)

	t.Run("Invalid organization data should rollback transaction", func(t *testing.T) {
		// Create an organization that would conflict
		conflictOrg := &models.Organization{
			Name:         "Conflict Org",
			SubDomain:    "conflict",
			ContactEmail: "conflict@test.com",
			Active:       true,
			PlanType:     "basic",
		}
		if err := ctx.DB.Create(conflictOrg).Error; err != nil {
			t.Fatalf("Failed to create conflict organization: %v", err)
		}

		w := tests.MakeRegistrationRequest(
			ctx.Router,
			"user@example.com",
			"password123",
			"Test",
			"User",
			"Test Organization",
			"contact@test.com",
			"conflict", // This will cause conflict
		)

		if w.Code != http.StatusConflict {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusConflict, w.Code, w.Body.String())
		}

		// Verify no orphaned data was created
		var userCount int64
		ctx.DB.Model(&models.User{}).Where("email = ?", "user@example.com").Count(&userCount)
		if userCount != 0 {
			t.Error("User should not have been created when organization creation fails")
		}

		var roleCount int64
		ctx.DB.Model(&models.Role{}).Where("organization_id NOT IN (?)", ctx.DB.Model(&models.Organization{}).Select("id")).Count(&roleCount)
		if roleCount != 0 {
			t.Error("No orphaned roles should exist")
		}
	})
}

func TestAuthHandler_Register_TokenValidation(t *testing.T) {
	ctx, err := tests.SetupTestContext()
	if err != nil {
		t.Fatalf("Failed to setup test context: %v", err)
	}
	defer tests.CleanupTestContext(ctx)

	t.Run("Generated tokens should be valid", func(t *testing.T) {
		w := tests.MakeRegistrationRequest(
			ctx.Router,
			"token@example.com",
			"password123",
			"Token",
			"User",
			"Token Organization",
			"contact@token.com",
			"tokenorg",
		)

		if w.Code != http.StatusCreated {
			t.Fatalf("Registration failed: %s", w.Body.String())
		}

		resp, err := tests.ParseRegistrationResponse(w)
		if err != nil {
			t.Fatalf("Failed to parse registration response: %v", err)
		}

		// Validate access token
		claims, err := ctx.JWTService.ValidateToken(resp.AccessToken)
		if err != nil {
			t.Errorf("Access token should be valid: %v", err)
		} else {
			if claims.Email != "token@example.com" {
				t.Errorf("Expected email in token claims: 'token@example.com', got '%s'", claims.Email)
			}
			if claims.Role != "owner" {
				t.Errorf("Expected role in token claims: 'owner', got '%s'", claims.Role)
			}
			if !claims.IsAccessToken() {
				t.Error("Token should be identified as access token")
			}
		}

		// Validate refresh token
		refreshClaims, err := ctx.JWTService.ValidateToken(resp.RefreshToken)
		if err != nil {
			t.Errorf("Refresh token should be valid: %v", err)
		} else {
			if !refreshClaims.IsRefreshToken() {
				t.Error("Token should be identified as refresh token")
			}
		}
	})
} 