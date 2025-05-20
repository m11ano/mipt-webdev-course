package activities

import (
	"context"

	productscl "github.com/m11ano/mipt-webdev-course/backend/clients/clgrpc/pkg/products"
	"github.com/m11ano/mipt-webdev-course/backend/temporal-app/pkg/e2temperr"
)

func (c *Controller) GetOrderBlockedProductsByOrderID(ctx context.Context, orderID int64) ([]*productscl.OrderBlockedProduct, error) {

	result, err := c.productsGRPC.Client.GetOrderBlockedProductsByOrderID(ctx, orderID)

	return result, e2temperr.ErrToTempErr(err)
}
