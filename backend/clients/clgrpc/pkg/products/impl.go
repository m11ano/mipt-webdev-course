package productscl

import (
	"log/slog"
	"time"

	"github.com/m11ano/mipt-webdev-course/backend/clients/clgrpc/dial"
	productsv1 "github.com/m11ano/mipt-webdev-course/backend/protos/gen/go/products"
	"google.golang.org/grpc"
)

type ClientImpl struct {
	api    productsv1.ProductsClient
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
		api:    productsv1.NewProductsClient(cc),
		logger: logger,
	}, cc, nil
}
