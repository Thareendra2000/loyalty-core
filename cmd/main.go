package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default values")
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create a simple HTTP handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"message": "Loyalty Core API Server is running!",
			"status": "healthy",
			"port": "%s",
			"endpoint": "%s"
		}`, port, r.URL.Path)
	})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"status": "healthy",
			"message": "Server is up and running"
		}`)
	})

	// API info endpoint
	http.HandleFunc("/api/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"service": "loyalty-core",
			"version": "1.0.0",
			"description": "Loyalty Core API Service"
		}`)
	})

	// Start the server
	serverAddr := ":" + port
	fmt.Printf("Server starting on port %s...\n", port)
	fmt.Printf("Server URL: http://localhost:%s\n", port)
	fmt.Printf("Health check: http://localhost:%s/health\n", port)
	fmt.Printf("API info: http://localhost:%s/api/info\n", port)
	
	log.Printf("Server is listening on port %s", port)
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}