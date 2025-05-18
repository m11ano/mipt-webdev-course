package bootstrap

import (
	productsgrpc "github.com/m11ano/mipt-webdev-course/backend/services/products/internal/delivery/grpc/products"
	"go.uber.org/fx"
)

var DeliveryGRPC = fx.Options(
	fx.Provide(NewGRPCServer),
	fx.Invoke(productsgrpc.Register),
)
