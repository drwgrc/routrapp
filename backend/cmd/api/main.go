package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"routrapp-api/internal/app"
	"routrapp-api/internal/config"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize application
	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatal("❌ Failed to initialize application:", err)
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 Starting server on port %s", cfg.Server.Port)
		if err := application.Start(); err != nil {
			log.Printf("❌ Server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := application.Shutdown(ctx); err != nil {
		log.Fatal("❌ Server forced to shutdown:", err)
	}

	log.Println("👋 Server exited")
}
