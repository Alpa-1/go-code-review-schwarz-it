package service

import (
	"coupon_service/internal/repository/memdb"
	"coupon_service/internal/service/entity"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	mockRepo := memdb.New()
	type args struct {
		repo Repository
	}
	tests := []struct {
		name string
		args args
		want Service
	}{
		{"initialize nil service", args{repo: nil}, Service{repo: nil}},
		{"initialize service", args{repo: mockRepo}, Service{repo: mockRepo}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_ApplyCoupon(t *testing.T) {
	type fields struct {
		repo Repository
	}
	type args struct {
		basketValue    int
		discount       int
		code           string
		minBasketValue int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantB   *entity.Basket
		wantErr bool
	}{
		{"Apply 10%", fields{memdb.New()}, args{100, 10, "Superdiscount", 55}, &entity.Basket{Value: 90, AppliedDiscount: 10}, false},
		{"Apply 33%", fields{memdb.New()}, args{100, 33, "Superdiscount", 55}, &entity.Basket{Value: 67, AppliedDiscount: 33}, false},
		{"Apply 100%", fields{memdb.New()}, args{100, 100, "Superdiscount", 55}, &entity.Basket{Value: 0, AppliedDiscount: 100}, false},
		{"Fail Apply on Negative Value", fields{memdb.New()}, args{-100, 10, "Superdiscount", 55}, nil, true},
		{"Fail Apply on MinBasketValue", fields{memdb.New()}, args{100, 10, "Superdiscount", 200}, nil, true},
		{"Fail Apply on Zero Value", fields{memdb.New()}, args{0, 10, "Superdiscount", 55}, nil, true},
		{"Fail Apply on Empty Code", fields{memdb.New()}, args{100, 10, "", 55}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				repo: tt.fields.repo,
			}
			s.CreateCoupon(tt.args.discount, tt.args.code, tt.args.minBasketValue)

			gotB, err := s.ApplyCoupon(tt.args.basketValue, tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApplyCoupon() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotB, tt.wantB) {
				t.Errorf("ApplyCoupon() gotB = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}

func TestService_CreateCoupon(t *testing.T) {
	type fields struct {
		repo Repository
	}
	type args struct {
		discount       int
		code           string
		minBasketValue int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Apply 10%", fields{memdb.New()}, args{10, "Superdiscount", 55}, false},
		{"Negative Discount", fields{memdb.New()}, args{-10, "Superdiscount", 55}, true},
		{"Discount > 100", fields{memdb.New()}, args{101, "Superdiscount", 55}, true},
		{"Negative MinBasketValue", fields{memdb.New()}, args{10, "Superdiscount", -55}, true},
		{"Empty Code", fields{memdb.New()}, args{10, "", 55}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				repo: tt.fields.repo,
			}

			err := s.CreateCoupon(tt.args.discount, tt.args.code, tt.args.minBasketValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCoupon() error = %v, wantError %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				coupon, err := s.repo.FindByCode(tt.args.code)
				if err != nil {
					t.Errorf("FindByCode() error = %v", err)
					return
				}

				if coupon.Discount != tt.args.discount {
					t.Errorf("CreateCoupon() discount = %v, want %v", coupon.Discount, tt.args.discount)
				}
				if coupon.Code != tt.args.code {
					t.Errorf("CreateCoupon() code = %v, want %v", coupon.Code, tt.args.code)
				}
				if coupon.MinBasketValue != tt.args.minBasketValue {
					t.Errorf("CreateCoupon() minBasketValue = %v, want %v", coupon.MinBasketValue, tt.args.minBasketValue)
				}
			}

		})
	}
}
