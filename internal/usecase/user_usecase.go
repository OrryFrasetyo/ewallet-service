package usecase

import (
	"context"
	"ewallet-service/internal/model"
	"ewallet-service/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	UserRepo *repository.UserRepository
}

func NewUserUsecase(repo *repository.UserRepository) *UserUsecase {
	return &UserUsecase{UserRepo: repo}
}

func (u *UserUsecase) Register(ctx context.Context, req model.RegisterRequest) (model.RegisterResponse, error) {
	// hash password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.RegisterResponse{}, err
	}

	newUser := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPass),
	}

	createdWallet, err := u.UserRepo.RegisterUser(ctx, newUser)
	if err != nil {
		return model.RegisterResponse{}, err
	}

	return model.RegisterResponse{
		ID:           newUser.ID,
		Name:         newUser.Name,
		Email:        newUser.Email,
		WalletNumber: createdWallet.WalletNumber,
		Balance:      createdWallet.Balance,
	}, nil
}
