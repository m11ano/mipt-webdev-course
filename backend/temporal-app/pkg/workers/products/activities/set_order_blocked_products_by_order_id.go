package activities

import (
	"context"

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
		return e2temperr.ErrToTempErr(err)
	}

	return nil
}
