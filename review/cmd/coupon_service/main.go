package main

import (
	"coupon_service/internal/api"
	"coupon_service/internal/config"
	"coupon_service/internal/repository/memdb"
	"coupon_service/internal/service"
	"errors"
	"log"
	"runtime"
)

var (
	cfg  = config.New()
	repo = memdb.New()
)

func main() {
	if err := checkCores(); err != nil && cfg.Environment == "prod" {
		log.Fatal(err)
		return
	}

	svc := service.New(repo)
	server := api.New(cfg, svc)
	log.Printf("Starting coupon service server on %v", cfg.API.Port)
	server.Start()
}

func checkCores() error {
	if 32 != runtime.NumCPU() {
		return errors.New("this api is meant to be run on 32 core machines.")
	}
	return nil
}
