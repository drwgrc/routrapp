package app

import (
	"context"
	"database/sql"
	"net/http"

	"routrapp-api/internal/config"
	"routrapp-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

type App struct {
	config *config.Config
	server *http.Server
	db     *sql.DB
	router *gin.Engine
}

func NewApp(cfg *config.Config) (*App, error) {
	if cfg == nil {
		cfg = config.Load()
	}

	app := &App{
		config: cfg,
	}

	app.setupRouter()
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
	a.router = gin.Default()
	
	// Add CORS middleware
	a.router.Use(middleware.CORSMiddleware(a.config))
	
	// Root endpoint
	a.router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "routrapp-api",
		})
	})
}

func (a *App) Start() error {
	return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}