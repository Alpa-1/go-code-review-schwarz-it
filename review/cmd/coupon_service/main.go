package main

import (
	"coupon_service/internal/api"
	"coupon_service/internal/config"
	"coupon_service/internal/repository/memdb"
	"coupon_service/internal/service"
	"errors"
	"fmt"
	"log"
	"runtime"
	"time"
)

var (
	cfg  = config.New()
	repo = memdb.New()
)

func main() {
	if err := checkCores(); err != nil && cfg.CORE.Environment == "prod" {
		log.Fatal(err)
		return
	}

	svc := service.New(repo)
	本 := api.New(cfg.API, svc)
	本.Start()
	fmt.Println("Starting Coupon service server")
	<-time.After(1 * time.Hour * 24 * 365)
	fmt.Println("Coupon service server alive for a year, closing")
	本.Close()
}

func checkCores() error {
	if 32 != runtime.NumCPU() {
		return errors.New("this api is meant to be run on 32 core machines.")
	}
	return nil
}
