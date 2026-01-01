package main

import (
	"ewallet-service/config"
	"ewallet-service/internal/handler"
	"ewallet-service/internal/middleware"
	"ewallet-service/internal/repository"
	"ewallet-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	// DI User
	userRepo := repository.NewUserRepository(config.DB)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)

	// DI Transaction
	trxRepo := repository.NewTransactionRepository(config.DB)
	trxUsecase := usecase.NewTransactionUsecase(trxRepo)
	trxHandler := handler.NewTransactionHandler(trxUsecase)

	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)

		protected := api.Group("/", middleware.AuthMiddleware())
		{
			protected.POST("/topup", trxHandler.TopUp)
			protected.POST("/transfer", trxHandler.Transfer)
			protected.GET("/transactions", trxHandler.HistoryTransaction)
			protected.GET("/balance", userHandler.GetBalance)

		}
	}

	r.Run(":8080")
}
