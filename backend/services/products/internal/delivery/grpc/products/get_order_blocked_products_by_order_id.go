package productsgrpc

import (
	"context"

	"github.com/m11ano/e"
	productsv1 "github.com/m11ano/mipt-webdev-course/backend/protos/gen/go/products"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverAPI) GetOrderBlockedProductsByOrderID(ctx context.Context, in *productsv1.GetOrderBlockedProductsByOrderIDRequest) (*productsv1.GetOrderBlockedProductsByOrderIDResponse, error) {
	if in.GetOrderId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty order_id")
	}

	items, err := s.productOrderBlockUC.GetOrderBlockedProducts(ctx, in.GetOrderId())
	if err != nil {
		if isAppErr, appErr := e.IsAppError(err); isAppErr {
			return nil, appErr.AsGRPCError()
		}
		return nil, err
	}

	out := &productsv1.GetOrderBlockedProductsByOrderIDResponse{
		Items: make([]*productsv1.OrderBlockedProduct, len(items)),
	}

	for i, item := range items {
		out.Items[i] = &productsv1.OrderBlockedProduct{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	return out, nil
}
