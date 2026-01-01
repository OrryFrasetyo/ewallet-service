package usecase_test

import (
	"context"
	"errors"
	"ewallet-service/internal/model"
	"ewallet-service/internal/repository/mocks"
	"ewallet-service/internal/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister_Success(t *testing.T) {
	// arrange
	mockRepo := new(mocks.UserRepositoryMock)

	u := usecase.NewUserUsecase(mockRepo)

	// data dummy
	req := model.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	expectedWallet := model.Wallet{
		WalletNumber: "10012345",
		Balance:      0,
	}

	mockRepo.On("EmailExists", mock.Anything, req.Email).Return(false, nil)

	mockRepo.On("RegisterUser", mock.Anything, mock.AnythingOfType("*model.User")).Return(expectedWallet, nil)

	// action
	res, err := u.Register(context.Background(), req)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, req.Name, res.Name)
	assert.Equal(t, "10012345", res.WalletNumber)

	mockRepo.AssertExpectations(t)
}

func TestRegister_EmailDuplicate(t *testing.T) {
	// arrange
	mockRepo := new(mocks.UserRepositoryMock)
	u := usecase.NewUserUsecase(mockRepo)

	req := model.RegisterRequest{
		Name:     "Duplikat",
		Email:    "ada@example.com",
		Password: "ada123",
	}

	mockRepo.On("EmailExists", mock.Anything, req.Email).Return(true, nil)

	// act
	res, err := u.Register(context.Background(), req)

	// assert
	assert.Error(t, err)
	assert.Equal(t, "Email ada@example.com sudah terdaftar", err.Error())
	assert.Empty(t, res.ID)
}

func TestLogin_Success(t *testing.T) {
	// arrange
	mockRepo := new(mocks.UserRepositoryMock)
	u := usecase.NewUserUsecase(mockRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	dummyUser := &model.User{
		ID:       1,
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}

	req := model.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(dummyUser, nil)

	// act
	res, err := u.Login(context.Background(), req)

	assert.NoError(t, err)
	assert.NotEmpty(t, res.AccessToken)
	assert.Equal(t, "Bearer", res.Type)

}

func TestLogin_WrongPassword(t *testing.T) {
	// arrange
	mockRepo := new(mocks.UserRepositoryMock)
	u := usecase.NewUserUsecase(mockRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("passwordBenar"), bcrypt.DefaultCost)

	dummyUser := &model.User{
		ID:       1,
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}

	req := model.LoginRequest{
		Email:    "test@example.com",
		Password: "passwordSalah",
	}

	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(dummyUser, nil)

	// act
	res, err := u.Login(context.Background(), req)

	// assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Email atau password salah")
	assert.Empty(t, res.AccessToken)
}

func TestGetBalance_Success(t *testing.T)  {
	// arrange
	mockRepo := new(mocks.UserRepositoryMock)
	u := usecase.NewUserUsecase(mockRepo)

	userID := 1
	expectedWallet := &model.Wallet{
		ID: 1,
		UserID: userID,
		Balance: 150000,
		WalletNumber: "100888",
	}

	// mock
	mockRepo.On("FindWalletByUserID", mock.Anything, userID).Return(expectedWallet, nil)

	// act
	res, err := u.GetBalance(context.Background(), userID)

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, float64(150000), res.Balance)
	assert.Equal(t, "100888", res.WalletNumber)

	mockRepo.AssertExpectations(t)
}

func TestGetBalance_Error(t *testing.T)  {
	mockRepo := new(mocks.UserRepositoryMock)
	u := usecase.NewUserUsecase(mockRepo)

	userID := 99
	expectedErr := errors.New("Database connection failed")

	// mocking
	mockRepo.On("FindWalletByUserID", mock.Anything, userID).Return(nil, expectedErr)

	// act
	res, err := u.GetBalance(context.Background(), userID)

	// assert
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "Database connection failed", err.Error())

	mockRepo.AssertExpectations(t)
}
