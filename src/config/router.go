package config

import "os"

func GetPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return "3000"
}

const (
	// Base API
	BASE_API_V1 = "/api/v1"

	// Swagger ROUTE
	SWAGGER_ROUTE = "/swagger/*any"

	// Auth ROUTE
	AUTH_SIGNIN_ROUTE   = "/auth/sign-in"
	AUTH_REGISTER_ROUTE = "/auth/register"

	// Data ROUTE
	USER_ROUTE   = "/data/user"
	WALLET_ROUTE = "/data/wallet"
)
