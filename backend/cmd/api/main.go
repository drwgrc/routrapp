package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"routrapp-api/internal/app"
	"routrapp-api/internal/config"
	"routrapp-api/internal/logger"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize application
	application, err := app.NewApp(cfg)
	if err != nil {
		logger.Fatalf("❌ Failed to initialize application: %v", err)
	}

	// Start server in a goroutine
	go func() {
		if err := application.Start(); err != nil {
			logger.Errorf("❌ Server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("🛑 Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := application.Shutdown(ctx); err != nil {
		logger.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	logger.Info("�� Server exited")
}
