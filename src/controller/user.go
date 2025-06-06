package api_controller

import (
	"go-simple-api/src/config"
	model_api "go-simple-api/src/model/api"
	model_database "go-simple-api/src/model/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// @Summary      Get Logged In User Profile
// @Description  Get profile details for the currently authenticated user by reading the JWT.
// @Tags         User
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} model_api.ApiResponseDoc "Successfully retrieved user profile"
// @Failure      401 {object} model_api.ApiResponseDoc "Unauthorized"
// @Failure      404 {object} model_api.ApiResponseDoc "User not found"
// @Router       /data/user [get]
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and a JWT token.
func GetUser(
	c *gin.Context,
) {
	idUserInterface, exists := c.Get("idUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, model_api.ApiResponseDoc{
			Status: &model_api.StatusModel{
				Code:    http.StatusUnauthorized,
				Message: "User identifier not found in context",
			},
		})
		return
	}

	idUser := idUserInterface.(string)

	var user model_database.User
	queryCondition := model_database.User{IDUser: idUser}

	if err := config.DB.
		Preload("Wallets").
		Where(&queryCondition).
		First(&user).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, model_api.ApiResponseDoc{
				Status: &model_api.StatusModel{
					Code:    http.StatusNotFound,
					Message: "User with the given ID not found",
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, model_api.ApiResponseDoc{
			Status: &model_api.StatusModel{
				Code:    http.StatusInternalServerError,
				Message: "Database error",
			},
		})
		return
	}

	response := model_api.ApiResponse[model_database.User]{
		Status: &model_api.StatusModel{
			Code:    http.StatusOK,
			Message: "User profile retrieved successfully",
		},
		Data: &user,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary      Edit User Profile
// @Description  Edit profile details for the currently authenticated user by reading the JWT dan Update Request.
// @Tags         User
// @Produce      json
// @Security     BearerAuth
// @Param        payload   body      model_api.UpdateUserRequest  true  "User Edit Payload"
// @Success      200 {object} model_api.ApiResponseDoc "Successfully retrieved user profile"
// @Failure      401 {object} model_api.ApiResponseDoc "Unauthorized"
// @Failure      404 {object} model_api.ApiResponseDoc "User not found"
// @Router       /data/user [put]
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and a JWT token.
func UpdateUser(
	c *gin.Context,
) {
	idUserInterface, exists := c.Get("idUser")

	if !exists {
		c.JSON(http.StatusUnauthorized, model_api.ApiResponseDoc{
			Status: &model_api.StatusModel{Code: http.StatusUnauthorized, Message: "User identifier not found in context"},
		})
		return
	}

	var input model_api.UpdateUserRequest
	var user model_database.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, model_api.ApiResponseDoc{
			Status: &model_api.StatusModel{Code: http.StatusBadRequest, Message: err.Error()},
		})
		return
	}

	idUser := idUserInterface.(string)
	queryCondition := model_database.User{
		IDUser: idUser,
	}

	if err := config.DB.Where(&queryCondition).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, model_api.ApiResponseDoc{
				Status: &model_api.StatusModel{Code: http.StatusNotFound, Message: "User with the given ID not found"},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, model_api.ApiResponseDoc{
			Status: &model_api.StatusModel{Code: http.StatusInternalServerError, Message: "Database error"},
		})
		return
	}

	if input.NameUser != "" {
		user.NameUser = input.NameUser
	}

	if input.PasswordUser != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.PasswordUser), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user.PasswordUser = string(hashedPassword)
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	response := model_api.ApiResponse[model_database.User]{
		Status: &model_api.StatusModel{Code: http.StatusOK, Message: "User profile updated successfully"},
		Data:   &user,
	}

	c.JSON(http.StatusOK, response)
}
