package bootstrap

import (
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/repository"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/usecase"
	"go.uber.org/fx"
)

var ProductModule = fx.Module(
	"product_module",
	fx.Provide(
		fx.Private,
		fx.Annotate(repository.NewProduct, fx.As(new(usecase.ProductRepository))),
	),
	fx.Provide(
		fx.Annotate(usecase.NewProductInpl, fx.As(new(usecase.Product))),
	),
)
