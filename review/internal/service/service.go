package service

import (
	. "coupon_service/internal/service/entity"
	"fmt"

	"github.com/google/uuid"
)

type Repository interface {
	FindByCode(string) (*Coupon, error)
	Save(Coupon) error
}

type Service struct {
	repo Repository
}

func New(repo Repository) Service {
	return Service{
		repo: repo,
	}
}

func (s Service) ApplyCoupon(basket Basket, code string) (*Basket, error) {
	b := &basket
	coupon, err := s.repo.FindByCode(code)
	if err != nil {
		return nil, fmt.Errorf("coupon for code '%s' not found: %w", code, err)
	}

	if b.Value > 0 && b.Value >= coupon.MinBasketValue {
		b.AppliedDiscount = coupon.Discount
		b.ApplicationSuccessful = true
		return b, nil
	}
	if b.Value > 0 && b.Value < coupon.MinBasketValue {
		return nil, fmt.Errorf("basket value is less than minimum basket value")
	}

	if b.Value == 0 {
		return nil, fmt.Errorf("basket value is zero")
	}

	return nil, fmt.Errorf("tried to apply discount to negative value")
}

// creates a new coupon
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

// validate a coupon by a given code
func (s Service) ValidateCoupon(code string) (Coupon, error) {
	coupon, e := s.repo.FindByCode(code)

	if e != nil {
		return Coupon{}, e
	}

	return *coupon, nil
}
