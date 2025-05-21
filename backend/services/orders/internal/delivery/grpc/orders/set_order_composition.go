package ordersgrpc

import (
	"context"

	"github.com/m11ano/e"
	ordersv1 "github.com/m11ano/mipt-webdev-course/backend/protos/gen/go/orders"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/domain"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/usecase"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverAPI) SetOrderComposition(ctx context.Context, in *ordersv1.SetOrderCompositionRequest) (*ordersv1.SetOrderCompositionResponse, error) {

	if in.GetOrderId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty order_id")
	}

	if in.GetIsOk() {
		params := usecase.SetOrderCompositionIn{
			OrderID: in.GetOrderId(),
		}

		switch v := in.GetOptionalProducts().(type) {
		case *ordersv1.SetOrderCompositionRequest_ItemsSet:
			products := make([]usecase.OrderProductWithPrice, len(v.ItemsSet.Items))
			for i, item := range v.ItemsSet.Items {
				price, err := decimal.NewFromString(item.GetPrice())
				if err != nil {
					return nil, e.ErrBadRequest.Wrap(err).AsGRPCError()
				}

				products[i] = usecase.OrderProductWithPrice{
					ID:       item.GetProductId(),
					Quantity: item.GetQuantity(),
					Price:    price,
				}
			}

			params.Products = &products

		case *ordersv1.SetOrderCompositionRequest_NoItems:
		default:

		}

		if in.GetOrderStatus() != nil {
			status, ok := domain.OrderStatusMap[in.GetOrderStatus().GetValue()]
			if !ok {
				return nil, e.ErrBadRequest.AsGRPCError()
			}

			params.Status = &status
		}

		err := s.orderUC.SetOrderComposition(ctx, params)

		if err != nil {
			if isAppErr, appErr := e.IsAppError(err); isAppErr {
				return nil, appErr.AsGRPCError()
			}
			return nil, err
		}
	} else {

		err := s.orderUC.RemoveOrderIfNew(ctx, in.GetOrderId())
		if err != nil {
			if isAppErr, appErr := e.IsAppError(err); isAppErr {
				return nil, appErr.AsGRPCError()
			}
			return nil, err
		}
	}

	return nil, nil
}
