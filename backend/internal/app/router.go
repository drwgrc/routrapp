package app

import (
	"routrapp-api/internal/handlers"
)

// RegisterRoutes registers all application routes
func (a *App) RegisterRoutes() {
	// Health check endpoint
	healthHandler := handlers.NewHealthHandler(a.db)
	a.router.GET("/health", healthHandler.Check)

	// User handler for testing error scenarios
	userHandler := handlers.NewUserHandler()

	// API group
	api := a.router.Group("/api")
	{
		// API version group
		v1 := api.Group("/v1")
		{
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