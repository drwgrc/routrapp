package app

import (
	"context"
	"fmt"
	"net/http"

	"routrapp-api/internal/config"
	"routrapp-api/internal/middleware"
	"routrapp-api/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	config *config.Config
	server *http.Server
	db     *gorm.DB
	router *gin.Engine
}

func NewApp(cfg *config.Config) (*App, error) {
	if cfg == nil {
		cfg = config.Load()
	}

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
		return nil, err
	}
	app.db = db

	// Auto-migrate models in development environment
	if app.config.Environment == "development" {
		err = app.db.AutoMigrate(
			&models.Tenant{},
			&models.User{},
			&models.Technician{},
			&models.Route{},
			&models.RouteStop{},
			&models.RouteActivity{},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to auto migrate: %w", err)
		}
	}

	app.setupRouter()
	app.RegisterRoutes() // Register all routes
	app.setupServer()
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
	
	// Add only essential middleware
	a.router.Use(middleware.CORSMiddleware(a.config))
	
	// Root endpoint
	a.router.GET("/", func(c *gin.Context) {
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
	return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	// Close database connection if needed
	sqlDB, err := a.db.DB()
	if err == nil {
		sqlDB.Close()
	}
	
	return a.server.Shutdown(ctx)
}