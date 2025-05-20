package ordersgrpc

import (
	"time"

	ordersv1 "github.com/m11ano/mipt-webdev-course/backend/protos/gen/go/orders"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/infra/config"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type serverAPI struct {
	ordersv1.UnimplementedOrdersServer
	cfg config.Config
	// productUC           usecase.Product
	// productOrderBlockUC usecase.ProductOrderBlock
}

func Register(gRPCServer *grpc.Server, cfg config.Config) {
	ordersv1.RegisterOrdersServer(gRPCServer, &serverAPI{
		cfg: cfg,
	})
}

func toProtoTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}
