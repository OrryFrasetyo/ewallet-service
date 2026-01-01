package usecase

import (
	"context"
	"ewallet-service/internal/model"
	"ewallet-service/internal/repository"
)

type TransactionUsecase struct {
	TransactionRepo repository.TransactionRepository
}

func NewTransactionUsecase(repo repository.TransactionRepository) *TransactionUsecase {
	return &TransactionUsecase{TransactionRepo: repo}
}

func (u *TransactionUsecase) TopUp(ctx context.Context, userID int, req model.TopUpRequest) (model.TopUpResponse, error) {
	return u.TransactionRepo.CreateTopUp(ctx, userID, req.Amount)
}

func (u *TransactionUsecase) Transfer(ctx context.Context, senderID int, req model.TransferRequest) (model.TransferResponse, error) {
	return u.TransactionRepo.Transfer(ctx, senderID, req)
}

func (u *TransactionUsecase) GetHistory(ctx context.Context, userID int) ([]model.Transaction, error) {
	return u.TransactionRepo.GetTransactionHistory(ctx, userID)
}
