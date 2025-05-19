package productsgcl

import (
	"log/slog"
	"time"

	"google.golang.org/grpc"

	productscl "github.com/m11ano/mipt-webdev-course/backend/clients/clgrpc/pkg/products"
)

type ClientConn struct {
	Client productscl.Client
	Conn   *grpc.ClientConn
}

func NewClientConn(addr string, retriesCount int, timeout time.Duration, logger *slog.Logger) (*ClientConn, error) {

	client, conn, err := productscl.NewClientImpl(addr, retriesCount, timeout, logger)
	if err != nil {
		return nil, err
	}

	return &ClientConn{
		Client: client,
		Conn:   conn,
	}, nil
}
