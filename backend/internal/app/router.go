package app

import (
	"routrapp-api/internal/api"
	"routrapp-api/internal/middleware"
)

// RegisterRoutes registers all application routes
func (a *App) RegisterRoutes() {
	// Health check endpoint
	healthHandler := api.NewHealthHandler(a.db)

	// User handler for testing error scenarios
	userHandler := api.NewUserHandler()

	// Auth handler
	authHandler := api.NewAuthHandler(a.db)

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
				auth.POST("/register", authHandler.Register)             // POST /api/v1/auth/register
				auth.POST("/login", authHandler.Login)                   // POST /api/v1/auth/login
				auth.POST("/refresh", authHandler.RefreshToken)          // POST /api/v1/auth/refresh
				auth.POST("/logout", middleware.AuthMiddleware(), authHandler.Logout) // POST /api/v1/auth/logout (requires auth)
				auth.POST("/change-password", middleware.AuthMiddleware(), authHandler.ChangePassword) // POST /api/v1/auth/change-password (requires auth)
			}

			// User endpoints for testing error scenarios
			users := v1.Group("/users")
			{
				users.GET("", userHandler.GetUsers)                    // GET /api/v1/users
				users.POST("", userHandler.CreateUser)                // POST /api/v1/users
				users.GET("/", userHandler.GetUserWithEmptyID)        // GET /api/v1/users/ - Bad request
				users.GET("/:id", userHandler.GetUser)                // GET /api/v1/users/:id
			}
			
			// Panic endpoint for testing recovery middleware
			v1.GET("/panic", userHandler.TriggerPanic) // GET /api/v1/panic
		}
	}
} 