package controller

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/infra/config"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/usecase"
)

type Controller struct {
	logger  *slog.Logger
	vldtr   *validator.Validate
	cfg     config.Config
	orderUC usecase.Order
}

func New(logger *slog.Logger, vldtr *validator.Validate, cfg config.Config, orderUC usecase.Order) *Controller {
	return &Controller{
		logger:  logger,
		vldtr:   vldtr,
		cfg:     cfg,
		orderUC: orderUC,
	}
}
