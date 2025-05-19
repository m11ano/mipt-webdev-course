package activities

import productsgcl "github.com/m11ano/mipt-webdev-course/backend/temporal-app/internal/clients/grpc/products"

type Controller struct {
	productsGRPC *productsgcl.ClientConn
}

func NewController(productsGRPC *productsgcl.ClientConn) *Controller {
	return &Controller{
		productsGRPC: productsGRPC,
	}
}
