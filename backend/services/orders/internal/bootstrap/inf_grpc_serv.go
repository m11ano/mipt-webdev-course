package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/infra/config"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewGRPCServer(cfg config.Config, logger *slog.Logger) *grpc.Server {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived,
			logging.PayloadSent,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			logger.Error("recovered from panic", slog.Any("panic", p))
			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	interceptors := []grpc.UnaryServerInterceptor{recovery.UnaryServerInterceptor(recoveryOpts...)}

	if !cfg.App.IsProd {
		interceptors = append(interceptors, logging.UnaryServerInterceptor(gRPCServerInterceptorLogger(logger), loggingOpts...))
	}

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors...))

	return gRPCServer
}

func gRPCServerInterceptorLogger(logger *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		logger.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func StartGRPCServer(grpcServer *grpc.Server, cfg config.Config, logger *slog.Logger, shutdowner fx.Shutdowner) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
	if err != nil {
		logger.Error("failed to listen gRPC", slog.Any("error", err))
		_ = shutdowner.Shutdown()
		return
	}

	logger.Info("gRPC started", slog.Int("port", cfg.GRPC.Port))

	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("failed to serve gRPC", slog.Any("error", err))
		_ = shutdowner.Shutdown()
	}
}
