package main

import (
	"go-simple-api/src/config"
	"go-simple-api/src/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "go-simple-api/docs"
)

// @title           Simple API with Gin
// @version         1.0
// @description     This is a sample server for a simple API.
// @host            localhost:3000
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and a JWT token.
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := config.GetPort()
	if port == "" {
		port = "8080"
	}

	config.ConnectDB()
	config.MigrateDB()

	r := gin.Default()
	routes.SetupRoutes(r)

	r.Run(":" + port)
}
