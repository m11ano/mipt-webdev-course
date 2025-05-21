package activities

import (
	"context"
	"log/slog"

	productscl "github.com/m11ano/mipt-webdev-course/backend/clients/clgrpc/pkg/products"
	"github.com/m11ano/mipt-webdev-course/backend/temporal-app/pkg/e2temperr"
	"github.com/samber/lo"
)

type SetOrderBlockedProductsByOrderIDItem struct {
	ProductID int64
	Quantity  int32
}

type SetOrderBlockedProductsByOrderIDIn struct {
	OrderID       int64
	OrderProducts []SetOrderBlockedProductsByOrderIDItem
}

func (c *Controller) SetOrderBlockedProductsByOrderID(ctx context.Context, in SetOrderBlockedProductsByOrderIDIn) error {

	err := c.productsGRPC.Client.SetOrderBlockedProductsByOrderID(ctx, productscl.SetOrderBlockedProductsByOrderIDIn{
		OrderID: in.OrderID,
		OrderProducts: lo.Map(in.OrderProducts, func(item SetOrderBlockedProductsByOrderIDItem, _ int) productscl.OrderBlockedProduct {
			return productscl.OrderBlockedProduct{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
			}
		}),
	})
	if err != nil {
		c.logger.Error("failed to set order blocked products by order id", slog.Any("error", err.Error()))

		return e2temperr.ErrToTempErr(err)
	}

	return nil
}
