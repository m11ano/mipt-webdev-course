package bootstrap

import (
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/repository"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/usecase"
	"go.uber.org/fx"
)

var FileModule = fx.Module(
	"file_module",
	fx.Provide(
		fx.Private,
		fx.Annotate(repository.NewFile, fx.As(new(usecase.FileRepository))),
	),
	fx.Provide(
		fx.Annotate(usecase.NewFileInpl, fx.As(new(usecase.File))),
	),
)
