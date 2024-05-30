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
		return nil, err
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

// returns a list of valid coupons based on the codes provided
func (s Service) GetCoupons(codes []string) ([]Coupon, error) {
	coupons := make([]Coupon, 0, len(codes))
	var e error = nil

	for idx, code := range codes {
		coupon, err := s.repo.FindByCode(code)
		if err != nil {
			if e == nil {
				e = fmt.Errorf("code: %s, index: %d", code, idx)
			} else {
				e = fmt.Errorf("%w; code: %s, index: %d", e, code, idx)
			}
		} else {
			coupons = append(coupons, *coupon)
		}
	}

	return coupons, e
}
