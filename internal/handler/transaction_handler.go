package handler

import (
	"ewallet-service/internal/model"
	"ewallet-service/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	TransactionUsecase *usecase.TransactionUsecase
}

func NewTransactionHandler(u *usecase.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{TransactionUsecase: u}
}

func (h *TransactionHandler) TopUp(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, WebResponse{
			Status:  "fail",
			Message: "User ID tidak ditemukan di token",
		})
		return
	}

	var req model.TopUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, WebResponse{
			Status:  "fail",
			Message: "Input tidak valid (min: 10000)",
			Error:   err.Error(),
		})
		return
	}

	res, err := h.TransactionUsecase.TopUp(c.Request.Context(), userID.(int), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, WebResponse{
			Status: "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, WebResponse{
		Status: "success",
		Message: "Topup berhasil",
		Data: res,
	})
}
