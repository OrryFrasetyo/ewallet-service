package middleware

import (
	"ewallet-service/internal/handler"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, handler.WebResponse{
				Status:  "fail",
				Message: "Token tidak ditemukan atau format salah",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		secretKey := []byte(os.Getenv("JWT_SECRET"))

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, gin.Error{Err: jwt.ErrSignatureInvalid}
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, handler.WebResponse{
				Status:  "fail",
				Message: "Token tidak valid atau kadaluarsa",
			})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := int(claims["user_id"].(float64))

			c.Set("userID", userID)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, handler.WebResponse{
				Status:  "fail",
				Message: "Gagal memproses token claims",
			})
			return 
		}

		c.Next()
	}
}
