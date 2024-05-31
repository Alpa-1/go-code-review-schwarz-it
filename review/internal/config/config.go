// Config package provides the centralized configuration for the coupon_service api.
package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	API struct {
		Host string `yaml:"host" env:"API_HOST" env-description:"API Service host" env-default:"localhost"`
		Port int    `yaml:"port" env:"API_PORT" env-description:"API Service port" env-default:"8080"`
	} `yaml:"api"`
	Environment string `yaml:"environment" env:"ENVIRONMENT" env-description:"Sets the running environment to either 'prod' or 'dev'" env-default:"dev"`
}

// New returns a new Config instance. It reads the configuration from the config.yml in the root directory file and then the environment variables.
func New() Config {
	var cfg Config
	err := cleanenv.ReadConfig("config.yml", &cfg)

	if err != nil {
		log.Fatal(err)
	}
	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}
