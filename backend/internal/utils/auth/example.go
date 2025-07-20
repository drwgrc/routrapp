package auth

import (
	"fmt"
	"log"
)

// ExampleUsage demonstrates how to use the JWT authentication system
func ExampleUsage() {
	// Create a JWT service
	jwtService := DefaultJWTService()
	
	// Sample user data
	userID := uint(123)
	organizationID := uint(456)
	email := "john.doe@example.com"
	role := "owner"
	
	// Generate access token
	accessToken, err := jwtService.GenerateAccessToken(userID, organizationID, email, role)
	if err != nil {
		log.Fatalf("Failed to generate access token: %v", err)
	}
	fmt.Printf("Generated Access Token: %s\n", accessToken[:50]+"...")
	
	// Generate refresh token
	refreshToken, err := jwtService.GenerateRefreshToken(userID, organizationID, email, role)
	if err != nil {
		log.Fatalf("Failed to generate refresh token: %v", err)
	}
	fmt.Printf("Generated Refresh Token: %s\n", refreshToken[:50]+"...")
	
	// Validate the access token
	claims, err := jwtService.ValidateToken(accessToken)
	if err != nil {
		log.Fatalf("Failed to validate access token: %v", err)
	}
	
	// Print the claims
	fmt.Printf("Token Claims:\n")
	fmt.Printf("  User ID: %d\n", claims.UserID)
	fmt.Printf("  Organization ID: %d\n", claims.OrganizationID)
	fmt.Printf("  Email: %s\n", claims.Email)
	fmt.Printf("  Role: %s\n", claims.Role)
	fmt.Printf("  Token Type: %s\n", claims.TokenType)
	fmt.Printf("  Is Access Token: %t\n", claims.IsAccessToken())
	fmt.Printf("  Is Refresh Token: %t\n", claims.IsRefreshToken())
	fmt.Printf("  Expires At: %s\n", claims.ExpiresAt.Time.Format("2006-01-02 15:04:05"))
	
	// Get user context
	userContext := claims.GetUserContext()
	fmt.Printf("User Context: %+v\n", userContext)
	
	// Test header extraction
	authHeader := "Bearer " + accessToken
	extractedToken, err := ExtractTokenFromHeader(authHeader)
	if err != nil {
		log.Fatalf("Failed to extract token from header: %v", err)
	}
	fmt.Printf("Extracted token matches: %t\n", extractedToken == accessToken)
}

/*
Example usage in a Gin handler:

func loginHandler(c *gin.Context) {
    // ... validate user credentials ...
    
    jwtService := auth.DefaultJWTService()
    
    // Generate tokens
    accessToken, err := jwtService.GenerateAccessToken(user.ID, user.OrganizationID, user.Email, user.Role.Name)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to generate access token"})
        return
    }
    
    refreshToken, err := jwtService.GenerateRefreshToken(user.ID, user.OrganizationID, user.Email, user.Role.Name)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to generate refresh token"})
        return
    }
    
    c.JSON(200, gin.H{
        "access_token":  accessToken,
        "refresh_token": refreshToken,
        "expires_in":    constants.JWT_ACCESS_TOKEN_EXPIRY,
    })
}

Example middleware usage:

// In your router setup:
protected := r.Group("/api/v1/protected")
protected.Use(middleware.AuthMiddleware())
{
    protected.GET("/profile", profileHandler)
}

// In your handler:
func profileHandler(c *gin.Context) {
    userID, _ := middleware.GetUserID(c)
    email, _ := middleware.GetUserEmail(c)
    role, _ := middleware.GetUserRole(c)
    
    c.JSON(200, gin.H{
        "user_id": userID,
        "email":   email,
        "role":    role,
    })
}
*/ 