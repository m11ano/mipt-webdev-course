package bootstrap

import (
	"log/slog"

	"github.com/m11ano/mipt-webdev-course/backend/temporal-app/internal/infra/config"
	productsw "github.com/m11ano/mipt-webdev-course/backend/temporal-app/pkg/workers/products"
	"go.uber.org/fx"

	tclient "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func NewTemporalClient(config config.Config, logger *slog.Logger) (tclient.Client, error) {
	c, err := tclient.Dial(tclient.Options{
		HostPort: config.Temporal.Endpoint,
		Logger:   logger.WithGroup("temporal"),
	})
	if err != nil {
		return nil, err
	}

	return c, nil
}

func RegisterWorkers(tClient tclient.Client) productsw.ProductsWorker {
	productsWorker := productsw.NewWorker(tClient)

	return productsWorker
}

func RunWorkers(logger *slog.Logger, shutdowner fx.Shutdowner, tClient tclient.Client, productsWorker productsw.ProductsWorker) {
	go func() {
		err := productsWorker.Run(worker.InterruptCh())
		if err != nil {
			logger.Error("error in productsWorker", slog.Any("error", err))
			shutdowner.Shutdown()
		}
	}()
}
