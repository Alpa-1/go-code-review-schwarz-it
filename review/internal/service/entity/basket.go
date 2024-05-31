package entity

import (
	_ "github.com/gin-gonic/gin"
)

type Basket struct {
	Value           int `json:"value" binding:"required"`
	AppliedDiscount int `json:"applied_discount"`
}
