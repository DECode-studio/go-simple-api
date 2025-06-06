package routes

import (
	"go-simple-api/src/config"
	api_controller "go-simple-api/src/controller"
	"go-simple-api/src/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine) {

	router.GET(config.ROOT_API,
		func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Simple GO Api",
			})
		},
	)

	router.GET(
		config.SWAGGER_ROUTE,
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	apiV1 := router.Group(config.BASE_API_V1)
	{
		// Auth Route
		apiV1.POST(
			config.AUTH_REGISTER_ROUTE,
			api_controller.Register,
		)
		apiV1.POST(
			config.AUTH_SIGNIN_ROUTE,
			api_controller.SignIn,
		)

		// User Route
		apiV1.GET(
			config.USER_ROUTE,
			middleware.AuthMiddleware,
			api_controller.GetUser,
		)
		apiV1.PUT(
			config.USER_ROUTE,
			middleware.AuthMiddleware,
			api_controller.UpdateUser,
		)

		// Wallet Route
		apiV1.GET(
			config.WALLET_ROUTE,
			middleware.AuthMiddleware,
			api_controller.GetWallet,
		)
		apiV1.POST(
			config.WALLET_ROUTE,
			middleware.AuthMiddleware,
			api_controller.AddWallet,
		)
		apiV1.DELETE(
			config.WALLET_ROUTE,
			middleware.AuthMiddleware,
			api_controller.DeleteWallet,
		)
	}
}
