package bootstrap

import (
	"log/slog"

	ordersgcl "github.com/m11ano/mipt-webdev-course/backend/temporal-app/internal/clients/grpc/orders"
	productsgcl "github.com/m11ano/mipt-webdev-course/backend/temporal-app/internal/clients/grpc/products"
	"github.com/m11ano/mipt-webdev-course/backend/temporal-app/internal/infra/config"
	"github.com/m11ano/mipt-webdev-course/backend/temporal-app/internal/infra/temporal"
	productsw "github.com/m11ano/mipt-webdev-course/backend/temporal-app/pkg/workers/products"
	"github.com/m11ano/mipt-webdev-course/backend/temporal-app/pkg/workers/products/activities"
	"github.com/m11ano/mipt-webdev-course/backend/temporal-app/pkg/workers/products/workflows"
	"go.uber.org/fx"

	tclient "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
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

func ProdiveTemporalActivities(productsGRPC *productsgcl.ClientConn, ordersGRPC *ordersgcl.ClientConn) *activities.Controller {
	activitiesController := activities.NewController(productsGRPC, ordersGRPC)
	return activitiesController
}

func RegisterWorkers(tClient temporal.TemporalClient, productsActivities *activities.Controller) productsw.ProductsWorker {
	productsWorker := productsw.NewWorker(tClient)

	productsWorker.RegisterWorkflow(workflows.SetOrderProductsAndBlock)

	productsWorker.RegisterActivity(productsActivities)

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
