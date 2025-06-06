package model_api

type SignInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	NameUser     string `json:"nameUser" binding:"required"`
	EmailUser    string `json:"emailUser" binding:"required,email"`
	PasswordUser string `json:"passwordUser" binding:"required,min=8"`
}

type UpdateUserRequest struct {
	NameUser     string `json:"nameUser"`
	PasswordUser string `json:"passwordUser"`
}

type AddWalletRequest struct {
	WalletAddress string `json:"walletAddress" binding:"required"`
}

type DeleteWalletRequest struct {
	IDWallet string `json:"idWallet" binding:"required"`
}
