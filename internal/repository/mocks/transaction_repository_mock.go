package mocks

import (
	"context"
	"ewallet-service/internal/model"

	"github.com/stretchr/testify/mock"
)

type TransactionRepositoryMock struct {
	mock.Mock
}

func (m *TransactionRepositoryMock) CreateTopUp(ctx context.Context, userID int, amount float64) (model.TopUpResponse, error) {
	args := m.Called(ctx, userID, amount)
	return args.Get(0).(model.TopUpResponse), args.Error(1)
}

func (m *TransactionRepositoryMock) Transfer(ctx context.Context, senderID int, req model.TransferRequest) (model.TransferResponse, error) {
	args := m.Called(ctx, senderID, req)
	return args.Get(0).(model.TransferResponse), args.Error(1)
}

func (m *TransactionRepositoryMock) GetTransactionHistory(ctx context.Context, userID int) ([]model.Transaction, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.Transaction), args.Error(1)
}
