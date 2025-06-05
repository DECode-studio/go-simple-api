package api_controller

import (
	"go-simple-api/src/config"
	model_api "go-simple-api/src/model/api"
	model_database "go-simple-api/src/model/database"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// @Summary      Get List Wallet Address Of User
// @Description  Get list wallet address of currently authenticated user by reading the JWT.
// @Tags         Wallet
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} model_api.ApiResponseDoc "Successfully retrieved user profile"
// @Failure      401 {object} model_api.ApiResponseDoc "Unauthorized"
// @Failure      404 {object} model_api.ApiResponseDoc "User not found"
// @Router       /data/wallet [get]
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and a JWT token.
func GetWallet(
	c *gin.Context,
) {
	idUserInterface, exists := c.Get("idUser")
	if !exists {

		c.JSON(http.StatusUnauthorized, model_api.ApiResponseDoc{
			Status: &model_api.StatusModel{Code: http.StatusUnauthorized, Message: "User identifier not found in context"},
		})
		return
	}

	var wallets []model_database.Wallet
	idUser := idUserInterface.(string)
	queryCondition := model_database.Wallet{
		IDUser: idUser,
	}

	if err := config.DB.Where(&queryCondition).Find(&wallets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model_api.ApiResponseDoc{
			Status: &model_api.StatusModel{
				Code:    http.StatusInternalServerError,
				Message: "Failed to fetch wallets",
			},
		})
		return
	}

	response := model_api.ApiResponse[[]model_database.Wallet]{
		Status: &model_api.StatusModel{
			Code:    http.StatusOK,
			Message: "Wallets retrieved successfully",
		},
		Data: &wallets,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary      Add Wallet Address Of User
// @Description  Add wallet address of currently authenticated user by reading the JWT.
// @Tags         Wallet
// @Produce      json
// @Security     BearerAuth
// @Param        payload   body      model_api.AddWalletRequest  true  "Add Wallet Payload"
// @Success      200 {object} model_api.ApiResponseDoc "Successfully retrieved user profile"
// @Failure      401 {object} model_api.ApiResponseDoc "Unauthorized"
// @Failure      404 {object} model_api.ApiResponseDoc "User not found"
// @Router       /data/wallet [post]
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and a JWT token.
func AddWallet(
	c *gin.Context,
) {
	idUserInterface, exists := c.Get("idUser")
	if !exists {

		c.JSON(http.StatusUnauthorized, model_api.ApiResponseDoc{
			Status: &model_api.StatusModel{Code: http.StatusUnauthorized, Message: "User identifier not found in context"},
		})
		return
	}

	var input model_api.AddWalletRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !common.IsHexAddress(input.WalletAddress) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid wallet address",
		})
		return
	}

	idUser := idUserInterface.(string)
	queryCondition := model_database.Wallet{
		IDUser:        idUser,
		AddressWallet: input.WalletAddress,
	}

	var existingWallet model_database.Wallet
	if err := config.DB.Where(&queryCondition).First(&existingWallet).Error; err == nil {
		c.JSON(http.StatusBadRequest, model_api.ApiResponseDoc{
			Status: &model_api.StatusModel{
				Code:    http.StatusBadRequest,
				Message: "Wallet address already exists for this user",
			},
		})
		return
	}

	newWallet := model_database.Wallet{
		IDWallet:      uuid.New().String(),
		IDUser:        idUser,
		AddressWallet: input.WalletAddress,
	}

	if err := config.DB.Create(&newWallet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add wallet"})
		return
	}

	response := model_api.ApiResponse[model_database.Wallet]{
		Status: &model_api.StatusModel{
			Code:    http.StatusOK,
			Message: "Wallet add created successfully",
		},
		Data: &newWallet,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary      Delete Wallet Address By ID Wallet
// @Description  Delete wallet address of ID WAllet and currently authenticated user by reading the JWT.
// @Tags         Wallet
// @Produce      json
// @Security     BearerAuth
// @Param        payload   body      model_api.DeleteWalletRequest  true  "Add Wallet Payload"
// @Success      200 {object} model_api.ApiResponseDoc "Successfully delete wallet address"
// @Failure      401 {object} model_api.ApiResponseDoc "Unauthorized"
// @Failure      404 {object} model_api.ApiResponseDoc "User not found"
// @Router       /data/wallet [delete]
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and a JWT token.
func DeleteWallet(
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

	var input model_api.DeleteWalletRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idUser := idUserInterface.(string)

	queryCondition := model_database.Wallet{
		IDWallet: input.IDWallet,
		IDUser:   idUser,
	}

	var wallet model_database.Wallet
	if err := config.DB.Where(&queryCondition).First(&wallet).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, model_api.ApiResponseDoc{
				Status: &model_api.StatusModel{
					Code:    http.StatusNotFound,
					Message: "Wallet not found for the given user",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, model_api.ApiResponseDoc{
			Status: &model_api.StatusModel{
				Code:    http.StatusInternalServerError,
				Message: "Database error while finding wallet",
			},
		})
		return
	}

	if err := config.DB.Delete(&wallet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model_api.ApiResponseDoc{
			Status: &model_api.StatusModel{
				Code:    http.StatusInternalServerError,
				Message: "Failed to delete wallet",
			},
		})
		return
	}

	c.JSON(http.StatusOK, model_api.ApiResponseDoc{
		Status: &model_api.StatusModel{
			Code:    http.StatusOK,
			Message: "Wallet deleted successfully",
		},
	})
}
