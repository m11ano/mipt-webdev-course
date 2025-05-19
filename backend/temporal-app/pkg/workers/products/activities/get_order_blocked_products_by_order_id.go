package activities

import (
	"context"

	productscl "github.com/m11ano/mipt-webdev-course/backend/clients/clgrpc/pkg/products"
)

func (c *Controller) GetOrderBlockedProductsByOrderID(ctx context.Context, orderID int64) ([]*productscl.OrderBlockedProduct, error) {

	return c.productsGRPC.Client.GetOrderBlockedProductsByOrderID(ctx, orderID)
}
