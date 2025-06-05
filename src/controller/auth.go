package api_controller

import (
	"go-simple-api/src/config"
	model_api "go-simple-api/src/model/api"
	model_database "go-simple-api/src/model/database"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// @Summary      User Sign-In
// @Description  Authenticates a user and returns a JWT token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        credentials body      model_api.SignInRequest   true  "User Credentials"
// @Success      200         {object}  model_api.ApiResponseDoc "Sign in successful"
// @Failure      400         {object}  model_api.ApiResponseDoc  "Bad Request"
// @Failure      401         {object}  model_api.ApiResponseDoc  "Unauthorized"
// @Failure      500         {object}  model_api.ApiResponseDoc  "Internal Server Error"
// @Router       /auth/sign-in [post]
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

// @Summary      User Registration
// @Description  Registers a new user with name, email, and password.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload   body      model_api.RegisterRequest  true  "User Registration Payload"
// @Success      201       {object}  model_api.ApiResponseDoc	"Successfully registered"
// @Failure      400       {object}  model_api.ApiResponseDoc	"Bad Request"
// @Failure      500       {object}  model_api.ApiResponseDoc	"Internal Server Error"
// @Router       /auth/register [post]
func Register(c *gin.Context) {
	var input model_api.RegisterRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.PasswordUser), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	newUser := model_database.User{
		IDUser:       uuid.New().String(),
		NameUser:     input.NameUser,
		EmailUser:    input.EmailUser,
		PasswordUser: string(hashedPassword),
	}

	if err := config.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	response := model_api.ApiResponse[any]{
		Status: &model_api.StatusModel{
			Code:    http.StatusCreated,
			Message: "register is success!",
		},
	}

	c.JSON(http.StatusCreated, response)
}
