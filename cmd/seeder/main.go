package main

import (
	"context"
	"ewallet-service/config"
	"ewallet-service/internal/model"
	"ewallet-service/internal/repository"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	config.ConnectDB()

	repo := repository.NewUserRepository(config.DB)

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("rahasia123"), bcrypt.DefaultCost)

	dummyUser := &model.User{
		Name: "Seeder User",
		Email: "seeder@example.com",
		Password: string(passwordHash),
	}

	wallet, err := repo.RegisterUser(context.Background(), dummyUser)
	if err != nil {
		log.Fatalf("Gagal seeding: %v", err)
	}

	fmt.Printf("🌱 Seeding success!\nUser: %s\nWallet: %s\nBalance: %.2f\n", dummyUser.Email, wallet.WalletNumber, wallet.Balance)
}