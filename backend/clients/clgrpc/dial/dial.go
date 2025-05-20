package dial

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
)

type Config struct {
	Addr         string
	RetriesCount int
	Timeout      time.Duration
}

func NewClientConn(cfg Config, logger *slog.Logger) (*grpc.ClientConn, error) {
	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.Internal, codes.Unavailable, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(cfg.RetriesCount)),
		grpcretry.WithPerRetryTimeout(cfg.Timeout),
	}

	logOpts := []logging.Option{
		logging.WithLogOnEvents(logging.PayloadReceived, logging.PayloadSent),
	}

	intercepts := []grpc.UnaryClientInterceptor{
		grpcretry.UnaryClientInterceptor(retryOpts...),
	}

	if logger != nil {
		intercepts = append(intercepts, logging.UnaryClientInterceptor(InterceptorLogger(logger), logOpts...))
	}

	cc, err := grpc.NewClient(cfg.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(intercepts...),
	)
	if err != nil {
		return nil, err
	}

	return cc, nil
}

func InterceptorLogger(logger *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		logger.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
