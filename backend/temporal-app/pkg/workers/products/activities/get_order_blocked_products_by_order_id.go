package activities

import (
	"context"
	"log/slog"

	productscl "github.com/m11ano/mipt-webdev-course/backend/clients/clgrpc/pkg/products"
	"github.com/m11ano/mipt-webdev-course/backend/temporal-app/pkg/e2temperr"
)

func (c *Controller) GetOrderBlockedProductsByOrderID(ctx context.Context, orderID int64) ([]*productscl.OrderBlockedProduct, error) {

	result, err := c.productsGRPC.Client.GetOrderBlockedProductsByOrderID(ctx, orderID)

	if err != nil {
		c.logger.Error("failed to get order blocked products by order id", slog.Any("error", err.Error()))
	}

	return result, e2temperr.ErrToTempErr(err)
}
