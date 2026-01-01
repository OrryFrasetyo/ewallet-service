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

type TransferRequest struct {
	TargetWalletNumber string  `json:"target_wallet_number" binding:"required"`
	Amount             float64 `json:"amount" binding:"required,min=1000"`
	Description        string  `json:"description"`
}

type TransferResponse struct {
	ID             string    `json:"id"`
	SenderBalance  float64   `json:"sender_balance"`
	ReceiverWallet string    `json:"receiver_wallet"`
	Amount         float64   `json:"amount"`
	CreatedAt      time.Time `json:"created_at"`
}
