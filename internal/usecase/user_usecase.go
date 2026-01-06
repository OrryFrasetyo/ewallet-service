package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"ewallet-service/internal/model"
	"ewallet-service/internal/repository"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	UserRepo repository.UserRepository
	Redis    *redis.Client
}

func NewUserUsecase(repo repository.UserRepository, rdb *redis.Client) *UserUsecase {
	return &UserUsecase{
		UserRepo: repo,
		Redis:    rdb}
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

func (u *UserUsecase) Login(ctx context.Context, req model.LoginRequest) (model.LoginResponse, error) {
	// search user by email
	user, err := u.UserRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return model.LoginResponse{}, errors.New("Email atau password salah")
	}

	// check password (hash vs plain)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return model.LoginResponse{}, errors.New("Email atau password salah")
	}

	// generate jwt token
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	// token with method HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return model.LoginResponse{}, fmt.Errorf("Gagal generate token: %v", err)
	}

	return model.LoginResponse{
		AccessToken: signedToken,
		Type:        "Bearer",
		ExpiresIn:   "24h",
	}, nil
}

func (u *UserUsecase) GetBalance(ctx context.Context, userID int) (*model.Wallet, error) {
	// key cache
	cacheKey := fmt.Sprintf("wallet:%d", userID)

	// check redis (cache hit)
	if u.Redis != nil {
		val, err := u.Redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var cachedWallet model.Wallet
			err = json.Unmarshal([]byte(val), &cachedWallet)
			if err == nil {
				fmt.Println("⚡ Cache Hit! Saldo diambil dari Redis")
				return &cachedWallet, nil
			}
		}
	}

	// cache miss (take from database)
	fmt.Println("🐢 Cache Miss. Mengambil saldo dari Database...")
	wallet, err := u.UserRepo.FindWalletByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// save to redis (set cache)
	if u.Redis != nil {
		data, _ := json.Marshal(wallet)
		err = u.Redis.Set(ctx, cacheKey, data, 5*time.Minute).Err()
		if err != nil {
			fmt.Println("⚠️ Gagal menyimpan ke Redis:", err)
		}
	}
	return wallet, nil
}
