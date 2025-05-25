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
	DB struct {
		URI                 string `yaml:"uri" env:"DB_URI"`
		MaxAttempt          int    `yaml:"max_attempts" env:"DB_MAX_ATTEMPTS" env-default:"10"`
		AttemptSleepSeconds int    `yaml:"attempt_sleep_seconds" env:"DB_ATTEMPT_SLEEP_SECONDS" env-default:"3"`
		MigrationsPath      string `yaml:"migrations_path" env:"DB_MIGRATIONS_PATH" env-default:"migrations"`
	} `yaml:"db"`
	HTTP struct {
		Prefix       string   `yaml:"prefix" env:"HTTP_PREFIX" env-default:""`
		ServiceName  string   `yaml:"service_name" env:"HTTP_SERVICE_NAME" env-default:""`
		Port         int      `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
		StopTimeout  int      `yaml:"stop_timeout" env:"HTTP_STOP_TIMEOUT" env-default:"5"`
		UnderProxy   bool     `yaml:"under_proxy" env:"HTTP_UNDER_PROXY" env-default:"false"`
		StartSwagger bool     `yaml:"start_swagger" env:"HTTP_START_SWAGGER" env-default:"false"`
		Cors         []string `yaml:"cors" env:"HTTP_CORS"`
	} `yaml:"http"`
	Auth struct {
		TokenTTLHours int `yaml:"token_ttl_hours" env:"AUTH_TOKEN_TTL_HOURS" env-default:"24"`
	} `yaml:"auth"`
	Secrets struct {
		JWT string `yaml:"jwt" env:"SECRETS_JWT" env-default:""`
	} `yaml:"secrets"`
}

func LoadConfig(file string) Config {
	var Config Config

	err := cleanenv.ReadConfig(file, &Config)
	if err != nil {
		log.Fatal("config error", err)
	}

	return Config
}
