package main

import (
	"time"

	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/bootstrap"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/infra/config"
	"go.uber.org/fx"
)

// @title Orders API
// @version 1.0
// @description API документация для orders
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.LoadConfig("config.yml")

	app := fx.New(
		fx.Options(
			fx.StartTimeout(time.Second*time.Duration(cfg.App.StartTimeout)),
			fx.StopTimeout(time.Second*time.Duration(cfg.App.StopTimeout)),
		),
		fx.Provide(func() config.Config {
			return cfg
		}),
		bootstrap.App,
	)

	app.Run()
}
