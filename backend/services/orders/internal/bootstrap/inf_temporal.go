package bootstrap

import (
	"context"
	"log/slog"

	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/infra/config"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/infra/temporal"
	tclient "go.temporal.io/sdk/client"
	"go.uber.org/fx"
)

func ProdiveTemporalAndConnect(config config.Config, logger *slog.Logger, shutdowner fx.Shutdowner) temporal.TemporalClient {
	c, err := tclient.Dial(tclient.Options{
		HostPort: config.Temporal.Endpoint,
		Logger:   logger.WithGroup("temporal"),
	})
	if err != nil {
		logger.Error("cant connect to temporal", slog.Any("error", err))
		shutdowner.Shutdown()
	}

	return c
}

func TemporalCheckHealth(ctx context.Context, tClient tclient.Client) error {
	_, err := tClient.CheckHealth(ctx, &tclient.CheckHealthRequest{})
	if err != nil {
		return err
	}

	return nil
}
