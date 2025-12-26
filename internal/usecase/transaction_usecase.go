package usecase

import (
	"context"
	"ewallet-service/internal/model"
	"ewallet-service/internal/repository"
)

type TransactionUsecase struct {
	TransactionRepo *repository.TransactionRepository
}

func NewTransactionUsecase(repo *repository.TransactionRepository) *TransactionUsecase {
	return &TransactionUsecase{TransactionRepo: repo}
}

func (u *TransactionUsecase) TopUp(ctx context.Context, userID int, req model.TopUpRequest) (model.TopUpResponse, error) {
	return u.TransactionRepo.CreateTopUp(ctx, userID, req.Amount)
}
