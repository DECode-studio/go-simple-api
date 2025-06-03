package middleware

import (
	"fmt"
	model_api "go-simple-api/src/model/api"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model_api.ApiResponseDoc{
			Status: &model_api.StatusModel{Code: http.StatusUnauthorized, Message: "Authorization header is required"},
		})
		return
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model_api.ApiResponseDoc{
			Status: &model_api.StatusModel{Code: http.StatusUnauthorized, Message: "Invalid authorization format. Use 'Bearer <token>'"},
		})
		return
	}
	tokenString := headerParts[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {

		if err == jwt.ErrTokenExpired {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model_api.ApiResponseDoc{
				Status: &model_api.StatusModel{Code: http.StatusUnauthorized, Message: "Token has expired"},
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, model_api.ApiResponseDoc{
			Status: &model_api.StatusModel{Code: http.StatusUnauthorized, Message: "Invalid token"},
		})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		idUser, userOk := claims["id"].(string)
		if !userOk {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model_api.ApiResponseDoc{
				Status: &model_api.StatusModel{Code: http.StatusUnauthorized, Message: "Invalid token claims"},
			})
			return
		}

		c.Set("idUser", idUser)

		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model_api.ApiResponseDoc{
			Status: &model_api.StatusModel{Code: http.StatusUnauthorized, Message: "Invalid token"},
		})
	}
}
