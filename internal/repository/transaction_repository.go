package repository

import (
	"context"
	"database/sql"
	"errors"
	"ewallet-service/internal/model"
	"fmt"
	"time"
)

type TransactionRepository struct {
	DB *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{DB: db}
}

func (r *TransactionRepository) CreateTopUp(ctx context.Context, userID int, amount float64) (model.TopUpResponse, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return model.TopUpResponse{}, err
	}
	defer tx.Rollback()

	var walletID int
	var walletNumber string
	var currentBalance float64

	queryCheck := "SELECT id, wallet_number, balance FROM wallets WHERE user_id = $1 FOR UPDATE"
	err = tx.QueryRowContext(ctx, queryCheck, userID).Scan(&walletID, &walletNumber, &currentBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.TopUpResponse{}, errors.New("Wallet tidak ditemukan")
		}
		return model.TopUpResponse{}, err
	}

	newBalance := currentBalance + amount
	queryUpdate := "UPDATE wallets SET balance = $1, updated_at = NOW() WHERE id = $2"
	_, err = tx.ExecContext(ctx, queryUpdate, newBalance, walletID)
	if err != nil {
		return model.TopUpResponse{}, fmt.Errorf("Gagal update saldo: %w", err)
	}

	var transactionID int
	queryHistory := `
		INSERT INTO transactions (wallet_id, transaction_type, amount, description, created_at) VALUES ($1, 'TOPUP', $2, 'Topup Saldo via API', NOW())
		RETURNING id
	`

	err = tx.QueryRowContext(ctx, queryHistory, walletID, amount).Scan(&transactionID)
	if err != nil {
		return model.TopUpResponse{}, fmt.Errorf("Gagal catat history: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return model.TopUpResponse{}, err
	}

	return model.TopUpResponse{
		ID: transactionID,
		WalletNumber: walletNumber,
		BalanceBefore: currentBalance,
		BalanceAfter: newBalance,
		TopUpAmount: amount,
		CreatedAt: time.Now(),
	}, nil
}
