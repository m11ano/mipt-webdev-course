package orderscl

import (
	"log/slog"
	"time"

	"github.com/m11ano/mipt-webdev-course/backend/clients/clgrpc/dial"
	ordersv1 "github.com/m11ano/mipt-webdev-course/backend/protos/gen/go/orders"
	"google.golang.org/grpc"
)

type ClientImpl struct {
	api    ordersv1.OrdersClient
	logger *slog.Logger
}

func NewClientImpl(addr string, retriesCount int, timeout time.Duration, logger *slog.Logger) (*ClientImpl, *grpc.ClientConn, error) {
	cfg := dial.Config{
		Addr:         addr,
		RetriesCount: retriesCount,
		Timeout:      timeout,
	}

	cc, err := dial.NewClientConn(cfg, logger)
	if err != nil {
		return nil, nil, err
	}

	return &ClientImpl{
		api:    ordersv1.NewOrdersClient(cc),
		logger: logger,
	}, cc, nil
}
