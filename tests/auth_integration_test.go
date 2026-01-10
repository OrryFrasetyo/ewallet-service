package tests

import (
	"bytes"
	"encoding/json"
	"ewallet-service/config"
	"ewallet-service/internal/handler"
	"ewallet-service/internal/repository"
	"ewallet-service/internal/usecase"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5435")
	os.Setenv("DB_USER", "myuser")
	os.Setenv("DB_PASSWORD", "mypassword")
	os.Setenv("DB_NAME", "ewalletdb")

	db := config.ConnectDB()

	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo, nil)
	userHandler := handler.NewUserHandler(userUsecase)

	r := gin.Default()
	api := r.Group("/api/v1")
	{
		api.POST("/register", userHandler.Register)
	}
	return r
}

func TestIntegration_Register_Success(t *testing.T) {
	router := setupRouter()

	requestBody := map[string]string{
		"name":     "Integration Test User",
		"email":    "integration@test.com",
		"password": "password123",
	}
	jsonValue, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])

	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5435")
	os.Setenv("DB_USER", "myuser")
	os.Setenv("DB_PASSWORD", "mypassword")
	os.Setenv("DB_NAME", "ewalletdb")

	db := config.ConnectDB()
	db.Exec("DELETE FROM wallets WHERE user_id IN (SELECT id FROM users WHERE email = 'integration@test.com')")
	db.Exec("DELETE FROM users WHERE email = 'integration@test.com'")
}
