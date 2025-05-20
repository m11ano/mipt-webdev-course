package bootstrap

import (
	ordersgrpc "github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/delivery/grpc/orders"
	"go.uber.org/fx"
)

var DeliveryGRPC = fx.Options(
	fx.Provide(NewGRPCServer),
	fx.Invoke(ordersgrpc.Register),
)
