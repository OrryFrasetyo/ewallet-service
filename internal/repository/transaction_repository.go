package repository

import (
	"context"
	"database/sql"
	"errors"
	"ewallet-service/internal/model"
	"fmt"
	"time"
)

type TransactionRepository interface {
	CreateTopUp(ctx context.Context, userID int, amount float64) (model.TopUpResponse, error)
	Transfer(ctx context.Context, senderID int, req model.TransferRequest) (model.TransferResponse, error)
	GetTransactionHistory(ctx context.Context, userID int) ([]model.Transaction, error)
}

type transactionRepositoryPostgres struct {
	DB *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepositoryPostgres{DB: db}
}

func (r *transactionRepositoryPostgres) CreateTopUp(ctx context.Context, userID int, amount float64) (model.TopUpResponse, error) {
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
		ID:            transactionID,
		WalletNumber:  walletNumber,
		BalanceBefore: currentBalance,
		BalanceAfter:  newBalance,
		TopUpAmount:   amount,
		CreatedAt:     time.Now(),
	}, nil
}

func (r *transactionRepositoryPostgres) Transfer(ctx context.Context, senderID int, req model.TransferRequest) (model.TransferResponse, error) {
	// start transaction
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return model.TransferResponse{}, err
	}
	defer tx.Rollback()

	// check sender wallet & saldo (locking)
	var senderWalletID int
	var senderBalance float64

	querySender := "SELECT id, balance FROM wallets WHERE user_id = $1 FOR UPDATE"
	err = tx.QueryRowContext(ctx, querySender, senderID).Scan(&senderWalletID, &senderBalance)
	if err != nil {
		return model.TransferResponse{}, errors.New("Wallet Pengirim tidak ditemukan")
	}

	// check the balance enough?
	if senderBalance < req.Amount {
		return model.TransferResponse{}, errors.New("Saldo tidak mencukupi")
	}

	// check receiver wallet (locking)
	var receiverWalletID int
	var receiverUserID int

	queryReceiver := "SELECT id, user_id FROM wallets WHERE wallet_number = $1 FOR UPDATE"
	err = tx.QueryRowContext(ctx, queryReceiver, req.TargetWalletNumber).Scan(&receiverWalletID, &receiverUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.TransferResponse{}, errors.New("Nomor wallet tujuan tidak ditemukan")
		}
		return model.TransferResponse{}, err
	}

	// validation : Don't transfer it to yourself
	if senderWalletID == receiverWalletID {
		return model.TransferResponse{}, errors.New("Tidak bisa transfer ke wallet sendiri")
	}

	// update sender balance (decrease)
	_, err = tx.ExecContext(ctx, "UPDATE wallets SET balance = balance - $1, updated_at = NOW() WHERE id = $2", req.Amount, senderWalletID)
	if err != nil {
		return model.TransferResponse{}, fmt.Errorf("Gagal potong saldo: %w", err)
	}

	// update receiver balance (increase)
	_, err = tx.ExecContext(ctx, "UPDATE wallets SET balance = balance + $1, updated_at = NOW() WHERE id = $2", req.Amount, receiverWalletID)
	if err != nil {
		return model.TransferResponse{}, fmt.Errorf("Gagal tambah saldo: %w", err)
	}

	// record history (double entry)
	// record for sender (money out)
	var createdAt time.Time
	queryHistoryOut := `
		INSERT INTO transactions (wallet_id, transaction_type, amount, description, created_at) VALUES ($1, 'TRANSFER_OUT', $2, $3, NOW()) RETURNING created_at
	`

	err = tx.QueryRowContext(ctx, queryHistoryOut, senderWalletID, req.Amount, "Transfer ke "+req.TargetWalletNumber).Scan(&createdAt)
	if err != nil {
		return model.TransferResponse{}, fmt.Errorf("Gagal catat history pengirim: %w", err)
	}

	// record for receiver (money in)
	queryHistoryIn := `
		INSERT INTO transactions (wallet_id, transaction_type, amount, description, created_at) VALUES ($1, 'TRANSFER_IN', $2, $3, NOW())
	`

	_, err = tx.ExecContext(ctx, queryHistoryIn, receiverWalletID, req.Amount, "Terima transfer")
	if err != nil {
		return model.TransferResponse{}, fmt.Errorf("Gagal catat history penerima: %w", err)
	}

	// commit all
	if err := tx.Commit(); err != nil {
		return model.TransferResponse{}, err
	}

	return model.TransferResponse{
		ID:             fmt.Sprintf("TRX-%d-%d", senderWalletID, time.Now().Unix()),
		SenderBalance:  senderBalance - req.Amount,
		ReceiverWallet: req.TargetWalletNumber,
		Amount:         req.Amount,
		CreatedAt:      createdAt,
	}, nil
}

func (r *transactionRepositoryPostgres) GetTransactionHistory(ctx context.Context, userID int) ([]model.Transaction, error) {
	query := `
		SELECT t.id, t.wallet_id, t.transaction_type, t.amount, t.description, t.created_at
		FROM transactions t
		JOIN wallets w ON t.wallet_id = w.id
		WHERE w.user_id = $1
		ORDER BY t.created_at DESC
		LIMIT 10
	`

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var t model.Transaction
		if err := rows.Scan(&t.ID, &t.WalletID, &t.TransactionType, &t.Amount, &t.Description, &t.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	if transactions == nil {
		transactions = []model.Transaction{}
	}

	return transactions, nil
}
