package bootstrap

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

var ErrGRPCServerNotConnected = errors.New("cant connect to grpc server")

func ConnectToGRPCServer(ctx context.Context, c1 *grpc.ClientConn, c2 *grpc.ClientConn) error {
	c1.Connect()
	if !c1.WaitForStateChange(ctx, connectivity.Idle) {
		c1.Close()
		return ErrGRPCServerNotConnected
	}

	c2.Connect()
	if !c2.WaitForStateChange(ctx, connectivity.Idle) {
		c2.Close()
		return ErrGRPCServerNotConnected
	}

	return nil
}
