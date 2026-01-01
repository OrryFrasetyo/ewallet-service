package usecase_test

import (
	"context"
	"errors"
	"ewallet-service/internal/model"
	"ewallet-service/internal/repository/mocks"
	"ewallet-service/internal/usecase"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTopUp_Success(t *testing.T) {
	// arrange
	mockRepo := new(mocks.TransactionRepositoryMock)
	u := usecase.NewTransactionUsecase(mockRepo)

	userID := 1
	req := model.TopUpRequest{
		Amount: 50000,
	}

	expectedRes := model.TopUpResponse{
		ID:           1,
		TopUpAmount:  50000,
		BalanceAfter: 50000,
		WalletNumber: "1001",
		CreatedAt:    time.Now(),
	}

	mockRepo.On("CreateTopUp", mock.Anything, userID, req.Amount).Return(expectedRes, nil)

	// act
	res, err := u.TopUp(context.Background(), userID, req)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, float64(50000), res.TopUpAmount)
	mockRepo.AssertExpectations(t)
}

func TestTransfer_Success(t *testing.T) {
	// arrange
	mockRepo := new(mocks.TransactionRepositoryMock)
	u := usecase.NewTransactionUsecase(mockRepo)

	senderID := 1
	req := model.TransferRequest{
		TargetWalletNumber: "100999",
		Amount:             25000,
	}

	expectedRes := model.TransferResponse{
		ID:             "TRX-123",
		Amount:         25000,
		ReceiverWallet: "100999",
	}

	mockRepo.On("Transfer", mock.Anything, senderID, req).Return(expectedRes, nil)

	// act
	res, err := u.Transfer(context.Background(), senderID, req)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, "100999", res.ReceiverWallet)
	mockRepo.AssertExpectations(t)
}

func TestTransfer_Failed_RepoError(t *testing.T) {
	// arrange
	mockRepo := new(mocks.TransactionRepositoryMock)
	u := usecase.NewTransactionUsecase(mockRepo)

	senderID := 1
	req := model.TransferRequest{
		TargetWalletNumber: "100999",
		Amount:             1000000,
	}

	expectedErr := errors.New("Saldo tidak mencukupi")
	mockRepo.On("Transfer", mock.Anything, senderID, req).Return(model.TransferResponse{}, expectedErr)

	// act
	res, err := u.Transfer(context.Background(), senderID, req)

	// assert
	assert.Error(t, err)
	assert.Equal(t, "Saldo tidak mencukupi", err.Error())
	assert.Empty(t, res.ID)

	mockRepo.AssertExpectations(t)
}

func TestGetHistory_Success(t *testing.T) {
	// arrange
	mockRepo := new(mocks.TransactionRepositoryMock)
	u := usecase.NewTransactionUsecase(mockRepo)

	userID := 1

	// data dummy
	expectedHistory := []model.Transaction{
		{
			ID:              1,
			Amount:          50000,
			TransactionType: "TOPUP",
			Description:     "Topup Saldo",
		},
		{
			ID:              2,
			Amount:          20000,
			TransactionType: "TRANSFER_OUT",
			Description:     "Bayar hutang",
		},
	}

	// mocking
	mockRepo.On("GetTransactionHistory", mock.Anything, userID).Return(expectedHistory, nil)

	// act
	res, err := u.GetHistory(context.Background(), userID)

	// assert
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, "TOPUP", res[0].TransactionType)
	assert.Equal(t, float64(20000), res[1].Amount)

	mockRepo.AssertExpectations(t)
}

func TestGetHistory_Empty(t *testing.T) {
	// arrange
	mockRepo := new(mocks.TransactionRepositoryMock)
	u := usecase.NewTransactionUsecase(mockRepo)

	userID := 2
	expectedHistory := []model.Transaction{}

	mockRepo.On("GetTransactionHistory", mock.Anything, userID).Return(expectedHistory, nil)

	// act
	res, err := u.GetHistory(context.Background(), userID)

	// assert
	assert.NoError(t, err)
	assert.Len(t, res, 0)
	assert.NotNil(t, res)
}
