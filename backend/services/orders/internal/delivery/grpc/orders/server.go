package ordersgrpc

import (
	ordersv1 "github.com/m11ano/mipt-webdev-course/backend/protos/gen/go/orders"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/infra/config"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/usecase"
	"google.golang.org/grpc"
)

type serverAPI struct {
	ordersv1.UnimplementedOrdersServer
	cfg            config.Config
	orderUC        usecase.Order
	orderProductUC usecase.OrderProduct
}

func Register(gRPCServer *grpc.Server, cfg config.Config, orderUC usecase.Order, orderProductUC usecase.OrderProduct) {
	ordersv1.RegisterOrdersServer(gRPCServer, &serverAPI{
		cfg:            cfg,
		orderUC:        orderUC,
		orderProductUC: orderProductUC,
	})
}
