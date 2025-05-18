package bootstrap

import (
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/internal/repository"
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/internal/usecase"
	"go.uber.org/fx"
)

var AccountModule = fx.Module(
	"account_module",
	fx.Provide(
		fx.Private,
		fx.Annotate(repository.NewAccount, fx.As(new(usecase.AccountRepository))),
	),
	fx.Provide(
		fx.Annotate(usecase.NewAccountInpl, fx.As(new(usecase.Account))),
	),
)
