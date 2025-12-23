package handler

import (
	"ewallet-service/internal/model"
	"ewallet-service/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserUsecase *usecase.UserUsecase
}

func NewUserHandler(u *usecase.UserUsecase) *UserHandler {
	return &UserHandler{UserUsecase: u}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req model.RegisterRequest

	// validation input json
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// c.Request.Context() penting untuk meneruskan context (timeout/cancellation)
	res, err := h.UserUsecase.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// response success
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "User registered successfully",
		"data":    res,
	})
}
