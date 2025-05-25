package bootstrap

import (
	"log/slog"
	"os"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/internal/infra/config"
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/internal/infra/db/txmngr"
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
