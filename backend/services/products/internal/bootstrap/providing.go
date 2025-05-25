package bootstrap

import (
	"log/slog"
	"os"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	ordersgcl "github.com/m11ano/mipt-webdev-course/backend/services/products/internal/clients/grpc/orders"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/infra/config"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/infra/db/txmngr"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/infra/storage"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/infra/storage/s3d"
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

func ProvideStorageClient(s3Client *s3.Client) storage.Client {
	return s3d.New(s3Client)
}

func ProvideS3Client(cfg config.Config) *s3.Client {
	return s3d.NewS3Client(cfg.Storage.S3Endpoint, cfg.Storage.S3Region, cfg.Storage.S3AccessKey, cfg.Storage.S3SecretKey)
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

func ProvideGRPCClientsConns(cfg config.Config, logger *slog.Logger) *ordersgcl.ClientConn {

	products, err := ordersgcl.NewClientConn(
		cfg.GRPC.Clients.Orders.Endpoint,
		cfg.GRPC.Clients.Orders.Retries,
		time.Duration(cfg.GRPC.Clients.Orders.TimeoutMS)*time.Millisecond,
		logger,
	)
	if err != nil {
		panic(err)
	}

	return products

}
