package model

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Wallet struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	Balance      float64   `json:"balance"`
	WalletNumber string    `json:"wallet_number"`
	CreatedAt    time.Time `json:"created_at"`
}

// DTO (Data Transfer Object) - what is sent to the client
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterResponse struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	WalletNumber string  `json:"wallet_number"`
	Balance      float64 `json:"balance"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	Type        string `json:"token_type"`
	ExpiresIn   string  `json:"expires_in"` 
}
