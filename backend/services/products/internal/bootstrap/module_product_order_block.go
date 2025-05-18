package bootstrap

import (
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/repository"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/usecase"
	"go.uber.org/fx"
)

var ProductOrderBlockModule = fx.Module(
	"product_order_block_module",
	fx.Provide(
		fx.Private,
		fx.Annotate(repository.NewProductOrderBlock, fx.As(new(usecase.ProductOrderBlockRepository))),
	),
	fx.Provide(
		fx.Annotate(usecase.NewProductOrderBlockInpl, fx.As(new(usecase.ProductOrderBlock))),
	),
)
