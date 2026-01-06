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
	db := config.ConnectDB()
	rdb := config.ConnectRedis()

	// DI User
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo, rdb)
	userHandler := handler.NewUserHandler(userUsecase)

	// DI Transaction
	trxRepo := repository.NewTransactionRepository(db)
	trxUsecase := usecase.NewTransactionUsecase(trxRepo, rdb)
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
