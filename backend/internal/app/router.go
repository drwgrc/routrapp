package app

import (
	"routrapp-api/internal/api"
)

// RegisterRoutes registers all application routes
func (a *App) RegisterRoutes() {
	// Health check endpoint
	healthHandler := api.NewHealthHandler(a.db)

	// User handler for testing error scenarios
	userHandler := api.NewUserHandler()

	// API group
	api := a.router.Group("/api")
	{
		// API version group
		v1 := api.Group("/v1")
		{
			// Health check endpoint
			v1.GET("/health", healthHandler.Check)

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