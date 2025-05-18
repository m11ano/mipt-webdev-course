package controller

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/internal/infra/config"
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/internal/usecase"
)

type Controller struct {
	logger    *slog.Logger
	vldtr     *validator.Validate
	cfg       config.Config
	accountUC usecase.Account
	authUC    usecase.Auth
}

func New(logger *slog.Logger, vldtr *validator.Validate, cfg config.Config, accountUC usecase.Account, authUC usecase.Auth) *Controller {
	return &Controller{
		logger:    logger,
		vldtr:     vldtr,
		cfg:       cfg,
		accountUC: accountUC,
		authUC:    authUC,
	}
}
