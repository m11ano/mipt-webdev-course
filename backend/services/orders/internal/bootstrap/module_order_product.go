package bootstrap

import (
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/repository"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/usecase"
	"go.uber.org/fx"
)

var OrderProductModule = fx.Module(
	"order_product_module",
	fx.Provide(
		fx.Private,
		fx.Annotate(repository.NewOrderProduct, fx.As(new(usecase.OrderProductRepository))),
	),
	fx.Provide(
		fx.Annotate(usecase.NewOrderProductInpl, fx.As(new(usecase.OrderProduct))),
	),
)
