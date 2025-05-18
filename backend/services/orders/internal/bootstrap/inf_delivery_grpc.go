package bootstrap

import (
	"go.uber.org/fx"
)

var DeliveryGRPC = fx.Options(
	fx.Provide(NewGRPCServer),
	//fx.Invoke(ordersgrpc.Register),
)
