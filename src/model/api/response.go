package model_api

type StatusModel struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ApiResponse[T any] struct {
	Status *StatusModel `json:"status,omitempty"`
	Data   *T           `json:"data,omitempty"`
}

type ApiResponseDoc struct {
	Status *StatusModel `json:"status,omitempty"`
	Data   interface{}  `json:"data,omitempty"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
