package usecase

import (
	"context"
	"ewallet-service/internal/model"
	"ewallet-service/internal/repository"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type TransactionUsecase struct {
	TransactionRepo repository.TransactionRepository
	Redis           *redis.Client
}

func NewTransactionUsecase(repo repository.TransactionRepository, rdb *redis.Client) *TransactionUsecase {
	return &TransactionUsecase{
		TransactionRepo: repo,
		Redis:           rdb,
	}
}

func (u *TransactionUsecase) TopUp(ctx context.Context, userID int, req model.TopUpRequest) (model.TopUpResponse, error) {
	res, err := u.TransactionRepo.CreateTopUp(ctx, userID, req.Amount)
	if err != nil {
		return model.TopUpResponse{}, err
	}

	// cache invalidation (delete cache balance user)
	// So that when the user checks the balance again, forced to take new data from the DB.
	if u.Redis != nil {
		cacheKey := fmt.Sprintf("wallet:%d", userID)
		u.Redis.Del(ctx, cacheKey)
		fmt.Printf("🧹 Cache Invalidated for User %d (After TopUp)\n", userID)
	}

	return res, nil
}

func (u *TransactionUsecase) Transfer(ctx context.Context, senderID int, req model.TransferRequest) (model.TransferResponse, error) {
	res, err := u.TransactionRepo.Transfer(ctx, senderID, req)
	if err != nil {
		return model.TransferResponse{}, err
	}

	// cache invalidation (delete cache sender)
	// As the sender's balance decreases, his old cache must be discarded.
	if u.Redis != nil {
		cacheKey := fmt.Sprintf("wallet:%d", senderID)
		u.Redis.Del(ctx, cacheKey)
		fmt.Printf("🧹 Cache Invalidated for User %d (After Transfer)\n", senderID)
	}

	return res, nil
}

func (u *TransactionUsecase) GetHistory(ctx context.Context, userID int) ([]model.Transaction, error) {
	return u.TransactionRepo.GetTransactionHistory(ctx, userID)
}
