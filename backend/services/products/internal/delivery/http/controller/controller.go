package controller

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/infra/config"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/usecase"
)

type Controller struct {
	logger    *slog.Logger
	vldtr     *validator.Validate
	cfg       config.Config
	fileUC    usecase.File
	productUC usecase.Product
}

func New(logger *slog.Logger, vldtr *validator.Validate, cfg config.Config, fileUC usecase.File, productUC usecase.Product) *Controller {
	return &Controller{
		logger:    logger,
		vldtr:     vldtr,
		cfg:       cfg,
		fileUC:    fileUC,
		productUC: productUC,
	}
}
