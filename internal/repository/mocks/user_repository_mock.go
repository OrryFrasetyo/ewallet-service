package mocks

import (
	"context"
	"ewallet-service/internal/model"

	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) EmailExists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *UserRepositoryMock) RegisterUser(ctx context.Context, user *model.User) (model.Wallet, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(model.Wallet), args.Error(1)
}

func (m *UserRepositoryMock) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *UserRepositoryMock) FindWalletByUserID(ctx context.Context, userID int) (*model.Wallet, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Wallet), args.Error(1)
}
