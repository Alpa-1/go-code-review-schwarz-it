package config

import (
	"coupon_service/internal/api"
	"log"

	"github.com/brumhard/alligotor"
)

type coreConfig struct {
	Environment string
}
type Config struct {
	API  api.Config
	CORE coreConfig
}

func New() Config {
	cfg := Config{}
	if err := alligotor.Get(&cfg); err != nil {
		log.Fatal(err)
	}
	return cfg
}
