package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"loyalty-core/config"
	"loyalty-core/models"
	"loyalty-core/storage"

	square "github.com/square/square-go-sdk"
)

type LoyaltyService struct {
	config        *config.Config
	userStorage   *storage.UserStorage
	squareService *SquareService
	transactions  map[string][]models.Transaction
	instanceID    string // Debug: track service instance
}

func NewLoyaltyService(cfg *config.Config) *LoyaltyService {
	var squareService *SquareService

	// Try to initialize Square service, but don't fail if it's not available
	if cfg.SquareAccessToken != "" && cfg.SquareAccessToken != "your-square-access-token" &&
		cfg.SquareLocationID != "" && cfg.SquareLocationID != "your-square-location-id" {
		var err error
		squareService, err = NewSquareService(cfg)
		if err != nil {
			log.Printf("Warning: Square service not available: %v", err)
			log.Printf("Running in fallback mode with in-memory storage")
		} else {
			log.Printf("Square service initialized successfully")
		}
	} else {
		log.Printf("Square credentials not configured, running in fallback mode")
	}

	service := &LoyaltyService{
		config:        cfg,
		userStorage:   storage.GetGlobalUserStorage(),
		squareService: squareService,
		transactions:  make(map[string][]models.Transaction),
	}

	return service
}

func (s *LoyaltyService) EarnPoints(userID string, points int, description string) (*models.Transaction, error) {
	user, err := s.userStorage.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Ensure user has a Square loyalty account
	if err := s.ensureSquareLoyaltyAccount(user); err != nil {
		return nil, fmt.Errorf("failed to ensure Square loyalty account: %w", err)
	}

	// Create transaction
	transaction := models.Transaction{
		ID:          s.generateID(),
		UserID:      userID,
		Type:        "earn",
		Points:      points,
		Description: description,
		CreatedAt:   time.Now(),
	}

	// If Square service is available, accumulate points in Square
	if s.squareService != nil {
		orderID := "order-" + transaction.ID // Create a mock order ID
		_, err := s.squareService.AccumulateLoyaltyPoints(user.LoyaltyID, points, orderID)
		if err != nil {
			return nil, fmt.Errorf("failed to accumulate points in Square: %w", err)
		}
	}

	// Update user points locally (for fallback)
	user.Points += points
	user.UpdatedAt = time.Now()

	// Save updated user
	if err := s.userStorage.UpdateUser(user); err != nil {
		return nil, err
	}

	// Store transaction
	s.transactions[userID] = append(s.transactions[userID], transaction)

	return &transaction, nil
}

func (s *LoyaltyService) RedeemPoints(userID string, points int, description string) (*models.Transaction, error) {
	user, err := s.userStorage.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Ensure user has a Square loyalty account
	if err := s.ensureSquareLoyaltyAccount(user); err != nil {
		return nil, fmt.Errorf("failed to ensure Square loyalty account: %w", err)
	}

	// Check balance from Square if available
	if s.squareService != nil {
		account, err := s.squareService.GetLoyaltyAccount(user.LoyaltyID)
		if err != nil {
			return nil, fmt.Errorf("failed to get Square account balance: %w", err)
		}

		if account.Balance != nil && *account.Balance < points {
			return nil, errors.New("insufficient points")
		}
	} else {
		// Fallback to local balance check
		if user.Points < points {
			return nil, errors.New("insufficient points")
		}
	}

	// Create transaction
	transaction := models.Transaction{
		ID:          s.generateID(),
		UserID:      userID,
		Type:        "redeem",
		Points:      points,
		Description: description,
		CreatedAt:   time.Now(),
	}

	// If Square service is available, redeem points in Square
	if s.squareService != nil {
		// Use adjust points to subtract points (negative value)
		_, err := s.squareService.AdjustLoyaltyPoints(user.LoyaltyID, -points, description)
		if err != nil {
			return nil, fmt.Errorf("failed to redeem points in Square: %w", err)
		}
	}

	// Update user points locally (for fallback)
	user.Points -= points
	user.UpdatedAt = time.Now()

	// Save updated user
	if err := s.userStorage.UpdateUser(user); err != nil {
		return nil, err
	}

	// Store transaction
	s.transactions[userID] = append(s.transactions[userID], transaction)

	return &transaction, nil
}

func (s *LoyaltyService) GetBalance(userID string) (*models.BalanceResponse, error) {
	user, err := s.userStorage.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Ensure user has a Square loyalty account
	if err := s.ensureSquareLoyaltyAccount(user); err != nil {
		return nil, fmt.Errorf("failed to ensure Square loyalty account: %w", err)
	}

	var balance int
	var transactions []models.Transaction

	// Get balance from Square if available
	if s.squareService != nil {
		account, err := s.squareService.GetLoyaltyAccount(user.LoyaltyID)
		if err != nil {
			return nil, fmt.Errorf("failed to get Square account balance: %w", err)
		}

		if account.Balance != nil {
			balance = *account.Balance
		}

		// Get transaction history from Square
		squareTransactions, err := s.GetTransactionHistory(userID)
		if err == nil {
			transactions = squareTransactions
		}
	} else {
		// Fallback to local balance
		balance = user.Points
		transactions = s.transactions[userID]
	}

	if transactions == nil {
		transactions = []models.Transaction{}
	}

	return &models.BalanceResponse{
		Points:       balance,
		Transactions: transactions,
	}, nil
}

func (s *LoyaltyService) GetTransactionHistory(userID string) ([]models.Transaction, error) {
	user, err := s.userStorage.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Ensure user has a Square loyalty account
	if err := s.ensureSquareLoyaltyAccount(user); err != nil {
		return nil, fmt.Errorf("failed to ensure Square loyalty account: %w", err)
	}

	// Get transaction history from Square if available
	if s.squareService != nil {
		events, err := s.squareService.SearchLoyaltyEvents(user.LoyaltyID, 50)
		if err != nil {
			return nil, fmt.Errorf("failed to get Square transaction history: %w", err)
		}

		// Convert Square events to our Transaction model
		transactions := make([]models.Transaction, 0, len(events))
		for _, event := range events {
			transaction := s.convertSquareEventToTransaction(event, userID)
			if transaction != nil {
				transactions = append(transactions, *transaction)
			}
		}

		return transactions, nil
	}

	// Fallback to local transactions
	transactions := s.transactions[userID]
	if transactions == nil {
		return []models.Transaction{}, nil
	}

	return transactions, nil
}

// ensureSquareLoyaltyAccount ensures the user has a Square loyalty account
func (s *LoyaltyService) ensureSquareLoyaltyAccount(user *models.User) error {
	if s.squareService == nil {
		return nil // Skip if Square service is not available
	}

	if user.LoyaltyID != "" {
		return nil // User already has a loyalty account
	}

	// Create a phone number for the user (required by Square)
	// In a real application, you would collect this from the user
	phoneNumber := "+1555" + strconv.Itoa(int(time.Now().Unix())%10000000) // Generate a mock phone number

	// Create loyalty account in Square
	account, err := s.squareService.CreateLoyaltyAccount(phoneNumber, user.FirstName, user.LastName)
	if err != nil {
		return fmt.Errorf("failed to create Square loyalty account: %w", err)
	}

	// Update user with the loyalty account ID
	user.LoyaltyID = *account.ID
	user.UpdatedAt = time.Now()

	// Save the updated user
	if err := s.userStorage.UpdateUser(user); err != nil {
		return fmt.Errorf("failed to save user with loyalty ID: %w", err)
	}

	return nil
}

// convertSquareEventToTransaction converts a Square loyalty event to our Transaction model
func (s *LoyaltyService) convertSquareEventToTransaction(event *square.LoyaltyEvent, userID string) *models.Transaction {
	if event == nil || event.ID == "" {
		return nil
	}

	var transactionType string
	var points int
	var description string

	// Determine transaction type and points based on Square event
	switch event.Type {
	case "ACCUMULATE_POINTS":
		transactionType = "earn"
		if event.AccumulatePoints != nil && event.AccumulatePoints.Points != nil {
			points = *event.AccumulatePoints.Points
		}
		description = "Points earned"
	case "ADJUST_POINTS":
		if event.AdjustPoints != nil {
			pointsValue := event.AdjustPoints.Points
			if pointsValue > 0 {
				transactionType = "earn"
				points = pointsValue
				description = "Points adjustment (earned)"
			} else {
				transactionType = "redeem"
				points = -pointsValue // Make positive for display
				description = "Points redeemed"
			}
		}
	case "CREATE_REWARD":
		transactionType = "redeem"
		if event.CreateReward != nil {
			points = event.CreateReward.Points
		}
		description = "Reward created"
	default:
		return nil // Skip unknown event types
	}

	// Parse created date
	createdAt := time.Now()
	if event.CreatedAt != "" {
		if parsedTime, err := time.Parse(time.RFC3339, event.CreatedAt); err == nil {
			createdAt = parsedTime
		}
	}

	return &models.Transaction{
		ID:          event.ID,
		UserID:      userID,
		Type:        transactionType,
		Points:      points,
		Description: description,
		CreatedAt:   createdAt,
	}
}

func (s *LoyaltyService) generateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
