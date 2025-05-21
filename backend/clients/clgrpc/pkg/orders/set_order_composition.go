package orderscl

import (
	"context"

	"github.com/m11ano/e"
	ordersv1 "github.com/m11ano/mipt-webdev-course/backend/protos/gen/go/orders"
	"github.com/samber/lo"
)

func (c *ClientImpl) SetOrderComposition(ctx context.Context, in SetOrderCompositionIn) error {

	_, err := c.api.SetOrderComposition(ctx, &ordersv1.SetOrderCompositionRequest{
		OrderId: in.OrderID,
		IsOk:    in.IsOk,
		Items: lo.Map(in.OrderProducts, func(item OrderCompositionItem, _ int) *ordersv1.OrderProduct {
			return &ordersv1.OrderProduct{
				ProductId: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Price.String(),
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
