package bootstrap

import (
	"context"
	"log/slog"

	ordersgcl "github.com/m11ano/mipt-webdev-course/backend/temporal-app/internal/clients/grpc/orders"
	productsgcl "github.com/m11ano/mipt-webdev-course/backend/temporal-app/internal/clients/grpc/products"
	"github.com/m11ano/mipt-webdev-course/backend/temporal-app/internal/infra/config"
	"github.com/m11ano/mipt-webdev-course/backend/temporal-app/internal/infra/temporal"
	"github.com/m11ano/mipt-webdev-course/backend/temporal-app/pkg/workers/products/activities"
	"go.uber.org/fx"
)

var App = fx.Options(
	// Инфраструктура
	fx.Provide(NewLogger),
	fx.WithLogger(ProvideFXLogger),
	fx.Provide(ProvideGRPCClientsConns),
	fx.Provide(ProdiveTemporalAndConnect),
	fx.Provide(ProdiveTemporalActivities),
	// Start && Stop invoke
	fx.Invoke(func(lc fx.Lifecycle, shutdowner fx.Shutdowner, logger *slog.Logger, config config.Config, productsGCl *productsgcl.ClientConn, ordersGCl *ordersgcl.ClientConn, tClient temporal.TemporalClient, productsActivities *activities.Controller) {

		productsWorker := RegisterWorkers(tClient, productsActivities)

		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				err := ConnectToGRPCServer(ctx, productsGCl.Conn, ordersGCl.Conn)
				if err != nil {
					return err
				}
				logger.Info("Connection to gRPC servers established")

				err = TemporalCheckHealth(ctx, tClient)
				if err != nil {
					logger.ErrorContext(ctx, "Cant connect to temporal")
					return err
				}
				logger.Info("Connection to temporal established")

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
