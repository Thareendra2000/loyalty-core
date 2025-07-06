package main

import (
	"fmt"
	"log"
	"net/http"

	"loyalty-core/config"
	"loyalty-core/routes"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default values")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Create main router
	mainRouter := routes.NewMainRouter(cfg)

	// Register all routes
	mainRouter.RegisterAllRoutes()

	// Start the server
	serverAddr := ":" + cfg.Port
	fmt.Printf("Server starting on port %s...\n", cfg.Port)

	log.Printf("Server is listening on port %s", cfg.Port)
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
