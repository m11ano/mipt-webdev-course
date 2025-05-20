package ordersgcl

import (
	"log/slog"
	"time"

	"google.golang.org/grpc"

	orderscl "github.com/m11ano/mipt-webdev-course/backend/clients/clgrpc/pkg/orders"
)

type ClientConn struct {
	Client orderscl.Client
	Conn   *grpc.ClientConn
}

func NewClientConn(addr string, retriesCount int, timeout time.Duration, logger *slog.Logger) (*ClientConn, error) {

	client, conn, err := orderscl.NewClientImpl(addr, retriesCount, timeout, logger)
	if err != nil {
		return nil, err
	}

	return &ClientConn{
		Client: client,
		Conn:   conn,
	}, nil
}
