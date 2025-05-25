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
	Cron struct {
		Tasks struct {
			RemoveFilesInStorage struct {
				Schedule           string `yaml:"schedule" env:"CRON_TASK_REMOVE_FILES_IN_STORAGE_SCHEDULE" env-default:""`
				AfterRemoveSecords int    `yaml:"after_remove_seconds" env:"CRON_TASK_REMOVE_FILES_IN_STORAGE_AFTER_REMOVE_SECONDS" env-default:"60"`
			} `yaml:"remove_files_in_storage"`

			DeleteNotAssignedFiles struct {
				Schedule         string `yaml:"schedule" env:"CRON_TASK_DELETE_NOT_ASSIGNED_FILES_SCHEDULE" env-default:""`
				AfterCreateHours int    `yaml:"after_create_hours" env:"CRON_TASK_DELETE_NOT_ASSIGNED_FILES_AFTER_CREATE_HOURS" env-default:"24"`
			} `yaml:"delete_not_assigned_files"`
		} `yaml:"tasks"`
	} `yaml:"cron"`
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
	GRPC struct {
		Port    int `yaml:"port" env:"GRPC_PORT" env-default:"50051"`
		Clients struct {
			Orders struct {
				Endpoint  string `yaml:"endpoint" env:"GRPC_CLIENTS_ORDERS_ENDPOINT" env-default:"127.0.0.1:8091"`
				Retries   int    `yaml:"retries" env:"GRPC_CLIENTS_ORDERS_RETRIES" env-default:"3"`
				TimeoutMS int    `yaml:"timeout_ms" env:"GRPC_CLIENTS_ORDERS_TIMEOUT_MS" env-default:"100"`
			} `yaml:"orders"`
		} `yaml:"clients"`
	} `yaml:"grpc"`
	Storage struct {
		S3Endpoint  string `yaml:"s3_endpoint" env:"STORAGE_S3_ENDPOINT" env-default:""`
		S3AccessKey string `yaml:"s3_access_key" env:"STORAGE_S3_ACCESS_KEY" env-default:""`
		S3SecretKey string `yaml:"s3_secret_key" env:"STORAGE_S3_SECRET_KEY" env-default:""`
		S3Region    string `yaml:"s3_region" env:"STORAGE_S3_REGION" env-default:""`
		S3URL       string `yaml:"s3_url" env:"STORAGE_S3_URL" env-default:""`
	} `yaml:"storage"`
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
