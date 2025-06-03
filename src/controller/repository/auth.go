package repository

import (
	"go-simple-api/src/config"
	model_api "go-simple-api/src/model/api"
	model_database "go-simple-api/src/model/database"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SignIn(c *gin.Context) {
	var input model_api.SignInRequest
	var user model_database.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, model_api.ApiResponseDoc{
			Status: &model_api.StatusModel{Code: http.StatusBadRequest, Message: err.Error()},
		})
		return
	}

	queryCondition := model_database.User{EmailUser: input.Email}

	if err := config.DB.Where(&queryCondition).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, model_api.ApiResponseDoc{
				Status: &model_api.StatusModel{Code: http.StatusUnauthorized, Message: "Invalid email or password"},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, model_api.ApiResponseDoc{
			Status: &model_api.StatusModel{Code: http.StatusInternalServerError, Message: "Database error"},
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordUser), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, model_api.ApiResponseDoc{
			Status: &model_api.StatusModel{Code: http.StatusUnauthorized, Message: "Invalid email or password"},
		})
		return
	}

	claims := jwt.MapClaims{
		"id":    user.IDUser,
		"name":  user.NameUser,
		"email": user.EmailUser,
		"exp":   time.Now().Add(time.Hour * 8).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model_api.ApiResponseDoc{
			Status: &model_api.StatusModel{Code: http.StatusInternalServerError, Message: "Failed to generate token"},
		})
		return
	}

	authResponse := model_api.AuthResponse{Token: tokenString}
	response := model_api.ApiResponse[model_api.AuthResponse]{
		Status: &model_api.StatusModel{Code: http.StatusOK, Message: "Sign in successful"},
		Data:   &authResponse,
	}

	c.JSON(http.StatusOK, response)
}
