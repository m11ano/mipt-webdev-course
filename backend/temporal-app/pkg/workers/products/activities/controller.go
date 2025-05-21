package activities

import (
	"log/slog"

	ordersgcl "github.com/m11ano/mipt-webdev-course/backend/temporal-app/internal/clients/grpc/orders"
	productsgcl "github.com/m11ano/mipt-webdev-course/backend/temporal-app/internal/clients/grpc/products"
)

type Controller struct {
	logger       *slog.Logger
	productsGRPC *productsgcl.ClientConn
	ordersGRPC   *ordersgcl.ClientConn
}

func NewController(logger *slog.Logger, productsGRPC *productsgcl.ClientConn, ordersGRPC *ordersgcl.ClientConn) *Controller {
	return &Controller{
		logger:       logger,
		productsGRPC: productsGRPC,
		ordersGRPC:   ordersGRPC,
	}
}
