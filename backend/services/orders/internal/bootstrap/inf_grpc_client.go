package bootstrap

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

var ErrGRPCServerNotConnected = errors.New("cant connect to grpc server")

func ConnectToGRPCServer(ctx context.Context, cc *grpc.ClientConn) error {
	cc.Connect()
	if !cc.WaitForStateChange(ctx, connectivity.Idle) {
		cc.Close()
		return ErrGRPCServerNotConnected
	}

	return nil
}
