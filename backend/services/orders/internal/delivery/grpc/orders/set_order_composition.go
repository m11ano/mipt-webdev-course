package ordersgrpc

import (
	"context"

	"github.com/m11ano/e"
	ordersv1 "github.com/m11ano/mipt-webdev-course/backend/protos/gen/go/orders"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/usecase"
	"github.com/shopspring/decimal"
)

func (s *serverAPI) SetOrderComposition(ctx context.Context, in *ordersv1.SetOrderCompositionRequest) (*ordersv1.SetOrderCompositionResponse, error) {

	if in.GetIsOk() {
		products := make([]usecase.OrderProductWithPriceIn, len(in.GetItems()))
		for i, item := range in.GetItems() {
			price, err := decimal.NewFromString(item.GetPrice())
			if err != nil {
				return nil, e.ErrBadRequest.Wrap(err).AsGRPCError()
			}

			products[i] = usecase.OrderProductWithPriceIn{
				ID:       item.GetProductId(),
				Quantity: item.GetQuantity(),
				Price:    price,
			}
		}

		err := s.orderUC.SetOrderComposition(ctx, usecase.SetOrderCompositionIn{
			OrderID:  in.GetOrderId(),
			Products: products,
		})

		if err != nil {
			if isAppErr, appErr := e.IsAppError(err); isAppErr {
				return nil, appErr.AsGRPCError()
			}
			return nil, err
		}
	} else {

		err := s.orderUC.RemoveNewOrder(ctx, in.GetOrderId())
		if err != nil {
			if isAppErr, appErr := e.IsAppError(err); isAppErr {
				return nil, appErr.AsGRPCError()
			}
			return nil, err
		}
	}

	return nil, nil
}
