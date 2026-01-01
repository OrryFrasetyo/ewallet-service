package repository

import (
	"context"
	"database/sql"
	"ewallet-service/internal/model"
	"fmt"
	"math/rand"
)

type UserRepository interface {
	RegisterUser(ctx context.Context, user *model.User) (model.Wallet, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindWalletByUserID(ctx context.Context, userID int) (*model.Wallet, error)
}

type userRepositoryPostgres struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepositoryPostgres{DB: db}
}

func (r *userRepositoryPostgres) RegisterUser(ctx context.Context, user *model.User) (model.Wallet, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return model.Wallet{}, err
	}

	defer tx.Rollback()

	sqlUser := "INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id"

	var userID int
	err = tx.QueryRowContext(ctx, sqlUser, user.Name, user.Email, user.Password).Scan(&userID)
	if err != nil {
		return model.Wallet{}, fmt.Errorf("Gagal insert user: %w", err)
	}
	user.ID = userID

	// generate wallet number (10 digits random)
	walletNumber := fmt.Sprintf("100%d", rand.Intn(9999999))

	// insert wallet (automatic balance 0)
	sqlWallet := "INSERT INTO wallets (user_id, wallet_number, balance) VALUES ($1, $2, 0) RETURNING id, balance, created_at"

	var wallet model.Wallet
	err = tx.QueryRowContext(ctx, sqlWallet, userID, walletNumber).Scan(&wallet.ID, &wallet.Balance, &wallet.CreatedAt)
	if err != nil {
		return model.Wallet{}, fmt.Errorf("Gagal insert wallet: %w", err)
	}

	wallet.UserID = userID
	wallet.WalletNumber = walletNumber

	if err := tx.Commit(); err != nil {
		return model.Wallet{}, err
	}

	return wallet, nil

}

func (r *userRepositoryPostgres) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)"
	err := r.DB.QueryRowContext(ctx, query, email).Scan(&exists)
	return exists, err
}

func (r *userRepositoryPostgres) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := "SELECT id, name, email, password FROM users WHERE email=$1"

	var user model.User

	err := r.DB.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepositoryPostgres) FindWalletByUserID(ctx context.Context, userID int) (*model.Wallet, error) {
	query := "SELECT id, user_id, balance, wallet_number, created_at FROM wallets WHERE user_id = $1"

	var w model.Wallet
	err := r.DB.QueryRowContext(ctx, query, userID).Scan(&w.ID, &w.UserID, &w.Balance, &w.WalletNumber, &w.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &w, nil
}
