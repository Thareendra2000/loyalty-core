package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"loyalty-core/config"
	"loyalty-core/models"
	"loyalty-core/storage"
)

type LoyaltyService struct {
	config       *config.Config
	userStorage  *storage.UserStorage
	transactions map[string][]models.Transaction
	instanceID   string // Debug: track service instance
}

func NewLoyaltyService(cfg *config.Config) *LoyaltyService {
	service := &LoyaltyService{
		config:       cfg,
		userStorage:  storage.GetGlobalUserStorage(),
		transactions: make(map[string][]models.Transaction),
	}

	return service
}

func (s *LoyaltyService) EarnPoints(userID string, points int, description string) (*models.Transaction, error) {
	user, err := s.userStorage.GetUserByID(userID)
	if err != nil {
		return nil, err
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

	// Update user points
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

	if user.Points < points {
		return nil, errors.New("insufficient points")
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

	// Update user points
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

	transactions := s.transactions[userID]
	if transactions == nil {
		transactions = []models.Transaction{}
	}

	return &models.BalanceResponse{
		Points:       user.Points,
		Transactions: transactions,
	}, nil
}

func (s *LoyaltyService) GetTransactionHistory(userID string) ([]models.Transaction, error) {
	// Verify user exists
	_, err := s.userStorage.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	transactions := s.transactions[userID]
	if transactions == nil {
		return []models.Transaction{}, nil
	}

	return transactions, nil
}

func (s *LoyaltyService) generateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
