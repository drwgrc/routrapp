package app

import (
	"routrapp-api/internal/handlers"
)

// RegisterRoutes registers all application routes
func (a *App) RegisterRoutes() {
	// Health check endpoint
	healthHandler := handlers.NewHealthHandler(a.db)
	a.router.GET("/health", healthHandler.Check)

	// API group
	api := a.router.Group("/api")
	{
		// API version group
		v1 := api.Group("/v1")
		_ = v1 // Placeholder until we add endpoints
	}
} 