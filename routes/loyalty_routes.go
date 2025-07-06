package routes

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"loyalty-core/config"
	"loyalty-core/models"
	"loyalty-core/services"
	"loyalty-core/utils"
)

type LoyaltyRoutes struct {
	loyaltyService *services.LoyaltyService
	authService    *services.AuthService
	config         *config.Config
}

func NewLoyaltyRoutes(cfg *config.Config) *LoyaltyRoutes {
	authService := services.NewAuthService(cfg)
	return &LoyaltyRoutes{
		loyaltyService: services.NewLoyaltyService(cfg),
		authService:    authService,
		config:         cfg,
	}
}

// EarnPoints handles earning points
func (lr *LoyaltyRoutes) EarnPoints(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	// Get user ID from token
	userID, err := lr.getUserIDFromToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	var req models.EarnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Validate points
	if req.Points <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Points must be greater than 0"})
		return
	}

	transaction, err := lr.loyaltyService.EarnPoints(userID, req.Points, req.Description)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transaction)
}

// RedeemPoints handles redeeming points
func (lr *LoyaltyRoutes) RedeemPoints(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	// Get user ID from token
	userID, err := lr.getUserIDFromToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	var req models.RedeemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Validate points
	if req.Points <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Points must be greater than 0"})
		return
	}

	transaction, err := lr.loyaltyService.RedeemPoints(userID, req.Points, req.Description)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transaction)
}

// GetBalance handles getting user's loyalty balance
func (lr *LoyaltyRoutes) GetBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	// Get user ID from token
	userID, err := lr.getUserIDFromToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	balance, err := lr.loyaltyService.GetBalance(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(balance)
}

// GetHistory handles getting user's transaction history
func (lr *LoyaltyRoutes) GetHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	// Get user ID from token
	userID, err := lr.getUserIDFromToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Get limit from query parameter
	limitStr := r.URL.Query().Get("limit")
	limit := 10 // default
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	transactions, err := lr.loyaltyService.GetTransactionHistory(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Apply limit
	if len(transactions) > limit {
		transactions = transactions[:limit]
	}

	response := map[string]interface{}{
		"transactions": transactions,
		"count":        len(transactions),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// getUserIDFromToken extracts user ID from JWT token
func (lr *LoyaltyRoutes) getUserIDFromToken(r *http.Request) (string, error) {
	// Get Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header required")
	}

	// Extract token from Bearer header
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return "", errors.New("invalid authorization header format")
	}

	// Validate token
	claims, err := utils.ValidateToken(tokenString, lr.config.JWTSecret)
	if err != nil {
		return "", errors.New("invalid token")
	}

	return claims.UserID, nil
}

// RegisterRoutes registers all loyalty routes
func (lr *LoyaltyRoutes) RegisterRoutes() {
	http.HandleFunc("/api/loyalty/earn", lr.EarnPoints)
	http.HandleFunc("/api/loyalty/redeem", lr.RedeemPoints)
	http.HandleFunc("/api/loyalty/balance", lr.GetBalance)
	http.HandleFunc("/api/loyalty/history", lr.GetHistory)

	log.Println("Loyalty routes registered")
}

// GetLoyaltyService returns the loyalty service instance (for testing)
func (lr *LoyaltyRoutes) GetLoyaltyService() *services.LoyaltyService {
	return lr.loyaltyService
}
