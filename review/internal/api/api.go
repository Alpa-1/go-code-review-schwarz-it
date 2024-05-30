package api

import (
	"context"
	"coupon_service/internal/config"
	"coupon_service/internal/service/entity"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Service interface {
	ApplyCoupon(entity.Basket, string) (*entity.Basket, error)
	CreateCoupon(int, string, int) error
	GetCoupons([]string) ([]entity.Coupon, error)
}

type API struct {
	srv *http.Server
	MUX *gin.Engine
	svc Service
	CFG config.Config
}

func New[T Service](cfg config.Config, svc T) API {
	gin.SetMode(gin.ReleaseMode)
	r := new(gin.Engine)
	r = gin.New()
	r.Use(gin.Recovery())

	return API{
		MUX: r,
		CFG: cfg,
		svc: svc,
	}.withServer().withRoutes()
}

func (a API) withServer() API {
	a.srv = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", a.CFG.API.Host, a.CFG.API.Port),
		Handler: a.MUX,
	}
	return a
}

func (a API) withRoutes() API {
	apiGroup := a.MUX.Group("/api")
	apiGroup.POST("/apply", a.Apply)
	apiGroup.POST("/create", a.Create)
	apiGroup.GET("/coupons", a.Get)
	return a
}

func (a API) Start() {
	if err := a.srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func (a API) Close() {
	<-time.After(5 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.srv.Shutdown(ctx); err != nil {
		log.Println(err)
	}
}
