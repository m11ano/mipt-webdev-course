package bootstrap

import (
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/repository"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/usecase"
	"go.uber.org/fx"
)

var OrderModule = fx.Module(
	"order_module",
	fx.Provide(
		fx.Private,
		fx.Annotate(repository.NewOrder, fx.As(new(usecase.OrderRepository))),
	),
	fx.Provide(
		fx.Annotate(usecase.NewOrderInpl, fx.As(new(usecase.Order))),
	),
)
