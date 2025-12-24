package usecase

import (
	"context"
	"ewallet-service/internal/model"
	"ewallet-service/internal/repository"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	UserRepo *repository.UserRepository
}

func NewUserUsecase(repo *repository.UserRepository) *UserUsecase {
	return &UserUsecase{UserRepo: repo}
}

func (u *UserUsecase) Register(ctx context.Context, req model.RegisterRequest) (model.RegisterResponse, error) {
	// check duplicate email
	emailExists, err := u.UserRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return model.RegisterResponse{}, err
	}
	if emailExists {
		return model.RegisterResponse{}, fmt.Errorf("Email %s sudah terdaftar", req.Email)
	}

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
