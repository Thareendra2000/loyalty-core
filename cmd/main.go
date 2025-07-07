package main

import (
	"fmt"
	"log"
	"net/http"

	"loyalty-core/config"
	"loyalty-core/models"
	"loyalty-core/routes"
	"loyalty-core/services"

	"github.com/joho/godotenv"
)

func createDemoData(cfg *config.Config) {
	fmt.Println("Creating demo data...")
	
	// Create auth service to add demo user
	authService := services.NewAuthService(cfg)
	
	// Demo user data
	demoUser := models.SignupRequest{
		Email:     "demo@loyalty.com",
		Password:  "password123",
		FirstName: "Demo",
		LastName:  "User",
	}
	
	// Try to create demo user
	response, err := authService.SignupUser(demoUser)
	if err != nil {
		if err.Error() == "user already exists" {
			fmt.Println("Demo user already exists")
		} else {
			fmt.Printf("⚠️  Failed to create demo user: %v\n", err)
		}
	} else {
		fmt.Println("Demo user created successfully")
		fmt.Printf("   Email: %s\n", response.User.Email)
		fmt.Printf("   Name: %s %s\n", response.User.FirstName, response.User.LastName)
	}
	
	fmt.Println("\nDemo Login Credentials:")
	fmt.Println("   Email: demo@loyalty.com")
	fmt.Println("   Password: password123")
	fmt.Println()
}

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

	// Create demo data for easy testing
	createDemoData(cfg)

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
