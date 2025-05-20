package activities

import (
	ordersgcl "github.com/m11ano/mipt-webdev-course/backend/temporal-app/internal/clients/grpc/orders"
	productsgcl "github.com/m11ano/mipt-webdev-course/backend/temporal-app/internal/clients/grpc/products"
)

type Controller struct {
	productsGRPC *productsgcl.ClientConn
	ordersGRPC   *ordersgcl.ClientConn
}

func NewController(productsGRPC *productsgcl.ClientConn, ordersGRPC *ordersgcl.ClientConn) *Controller {
	return &Controller{
		productsGRPC: productsGRPC,
		ordersGRPC:   ordersGRPC,
	}
}
