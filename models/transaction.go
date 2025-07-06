package models

import (
	"time"
)

type Transaction struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Type        string    `json:"type"` // "earn" or "redeem"
	Points      int       `json:"points"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

type EarnRequest struct {
	Points      int    `json:"points" binding:"required"`
	Description string `json:"description"`
}

type RedeemRequest struct {
	Points      int    `json:"points" binding:"required"`
	Description string `json:"description"`
}

type BalanceResponse struct {
	Points       int           `json:"points"`
	Transactions []Transaction `json:"transactions"`
}
