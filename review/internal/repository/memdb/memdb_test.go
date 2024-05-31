package memdb

import (
	"coupon_service/internal/service/entity"
	"reflect"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name string
		args args
		want *Repository
	}{
		{"initialize repository", args{}, &Repository{entries: make(map[string]entity.Coupon, 0)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
				return
			}
		})
	}

}

func TestRepository_Save(t *testing.T) {
	mockRepo := New()
	type args struct {
		coupon entity.Coupon
	}
	tests := []struct {
		name string
		args args
		want entity.Coupon
	}{
		{"save coupon", args{coupon: entity.Coupon{Code: "SAVE10", Discount: 10}}, entity.Coupon{Code: "SAVE10", Discount: 10}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mockRepo.Save(tt.want); !reflect.DeepEqual(got, nil) {
				t.Errorf("Save() = %v, want %v", got, nil)
				return
			}
		})
	}
}

func TestRepository_FindByCode(t *testing.T) {
	mockRepo := New()
	type args struct {
		entryCoupon entity.Coupon
		code        string
	}
	tests := []struct {
		name string
		args args
		want entity.Coupon
	}{
		{"find coupon by code", args{entryCoupon: entity.Coupon{Code: "SAVE10", Discount: 10}, code: "SAVE10"}, entity.Coupon{Code: "SAVE10", Discount: 10}},
		{"Coupon does not exist", args{entryCoupon: entity.Coupon{Code: "SAVE10", Discount: 10}, code: "SAVE20"}, entity.Coupon{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.Save(tt.args.entryCoupon)
			if got, _ := mockRepo.FindByCode(tt.args.code); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindByCode() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func TestRepository_ConcurrentAccess(t *testing.T) {
	repo := New()
	coupon := entity.Coupon{Code: "SAVE10", Discount: 10}
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			repo.Save(coupon)
		}()
	}

	wg.Wait()

	savedCoupon, err := repo.FindByCode(coupon.Code)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if savedCoupon != coupon {
		t.Errorf("expected %v, got %v", coupon, savedCoupon)
	}
}
