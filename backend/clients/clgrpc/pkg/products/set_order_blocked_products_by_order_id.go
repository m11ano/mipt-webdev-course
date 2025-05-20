package productscl

import (
	"context"

	"github.com/m11ano/e"
	productsv1 "github.com/m11ano/mipt-webdev-course/backend/protos/gen/go/products"
	"github.com/samber/lo"
)

func (c *ClientImpl) SetOrderBlockedProductsByOrderID(ctx context.Context, in SetOrderBlockedProductsByOrderIDIn) error {

	_, err := c.api.SetOrderBlockedProductsByOrderID(ctx, &productsv1.SetOrderBlockedProductsByOrderIDRequest{
		OrderId: in.OrderID,
		Items: lo.Map(in.OrderProducts, func(item OrderBlockedProduct, _ int) *productsv1.OrderProduct {
			return &productsv1.OrderProduct{
				ProductId: item.ProductID,
				Quantity:  item.Quantity,
			}
		}),
	})

	if err != nil {
		if ok, lgErr := e.ErrConvertGRPCToLogic(err); ok {
			return lgErr
		}

		return err
	}

	return nil
}
