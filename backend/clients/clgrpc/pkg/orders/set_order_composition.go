package orderscl

import (
	"context"

	"github.com/m11ano/e"
	ordersv1 "github.com/m11ano/mipt-webdev-course/backend/protos/gen/go/orders"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (c *ClientImpl) SetOrderComposition(ctx context.Context, in SetOrderCompositionIn) error {

	req := &ordersv1.SetOrderCompositionRequest{
		OrderId: in.OrderID,
		IsOk:    in.IsOk,
		OptionalProducts: &ordersv1.SetOrderCompositionRequest_NoItems{
			NoItems: &emptypb.Empty{},
		},
	}

	if in.OrderStatus != nil {
		req.OrderStatus = wrapperspb.String(*in.OrderStatus)
	}

	if in.OrderProducts != nil {
		req.OptionalProducts = &ordersv1.SetOrderCompositionRequest_ItemsSet{
			ItemsSet: &ordersv1.OrderProductList{
				Items: lo.Map(*in.OrderProducts, func(item OrderCompositionItem, _ int) *ordersv1.OrderProduct {
					return &ordersv1.OrderProduct{
						ProductId: item.ProductID,
						Quantity:  item.Quantity,
						Price:     item.Price.String(),
					}
				}),
			},
		}
	}

	_, err := c.api.SetOrderComposition(ctx, req)

	if err != nil {
		if ok, lgErr := e.ErrConvertGRPCToLogic(err); ok {
			return lgErr
		}

		return err
	}

	return nil
}
