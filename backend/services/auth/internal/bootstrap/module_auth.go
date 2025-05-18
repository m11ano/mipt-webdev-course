package bootstrap

import (
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/internal/usecase"
	"go.uber.org/fx"
)

var AuthModule = fx.Module(
	"auth_module",
	fx.Provide(
		fx.Annotate(usecase.NewAuthInpl, fx.As(new(usecase.Auth))),
	),
)
