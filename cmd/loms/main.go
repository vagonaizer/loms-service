package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vagonaizer/loms/internal/app"
	"github.com/vagonaizer/loms/internal/config"
)

const configPath = "config/config.yaml"

func main() {
	log.Println("Starting LOMS service...")

	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	log.Printf("Configuration loaded, gRPC port: %d", cfg.GRPC.Port)

	// Initialize application
	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}
	log.Println("Application initialized successfully")

	// Start application
	go func() {
		log.Printf("Starting gRPC server on port %d...", cfg.GRPC.Port)
		if err := application.Run(); err != nil {
			log.Fatalf("failed to start application: %v", err)
		}
	}()

	log.Println("Server is running. Press Ctrl+C to stop.")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	// Graceful shutdown
	application.Stop()
	log.Println("Server stopped")
}
