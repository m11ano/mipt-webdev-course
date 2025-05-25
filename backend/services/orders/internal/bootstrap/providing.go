package bootstrap

import (
	"log/slog"
	"os"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	productsgcl "github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/clients/grpc/products"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/infra/config"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/infra/db/txmngr"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/infra/temporal"
	productstc "github.com/m11ano/mipt-webdev-course/backend/temporal-app/pkg/workers/products/client"
	"go.uber.org/fx/fxevent"
)

func ProvideFXLogger(config config.Config) fxevent.Logger {
	if !config.App.UseFxLogger {
		return fxevent.NopLogger
	}
	return &fxevent.ConsoleLogger{
		W: os.Stdout,
	}
}

func ProvidePGXPoolWithTxMgr(pgxpool *pgxpool.Pool) (*manager.Manager, *trmpgx.CtxGetter) {
	return txmngr.New(pgxpool)
}

func ProvideFiberApp(cfg config.Config, logger *slog.Logger) *fiber.App {
	fiberApp := NewHTTPFiber(HTTPConfig{
		UnderProxy:       cfg.HTTP.UnderProxy,
		UseTraceID:       true,
		UseLogger:        true,
		BodyLimit:        -1,
		CorsAllowOrigins: cfg.HTTP.Cors,
	}, logger)
	return fiberApp
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

func ProvideTemporalClients(tClient temporal.TemporalClient) productstc.Client {
	productsClient := productstc.NewClientImpl(tClient)

	return productsClient
}
