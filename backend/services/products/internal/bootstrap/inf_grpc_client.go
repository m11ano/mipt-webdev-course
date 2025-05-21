package bootstrap

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

var ErrGRPCServerNotCreated = errors.New("cant create connect to grpc server")

func ConnectToGRPCServer(ctx context.Context, cc *grpc.ClientConn) error {
	cc.Connect()
	if !cc.WaitForStateChange(ctx, connectivity.Idle) {
		cc.Close()
		return ErrGRPCServerNotCreated
	}

	return nil
}
