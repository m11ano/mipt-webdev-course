package bootstrap

import (
	"log/slog"
	"os"
	"time"

	"github.com/m11ano/mipt-webdev-course/backend/temporal-app/internal/infra/config"
	"go.uber.org/fx/fxevent"

	productsgcl "github.com/m11ano/mipt-webdev-course/backend/temporal-app/internal/clients/grpc/products"
)

func ProvideFXLogger(config config.Config) fxevent.Logger {
	if !config.App.UseFxLogger {
		return fxevent.NopLogger
	}
	return &fxevent.ConsoleLogger{
		W: os.Stdout,
	}
}

func ProvideGRPCClientsConns(cfg config.Config, logger *slog.Logger) *productsgcl.ClientConn {

	products, err := productsgcl.NewClientConn(
		cfg.GRPC.Clients.Products.Endpoint,
		cfg.GRPC.Clients.Products.Retries,
		time.Duration(cfg.GRPC.Clients.Products.TimeoutMS)*time.Millisecond,
		logger,
	)
	if err != nil {
		panic(err)
	}

	return products

}
