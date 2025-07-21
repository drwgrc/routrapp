package app

import (
	"routrapp-api/internal/api"
	"routrapp-api/internal/middleware"
)

// RegisterRoutes registers all application routes
func (a *App) RegisterRoutes() {
	// Health check endpoint
	healthHandler := api.NewHealthHandler(a.db)

	// User handler
	userHandler := api.NewUserHandler(a.db)

	// Auth handler with configured JWT service
	authHandler := api.NewAuthHandlerWithJWT(a.db, a.jwtService)

	// API group
	api := a.router.Group("/api")
	{
		// API version group
		v1 := api.Group("/v1")
		{
			// Health check endpoint
			v1.GET("/health", healthHandler.Check)

			// Auth endpoints (no authentication required for registration and login)
			auth := v1.Group("/auth")
			{
				auth.POST("/register", authHandler.RegisterOrganization)  // POST /api/v1/auth/register (organization registration)
				auth.POST("/register-user", authHandler.Register)         // POST /api/v1/auth/register-user (user registration to existing org)
				auth.POST("/login", authHandler.Login)                    // POST /api/v1/auth/login
				auth.POST("/refresh", authHandler.RefreshToken)           // POST /api/v1/auth/refresh
				auth.GET("/me", middleware.AuthMiddlewareWithJWT(a.jwtService), authHandler.GetCurrentUser) // GET /api/v1/auth/me (requires auth)
				auth.POST("/logout", middleware.AuthMiddlewareWithJWT(a.jwtService), authHandler.Logout) // POST /api/v1/auth/logout (requires auth)
				auth.POST("/change-password", middleware.AuthMiddlewareWithJWT(a.jwtService), authHandler.ChangePassword) // POST /api/v1/auth/change-password (requires auth)
			}

			// User endpoints
			users := v1.Group("/users")
			{
				users.GET("", userHandler.GetUsers)                                                 // GET /api/v1/users
				users.POST("", userHandler.CreateUser)                                             // POST /api/v1/users
				users.GET("/", userHandler.GetUserWithEmptyID)                                     // GET /api/v1/users/ - Bad request
				users.GET("/:id", userHandler.GetUser)                                             // GET /api/v1/users/:id
				users.PUT("/profile", middleware.AuthMiddlewareWithJWT(a.jwtService), userHandler.UpdateProfile)     // PUT /api/v1/users/profile (requires auth)
			}
			
			// Panic endpoint for testing recovery middleware
			v1.GET("/panic", userHandler.TriggerPanic) // GET /api/v1/panic
		}
	}
} 