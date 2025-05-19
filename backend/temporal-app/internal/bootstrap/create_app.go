package bootstrap

import (
	"context"
	"log/slog"

	"github.com/m11ano/mipt-webdev-course/backend/temporal-app/internal/infra/config"
	"go.uber.org/fx"
)

var App = fx.Options(
	// Инфраструктура
	fx.Provide(NewLogger),
	fx.WithLogger(ProvideFXLogger),
	// Start && Stop invoke
	fx.Invoke(func(lc fx.Lifecycle, shutdowner fx.Shutdowner, logger *slog.Logger, config config.Config) {
		tClient, err := NewTemporalClient(config, logger)
		if err != nil {
			logger.Error("cant connect to temporal", slog.Any("error", err))
			shutdowner.Shutdown()
			return
		}
		productsWorker := RegisterWorkers(tClient)

		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				RunWorkers(logger, shutdowner, tClient, productsWorker)

				return nil
			},
			OnStop: func(_ context.Context) error {

				tClient.Close()

				return nil
			},
		})
	}),
)
