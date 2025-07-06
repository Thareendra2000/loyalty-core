package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"loyalty-core/config"
	"loyalty-core/models"
	"loyalty-core/services"
)

type AuthRoutes struct {
	authService *services.AuthService
}

func NewAuthRoutes(cfg *config.Config) *AuthRoutes {
	return &AuthRoutes{
		authService: services.NewAuthService(cfg),
	}
}

// Signup handles user registration
func (ar *AuthRoutes) Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	var req models.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	response, err := ar.authService.SignupUser(req)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "user already exists" {
			status = http.StatusConflict
		} else if err.Error() == "internal server error" {
			status = http.StatusInternalServerError
		}

		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login handles user authentication
func (ar *AuthRoutes) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	response, err := ar.authService.LoginUser(req)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "invalid credentials" {
			status = http.StatusUnauthorized
		} else if err.Error() == "internal server error" {
			status = http.StatusInternalServerError
		}

		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// RegisterRoutes registers all auth routes
func (ar *AuthRoutes) RegisterRoutes() {
	http.HandleFunc("/api/auth/signup", ar.Signup)
	http.HandleFunc("/api/auth/login", ar.Login)
	log.Println("Auth routes registered")
}

// GetAuthService returns the auth service instance (for testing)
func (ar *AuthRoutes) GetAuthService() *services.AuthService {
	return ar.authService
}
