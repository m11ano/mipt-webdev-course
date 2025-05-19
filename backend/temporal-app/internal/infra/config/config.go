package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App struct {
		StartTimeout int  `yaml:"start_timeout" env:"APP_START_TIMEOUT" env-default:"60"`
		StopTimeout  int  `yaml:"stop_timeout" env:"APP_STOP_TIMEOUT" env-default:"10"`
		IsProd       bool `yaml:"is_prod" env:"APP_IS_PROD" env-default:"false"`
		UseLogger    bool `yaml:"use_logger" env:"APP_USE_LOGGER" env-default:"true"`
		UseFxLogger  bool `yaml:"use_fx_logger" env:"APP_USE_FX_LOGGER" env-default:"false"`
	} `yaml:"app"`
	Temporal struct {
		Endpoint string `yaml:"endpoint" env:"TEMPORAL_ENDPOINT" env-default:"127.0.0.1:7233"`
	} `yaml:"temporal"`
	GRPC struct {
		Port    int `yaml:"port" env:"GRPC_PORT" env-default:"50051"`
		Clients struct {
			Products struct {
				Endpoint  string `yaml:"endpoint" env:"GRPC_CLIENTS_PRODUCTS_ENDPOINT" env-default:"127.0.0.1:8090"`
				Retries   int    `yaml:"retries" env:"GRPC_CLIENTS_PRODUCTS_RETRIES" env-default:"3"`
				TimeoutMS int    `yaml:"timeout_ms" env:"GRPC_CLIENTS_PRODUCTS_TIMEOUT_MS" env-default:"100"`
			} `yaml:"products"`
		} `yaml:"clients"`
	} `yaml:"grpc"`
}

func LoadConfig(file string) Config {
	var Config Config

	err := cleanenv.ReadConfig(file, &Config)
	if err != nil {
		log.Fatal("config error", err)
	}

	return Config
}
