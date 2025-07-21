package app

import (
	"context"
	"fmt"
	"net/http"

	"routrapp-api/internal/config"
	"routrapp-api/internal/logger"
	"routrapp-api/internal/middleware"
	"routrapp-api/internal/models"
	"routrapp-api/internal/utils/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	config     *config.Config
	server     *http.Server
	db         *gorm.DB
	router     *gin.Engine
	jwtService *auth.JWTService
}

func NewApp(cfg *config.Config) (*App, error) {
	if cfg == nil {
		cfg = config.Load()
	}

	// Initialize logger first
	logger.InitLogger(cfg.Environment)
	logger.Infof("Starting application in %s environment", cfg.Environment)

	// Set Gin mode based on environment before creating any engine
	switch cfg.Environment {
	case "production", "staging":
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	app := &App{
		config: cfg,
	}

	// Initialize database connection
	db, err := config.InitDatabase(&cfg.Database)
	if err != nil {
		logger.Errorf("Failed to initialize database: %v", err)
		return nil, err
	}
	app.db = db
	logger.Info("Database connection established")

	// Initialize JWT service with configuration
	app.jwtService = auth.NewJWTService(cfg.JWT.Secret)
	logger.Infof("JWT service initialized with secret from configuration")

	// Auto-migrate models in development environment
	if app.config.Environment == "development" {
		logger.Info("Running database migrations for development environment")
		err = app.db.AutoMigrate(models.AllModels()...)
		if err != nil {
			logger.Errorf("Failed to auto migrate: %v", err)
			return nil, fmt.Errorf("failed to auto migrate: %w", err)
		}
		logger.Info("Database migrations completed")
	}

	app.setupRouter()
	app.RegisterRoutes() // Register all routes
	app.setupServer()
	logger.Infof("Application initialized successfully on port %s", cfg.Server.Port)
	return app, nil
}

func (a *App) setupServer() {
	a.server = &http.Server{
		Addr:         ":" + a.config.Server.Port,
		Handler:      a.router,
		ReadTimeout:  a.config.Server.ReadTimeout,
		WriteTimeout: a.config.Server.WriteTimeout,
	}
}

func (a *App) setupRouter() {
	// Create a new Gin engine without any default middleware
	a.router = gin.New()
	
	// Add middleware stack in correct order
	a.router.Use(middleware.RecoveryMiddleware())       // First: Recovery from panics
	a.router.Use(middleware.RequestIDMiddleware())      // Second: Add request IDs
	a.router.Use(middleware.LoggerMiddleware())         // Third: Log requests
	a.router.Use(middleware.CORSMiddleware(a.config))   // Fourth: CORS handling
	a.router.Use(middleware.ErrorHandlerMiddleware())   // Last: Error handling
	
	// Root endpoint
	a.router.GET("/", func(c *gin.Context) {
		logger.WithContext(c).Info("Root endpoint accessed")
		c.JSON(200, gin.H{
			"message": "routrapp-api",
			"status": "running",
		})
	})
}

// GetDB returns the database connection
func (a *App) GetDB() *gorm.DB {
	return a.db
}

func (a *App) Start() error {
	logger.Infof("🚀 Starting server on port %s", a.config.Server.Port)
	return a.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (a *App) Shutdown(ctx context.Context) error {
	logger.Info("🛑 Shutting down server...")
	return a.server.Shutdown(ctx)
}