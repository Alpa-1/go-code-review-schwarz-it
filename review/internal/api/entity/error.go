package entity

type APIError struct {
	Code    int    `json:"code" binding:"required"`
	Message string `json:"message" binding:"required"`
}
