package bootstrap

import (
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/repository"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/usecase"
	"go.uber.org/fx"
)

var ProductSliderImageModule = fx.Module(
	"product_slider_image_module",
	fx.Provide(
		fx.Private,
		fx.Annotate(repository.NewProductSliderImage, fx.As(new(usecase.ProductSliderImageRepository))),
	),
	fx.Provide(
		fx.Annotate(usecase.NewProductSliderImageInpl, fx.As(new(usecase.ProductSliderImage))),
	),
)
