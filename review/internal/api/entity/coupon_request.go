package entity

type CouponRequest struct {
	Code string `json:"code" binding:"required"`
}
