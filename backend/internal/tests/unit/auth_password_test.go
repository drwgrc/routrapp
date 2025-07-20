package unit_test

import (
	"testing"

	"routrapp-api/internal/utils/auth"
)

func TestValidatePassword(t *testing.T) {
	testCases := []struct {
		name     string
		password string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "Valid strong password",
			password: "StrongPass123!",
			wantErr:  false,
		},
		{
			name:     "Valid password with symbols",
			password: "MyP@ssw0rd$",
			wantErr:  false,
		},
		{
			name:     "Empty password",
			password: "",
			wantErr:  true,
			errMsg:   "password cannot be empty",
		},
		{
			name:     "Too short",
			password: "Short1!",
			wantErr:  true,
			errMsg:   "password must be at least 8 characters long",
		},
		{
			name:     "No uppercase letter",
			password: "lowercase123!",
			wantErr:  true,
			errMsg:   "password must contain at least one uppercase letter",
		},
		{
			name:     "No lowercase letter",
			password: "UPPERCASE123!",
			wantErr:  true,
			errMsg:   "password must contain at least one lowercase letter",
		},
		{
			name:     "No number",
			password: "NoNumbersHere!",
			wantErr:  true,
			errMsg:   "password must contain at least one number",
		},
		{
			name:     "No special character",
			password: "NoSpecialChar123",
			wantErr:  true,
			errMsg:   "password must contain at least one special character",
		},
		{
			name:     "Too long password",
			password: "ThisPasswordIsWayTooLongAndExceedsTheMaximumLengthOfTwoHundredAndFiftyFiveCharactersWhichIsDefinedAsTheMaximumLengthForAPasswordInOurSystemAndThisStringContinuesToGoOnAndOnUntilItReachesTheRequiredLengthToTestTheTooLongValidationRuleThisPasswordIsWayTooLongAndExceedsTheMaximumLengthOf255Chars123!@",
			wantErr:  true,
			errMsg:   "password must not exceed 255 characters",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := auth.ValidatePassword(tc.password)
			
			if tc.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if err.Error() != tc.errMsg {
					t.Errorf("Expected error message '%s', got '%s'", tc.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestValidatePasswordWithCustomRequirements(t *testing.T) {
	// Test with custom requirements (no special chars required)
	requirements := auth.PasswordRequirements{
		MinLength:      6,
		MaxLength:      20,
		RequireUpper:   true,
		RequireLower:   true,
		RequireDigit:   true,
		RequireSpecial: false,
	}

	testCases := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid with custom requirements",
			password: "SimplePass123",
			wantErr:  false,
		},
		{
			name:     "Too short for custom requirements",
			password: "Ab1",
			wantErr:  true,
		},
		{
			name:     "Too long for custom requirements",
			password: "ThisPasswordIsTooLongForCustomRequirements123",
			wantErr:  true,
		},
		{
			name:     "Missing uppercase",
			password: "simplepass123",
			wantErr:  true,
		},
		{
			name:     "Valid without special chars (allowed by custom)",
			password: "ValidPass123",
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := auth.ValidatePasswordWithRequirements(tc.password, requirements)
			
			if tc.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			} else if !tc.wantErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestGetPasswordStrength(t *testing.T) {
	testCases := []struct {
		name             string
		password         string
		expectedStrength string
	}{
		{
			name:             "Very weak password",
			password:         "123",
			expectedStrength: "Very Weak",
		},
		{
			name:             "Weak password",
			password:         "password",
			expectedStrength: "Weak",
		},
		{
			name:             "Medium password",
			password:         "Password1",
			expectedStrength: "Weak", // Only has 3 points (length, upper, lower, digit) - common password penalty doesn't apply
		},
		{
			name:             "Strong password",
			password:         "StrongPass123!",
			expectedStrength: "Very Strong", // Has 6 points (length >=8, length >=12, upper, lower, digit, special)
		},
		{
			name:             "Very strong password",
			password:         "VeryStr0ngP@ssw0rd2024!",
			expectedStrength: "Very Strong",
		},
		{
			name:             "Common password penalized",
			password:         "Password123!",
			expectedStrength: "Medium", // Penalized for containing "password"
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			strength := auth.GetPasswordStrength(tc.password)
			if strength != tc.expectedStrength {
				t.Errorf("Expected strength '%s', got '%s'", tc.expectedStrength, strength)
			}
		})
	}
}

func TestIsCommonPassword(t *testing.T) {
	testCases := []struct {
		name     string
		password string
		expected bool
	}{
		{
			name:     "Common password - password",
			password: "password",
			expected: true,
		},
		{
			name:     "Common password - 123456",
			password: "123456",
			expected: true,
		},
		{
			name:     "Common password case insensitive",
			password: "PASSWORD",
			expected: true,
		},
		{
			name:     "Common password - qwerty",
			password: "qwerty",
			expected: true,
		},
		{
			name:     "Strong unique password",
			password: "StrongUniquePass123!",
			expected: false,
		},
		{
			name:     "Not common password",
			password: "MySecur3P@ss",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := auth.IsCommonPassword(tc.password)
			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestHashPasswordAndVerify(t *testing.T) {
	password := "TestPassword123!"
	
	// Test hashing
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	
	if hashedPassword == "" {
		t.Error("Hashed password should not be empty")
	}
	
	if hashedPassword == password {
		t.Error("Hashed password should be different from original")
	}
	
	// Test verification with correct password
	err = auth.VerifyPassword(password, hashedPassword)
	if err != nil {
		t.Errorf("Password verification failed: %v", err)
	}
	
	// Test verification with wrong password
	err = auth.VerifyPassword("WrongPassword", hashedPassword)
	if err == nil {
		t.Error("Expected verification to fail with wrong password")
	}
	
	// Test with empty password
	_, err = auth.HashPassword("")
	if err == nil {
		t.Error("Expected error when hashing empty password")
	}
	
	// Test verification with empty password
	err = auth.VerifyPassword("", hashedPassword)
	if err == nil {
		t.Error("Expected error when verifying empty password")
	}
	
	// Test verification with empty hash
	err = auth.VerifyPassword(password, "")
	if err == nil {
		t.Error("Expected error when verifying with empty hash")
	}
}

func TestIsValidPassword_Legacy(t *testing.T) {
	testCases := []struct {
		name     string
		password string
		expected bool
	}{
		{
			name:     "Valid 8+ character password",
			password: "password",
			expected: true,
		},
		{
			name:     "Too short password",
			password: "pass",
			expected: false,
		},
		{
			name:     "Empty password",
			password: "",
			expected: false,
		},
		{
			name:     "Exactly 8 characters",
			password: "12345678",
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := auth.IsValidPassword(tc.password)
			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestDefaultPasswordRequirements(t *testing.T) {
	req := auth.DefaultPasswordRequirements()
	
	if req.MinLength != 8 {
		t.Errorf("Expected MinLength 8, got %d", req.MinLength)
	}
	
	if req.MaxLength != 255 {
		t.Errorf("Expected MaxLength 255, got %d", req.MaxLength)
	}
	
	if !req.RequireUpper {
		t.Error("Expected RequireUpper to be true")
	}
	
	if !req.RequireLower {
		t.Error("Expected RequireLower to be true")
	}
	
	if !req.RequireDigit {
		t.Error("Expected RequireDigit to be true")
	}
	
	if !req.RequireSpecial {
		t.Error("Expected RequireSpecial to be true")
	}
} 