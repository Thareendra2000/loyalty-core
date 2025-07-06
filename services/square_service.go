package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"loyalty-core/config"

	square "github.com/square/square-go-sdk"
	"github.com/square/square-go-sdk/client"
	"github.com/square/square-go-sdk/core"
	"github.com/square/square-go-sdk/loyalty"
	"github.com/square/square-go-sdk/option"
)

// SquareService handles Square API interactions
type SquareService struct {
	config     *config.Config
	client     *client.Client
	programID  string
	locationID string
}

// NewSquareService creates a new Square service instance
func NewSquareService(cfg *config.Config) (*SquareService, error) {
	// Validate required configuration
	if cfg.SquareAccessToken == "" || cfg.SquareAccessToken == "your-square-access-token" {
		return nil, errors.New("Square access token is required")
	}

	if cfg.SquareLocationID == "" || cfg.SquareLocationID == "your-square-location-id" {
		return nil, errors.New("Square location ID is required")
	}

	// Initialize Square client
	squareClient := client.NewClient(
		option.WithToken(cfg.SquareAccessToken),
		option.WithHTTPClient(getHTTPClient(cfg.SquareEnvironment)),
	)

	service := &SquareService{
		config:     cfg,
		client:     squareClient,
		locationID: cfg.SquareLocationID,
	}

	// Get the loyalty program ID
	programID, err := service.getLoyaltyProgramID()
	if err != nil {
		return nil, fmt.Errorf("failed to get loyalty program ID: %w", err)
	}
	service.programID = programID

	return service, nil
}

// getLoyaltyProgramID retrieves the main loyalty program ID
func (s *SquareService) getLoyaltyProgramID() (string, error) {
	ctx := context.Background()

	// Use "main" as the program ID to get the default loyalty program
	request := &loyalty.GetProgramsRequest{
		ProgramID: "main",
	}

	response, err := s.client.Loyalty.Programs.Get(ctx, request)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve loyalty program: %w", err)
	}

	if response.Program == nil {
		return "", errors.New("no loyalty program found")
	}

	return *response.Program.ID, nil
}

// CreateLoyaltyAccount creates a loyalty account in Square
func (s *SquareService) CreateLoyaltyAccount(phoneNumber, givenName, familyName string) (*square.LoyaltyAccount, error) {
	ctx := context.Background()

	// Generate idempotency key
	idempotencyKey := fmt.Sprintf("create-loyalty-%d", time.Now().UnixNano())

	// Create the loyalty account request
	request := &loyalty.CreateLoyaltyAccountRequest{
		LoyaltyAccount: &square.LoyaltyAccount{
			ProgramID: s.programID,
			Mapping: &square.LoyaltyAccountMapping{
				PhoneNumber: &phoneNumber,
			},
		},
		IdempotencyKey: idempotencyKey,
	}

	response, err := s.client.Loyalty.Accounts.Create(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to create loyalty account: %w", err)
	}

	if response.LoyaltyAccount == nil {
		return nil, errors.New("no loyalty account returned from Square")
	}

	return response.LoyaltyAccount, nil
}

// AccumulateLoyaltyPoints adds points to a loyalty account
func (s *SquareService) AccumulateLoyaltyPoints(accountID string, points int, orderID string) (*square.LoyaltyEvent, error) {
	ctx := context.Background()

	// Generate idempotency key
	idempotencyKey := fmt.Sprintf("accumulate-points-%s-%d", accountID, time.Now().UnixNano())

	request := &loyalty.AccumulateLoyaltyPointsRequest{
		AccountID: accountID,
		AccumulatePoints: &square.LoyaltyEventAccumulatePoints{
			LoyaltyProgramID: &s.programID,
			Points:           &points,
			OrderID:          &orderID,
		},
		IdempotencyKey: idempotencyKey,
		LocationID:     s.locationID,
	}

	response, err := s.client.Loyalty.Accounts.AccumulatePoints(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to accumulate loyalty points: %w", err)
	}

	if response.Event == nil {
		return nil, errors.New("no loyalty event returned from Square")
	}

	return response.Event, nil
}

// CreateLoyaltyReward creates a loyalty reward (redeems points)
func (s *SquareService) CreateLoyaltyReward(accountID string, rewardTierID string, orderID string) (*square.LoyaltyReward, error) {
	ctx := context.Background()

	// Generate idempotency key
	idempotencyKey := fmt.Sprintf("create-reward-%s-%d", accountID, time.Now().UnixNano())

	request := &loyalty.CreateLoyaltyRewardRequest{
		Reward: &square.LoyaltyReward{
			LoyaltyAccountID: accountID,
			RewardTierID:     rewardTierID,
			OrderID:          &orderID,
		},
		IdempotencyKey: idempotencyKey,
	}

	response, err := s.client.Loyalty.Rewards.Create(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to create loyalty reward: %w", err)
	}

	if response.Reward == nil {
		return nil, errors.New("no loyalty reward returned from Square")
	}

	return response.Reward, nil
}

// GetLoyaltyAccount retrieves a loyalty account by ID
func (s *SquareService) GetLoyaltyAccount(accountID string) (*square.LoyaltyAccount, error) {
	ctx := context.Background()

	request := &loyalty.GetAccountsRequest{
		AccountID: accountID,
	}

	response, err := s.client.Loyalty.Accounts.Get(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve loyalty account: %w", err)
	}

	if response.LoyaltyAccount == nil {
		return nil, errors.New("no loyalty account found")
	}

	return response.LoyaltyAccount, nil
}

// SearchLoyaltyAccounts searches for loyalty accounts by phone number
func (s *SquareService) SearchLoyaltyAccounts(phoneNumber string) ([]*square.LoyaltyAccount, error) {
	ctx := context.Background()

	request := &loyalty.SearchLoyaltyAccountsRequest{
		Query: &square.SearchLoyaltyAccountsRequestLoyaltyAccountQuery{
			Mappings: []*square.LoyaltyAccountMapping{
				{
					PhoneNumber: &phoneNumber,
				},
			},
		},
	}

	response, err := s.client.Loyalty.Accounts.Search(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to search loyalty accounts: %w", err)
	}

	return response.LoyaltyAccounts, nil
}

// SearchLoyaltyEvents searches for loyalty events (transaction history)
func (s *SquareService) SearchLoyaltyEvents(accountID string, limit int) ([]*square.LoyaltyEvent, error) {
	ctx := context.Background()

	request := &square.SearchLoyaltyEventsRequest{
		Query: &square.LoyaltyEventQuery{
			Filter: &square.LoyaltyEventFilter{
				LoyaltyAccountFilter: &square.LoyaltyEventLoyaltyAccountFilter{
					LoyaltyAccountID: accountID,
				},
			},
		},
		Limit: &limit,
	}

	response, err := s.client.Loyalty.SearchEvents(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to search loyalty events: %w", err)
	}

	return response.Events, nil
}

// AdjustLoyaltyPoints adjusts points in a loyalty account (for manual point redemption)
func (s *SquareService) AdjustLoyaltyPoints(accountID string, points int, reason string) (*square.LoyaltyEvent, error) {
	ctx := context.Background()

	// Generate idempotency key
	idempotencyKey := fmt.Sprintf("adjust-points-%s-%d", accountID, time.Now().UnixNano())

	request := &loyalty.AdjustLoyaltyPointsRequest{
		AccountID: accountID,
		AdjustPoints: &square.LoyaltyEventAdjustPoints{
			LoyaltyProgramID: &s.programID,
			Points:           points,
			Reason:           &reason,
		},
		IdempotencyKey: idempotencyKey,
	}

	response, err := s.client.Loyalty.Accounts.Adjust(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to adjust loyalty points: %w", err)
	}

	if response.Event == nil {
		return nil, errors.New("no loyalty event returned from Square")
	}

	return response.Event, nil
}

// getHTTPClient returns the appropriate HTTP client based on environment
func getHTTPClient(_ string) core.HTTPClient {
	// Return default HTTP client
	// In production, you might want to configure timeout, retry logic, etc.
	return nil // This will use the default HTTP client
}
