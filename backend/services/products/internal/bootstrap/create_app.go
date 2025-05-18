package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/infra/config"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/infra/db/migrations"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/infra/storage"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/infra/storage/s3d"
	"go.uber.org/fx"
)

var App = fx.Options(
	// Инфраструктура
	fx.Provide(NewLogger),
	fx.WithLogger(ProvideFXLogger),
	fx.Provide(NewPgxv5),
	fx.Provide(ProvideFiberApp),
	fx.Provide(ProvidePGXPoolWithTxMgr),
	fx.Provide(ProvideS3Client),
	fx.Provide(ProvideStorageClient),
	// Бизнес логика
	FileModule,
	ProductModule,
	ProductSliderImageModule,
	ProductOrderBlockModule,
	// Delivery
	DeliveryHTTP,
	// Start && Stop invoke
	fx.Invoke(func(lc fx.Lifecycle, shutdowner fx.Shutdowner, logger *slog.Logger, config config.Config, dbpool *pgxpool.Pool, fiberApp *fiber.App, s3Client *s3.Client, storageClient storage.Client) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				err := Pgxv5TestConnection(ctx, dbpool, logger, config.DB.MaxAttempt, config.DB.AttemptSleepSeconds)
				if err != nil {
					return err
				}
				logger.Info("Postgress connected")

				err = migrations.RunMigrations(ctx, dbpool, config, logger)
				if err != nil {
					return err
				}

				err = s3d.PingS3Client(ctx, s3Client)
				if err != nil {
					return err
				}
				logger.Info("S3 connected")

				err = storage.CreateBuckets(ctx, storageClient)
				if err != nil {
					return err
				}
				logger.Info("Storage buckets created")

				if config.HTTP.Port > 0 {
					go func() {
						if err := fiberApp.Listen(fmt.Sprintf(":%d", config.HTTP.Port)); err != nil {
							logger.Error("failed to start fiber", slog.Any("error", err), slog.Any("trackeback", string(debug.Stack())))
							err := shutdowner.Shutdown()
							if err != nil {
								logger.Error("failed to shutdown", slog.Any("error", err), slog.Any("trackeback", string(debug.Stack())))
							}
						}
					}()
				}

				return nil
			},
			OnStop: func(_ context.Context) error {
				if config.HTTP.Port > 0 {
					logger.Info("stopping HTTP Fiber")
					err := fiberApp.ShutdownWithTimeout(time.Duration(config.HTTP.StopTimeout) * time.Second)
					if err != nil {
						logger.Error("failed to stop fiber", slog.Any("error", err), slog.Any("trackeback", string(debug.Stack())))
					}
				}

				logger.Info("stopping Postgress")
				dbpool.Close()

				return nil
			},
		})
	}),
)
