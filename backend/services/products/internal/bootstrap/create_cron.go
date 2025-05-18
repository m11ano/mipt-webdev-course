package bootstrap

import (
	"context"
	"log/slog"
	"runtime/debug"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/infra/config"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/infra/db/migrations"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/infra/storage"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/infra/storage/s3d"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/usecase"
	"github.com/robfig/cron/v3"
	"go.uber.org/fx"
)

var Cron = fx.Options(
	// Инфраструктура
	fx.Provide(NewLogger),
	fx.WithLogger(ProvideFXLogger),
	fx.Provide(NewPgxv5),
	fx.Provide(ProvidePGXPoolWithTxMgr),
	fx.Provide(ProvideS3Client),
	fx.Provide(ProvideStorageClient),
	// Бизнес логика
	FileModule,
	// Start && Stop invoke
	fx.Invoke(func(lc fx.Lifecycle, shutdowner fx.Shutdowner, logger *slog.Logger, config config.Config, dbpool *pgxpool.Pool, s3Client *s3.Client, storageClient storage.Client, fileUsecase usecase.File) {
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

				c := cron.New(cron.WithSeconds())

				//Удаление помеченных к удалению файлов из хранилища
				_, err = c.AddFunc(config.Cron.Tasks.RemoveFilesInStorage.Schedule, func() {
					defer func() {
						if r := recover(); r != nil {
							logger.Error("panic", slog.Any("error", r), slog.Any("trackeback", string(debug.Stack())))
						}
					}()

					logger.Info("start task remove files from storage")

					now := time.Now()
					maxDeletedAt := now.Add(-1 * time.Duration(config.Cron.Tasks.RemoveFilesInStorage.AfterRemoveSecords) * time.Second)
					fileUsecase.TaskRemoveFilesFromStorage(context.Background(), maxDeletedAt)
				})
				if err != nil {
					return err
				}

				//Удаление неиспользованных файлов
				_, err = c.AddFunc(config.Cron.Tasks.DeleteNotAssignedFiles.Schedule, func() {
					defer func() {
						if r := recover(); r != nil {
							logger.Error("panic", slog.Any("error", r), slog.Any("trackeback", string(debug.Stack())))
						}
					}()

					logger.Info("start task delete not assigned files")

					now := time.Now()
					maxCreatedAt := now.Add(-1 * time.Duration(config.Cron.Tasks.DeleteNotAssignedFiles.AfterCreateHours) * time.Hour)
					fileUsecase.TaskDeleteNotAssignedFiles(context.Background(), maxCreatedAt)
				})
				if err != nil {
					return err
				}

				c.Start()

				return nil
			},
			OnStop: func(_ context.Context) error {
				logger.Info("stopping Postgress")
				dbpool.Close()

				return nil
			},
		})
	}),
)
