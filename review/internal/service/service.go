// Package service provides the business logic for the coupon service. It relies on the repository to interact with the database.
package service

import (
	. "coupon_service/internal/service/entity"
	"fmt"

	"github.com/google/uuid"
)

type Repository interface {
	FindByCode(string) (Coupon, error)
	Save(Coupon) error
}

type Service struct {
	repo Repository
}

// New creates a new service with the provided repository.
func New(repo Repository) Service {
	return Service{
		repo: repo,
	}
}

// ApplyCoupon applies a coupon to a basket and returns the updated basket if the code was applicable.
//
//	basketValue: the value of the basket in cents
//	code: the coupon code
func (s Service) ApplyCoupon(basketValue int, code string) (*Basket, error) {
	basket := Basket{Value: basketValue, AppliedDiscount: 0}
	b := &basket

	coupon, err := s.repo.FindByCode(code)
	if err != nil {
		return nil, fmt.Errorf("coupon for code '%s' not found: %w", code, err)
	}

	if basketValue > 0 && basketValue >= coupon.MinBasketValue {
		b.Value = int(basketValue - (basketValue * coupon.Discount / 100))
		b.AppliedDiscount = coupon.Discount
		return b, nil
	}
	if basketValue > 0 && basketValue < coupon.MinBasketValue {
		return nil, fmt.Errorf("basket value is less than minimum basket value")
	}

	if basketValue == 0 {
		return nil, fmt.Errorf("basket value is zero")
	}

	return nil, fmt.Errorf("tried to apply discount to negative value")
}

// CreateCoupon creates a new coupon
//
//	discount: the discount percentage in whole percentages
//	code: the coupon code
//	minBasketValue: the minimum basket value required to apply the coupon in cents
func (s Service) CreateCoupon(discount int, code string, minBasketValue int) error {
	if discount < 0 {
		return fmt.Errorf("discount cannot be negative")
	}

	if discount > 100 {
		return fmt.Errorf("discount cannot be greater than 100")
	}

	if minBasketValue < 0 {
		return fmt.Errorf("minimum basket value cannot be negative")
	}

	if code == "" {
		return fmt.Errorf("code cannot be empty")
	}

	coupon := Coupon{
		Discount:       discount,
		Code:           code,
		MinBasketValue: minBasketValue,
		ID:             uuid.NewString(),
	}

	if err := s.repo.Save(coupon); err != nil {
		return err
	}
	return nil
}

// ValidateCoupon retrieves a coupon by its code from the repository.
func (s Service) ValidateCoupon(code string) (Coupon, error) {
	coupon, e := s.repo.FindByCode(code)

	if e != nil {
		return Coupon{}, e
	}

	return coupon, nil
}
