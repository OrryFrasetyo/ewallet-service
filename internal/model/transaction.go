package model

import "time"

type Transaction struct {
	ID              int       `json:"id"`
	WalletID        int       `json:"wallet_id"`
	TransactionType string    `json:"transaction_type"`
	Amount          float64   `json:"amount"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
}

type TopUpRequest struct {
	Amount float64 `json:"amount" binding:"required,min=10000"`
}

type TopUpResponse struct {
	ID            int       `json:"id"`
	WalletNumber  string    `json:"wallet_number"`
	BalanceBefore float64   `json:"balance_before"`
	BalanceAfter  float64   `json:"balance_after"`
	TopUpAmount   float64   `json:"topup_amount"`
	CreatedAt     time.Time `json:"created_at"`
}
