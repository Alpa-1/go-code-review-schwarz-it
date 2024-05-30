package entity

type ApplicationRequest struct {
	Code        string `json:"code" binding:"required"`
	BasketValue int    `json:"basket_value" binding:"required"`
}
