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
		c.JSON(http.StatusBadRequest, WebResponse{
			Status:  "fail",
			Message: "Input tidak valid",
			Error:   err.Error(),
		})
		return
	}

	// c.Request.Context() penting untuk meneruskan context (timeout/cancellation)
	res, err := h.UserUsecase.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusConflict, WebResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// response success
	c.JSON(http.StatusCreated, WebResponse{
		Status:  "success",
		Message: "User registered successfully",
		Data:    res,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, WebResponse{
			Status:  "fail",
			Message: "Input tidak valid",
			Error:   err.Error(),
		})
	}

	res, err := h.UserUsecase.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, WebResponse{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, WebResponse{
		Status:  "success",
		Message: "Login berhasil",
		Data:    res,
	})
}
