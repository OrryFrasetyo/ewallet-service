package main

import (
	"context"
	"ewallet-service/config"
	"ewallet-service/internal/handler"
	"ewallet-service/internal/middleware"
	"ewallet-service/internal/repository"
	"ewallet-service/internal/service"
	"ewallet-service/internal/usecase"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.ConnectDB()
	rdb := config.ConnectRedis()

	aiService, err := service.NewAIService(context.Background())
	if err != nil {
		log.Fatal("Gagal init AI Service: ", err)
	}

	// DI User
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo, rdb)
	userHandler := handler.NewUserHandler(userUsecase)

	// DI Transaction
	trxRepo := repository.NewTransactionRepository(db)
	trxUsecase := usecase.NewTransactionUsecase(trxRepo, rdb, aiService)
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
			protected.GET("/financial-health", trxHandler.GetFinancialHealth)
			protected.GET("/balance", userHandler.GetBalance)

		}
	}

	r.Run(":8080")
}
