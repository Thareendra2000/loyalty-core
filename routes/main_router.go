package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"loyalty-core/config"
)

type MainRouter struct {
	cfg           *config.Config
	authRoutes    *AuthRoutes
	loyaltyRoutes *LoyaltyRoutes
}

func NewMainRouter(cfg *config.Config) *MainRouter {
	return &MainRouter{
		cfg:           cfg,
		authRoutes:    NewAuthRoutes(cfg),
		loyaltyRoutes: NewLoyaltyRoutes(cfg),
	}
}

func (mr *MainRouter) RegisterAllRoutes() {
	// Register auth routes
	mr.authRoutes.RegisterRoutes()

	// Register loyalty routes
	mr.loyaltyRoutes.RegisterRoutes()

	// Register general routes
	mr.registerGeneralRoutes()
}

func (mr *MainRouter) registerGeneralRoutes() {
	// Root endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"message": "Loyalty Core API Server is running!",
			"status": "healthy",
			"port": "%s",
			"endpoint": "%s"
		}`, mr.cfg.Port, r.URL.Path)
	})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "healthy",
			"message": "Server is up and running",
		})
	})

	// API info endpoint
	http.HandleFunc("/api/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"service":     "loyalty-core",
			"version":     "1.0.0",
			"description": "Loyalty Core API Service",
		})
	})

	// API endpoints info
	http.HandleFunc("/api/endpoints", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		endpoints := map[string]interface{}{
			"auth": map[string]string{
				"signup":  "POST /api/auth/signup",
				"login":   "POST /api/auth/login",
				"profile": "GET /api/auth/profile",
			},
			"loyalty": map[string]string{
				"earn":    "POST /api/loyalty/earn",
				"redeem":  "POST /api/loyalty/redeem",
				"balance": "GET /api/loyalty/balance",
				"history": "GET /api/loyalty/history",
			},
			"general": map[string]string{
				"health": "GET /health",
				"info":   "GET /api/info",
			},
		}
		json.NewEncoder(w).Encode(endpoints)
	})
}

func (mr *MainRouter) GetAuthService() *AuthRoutes {
	return mr.authRoutes
}
