package main

import (
	"ewallet-service/config"
	"ewallet-service/internal/handler"
	"ewallet-service/internal/repository"
	"ewallet-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	// setup layers (dependency injection)
	userRepo := repository.NewUserRepository(config.DB)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)

	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.POST("/register", userHandler.Register)
	}

	r.Run(":8080")
}
