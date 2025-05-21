package productsgrpc

import (
	"context"

	"github.com/m11ano/e"
	productsv1 "github.com/m11ano/mipt-webdev-course/backend/protos/gen/go/products"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverAPI) SetOrderBlockedProductsByOrderID(ctx context.Context, in *productsv1.SetOrderBlockedProductsByOrderIDRequest) (*productsv1.SetOrderBlockedProductsByOrderIDResponse, error) {

	if in.GetOrderId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty order_id")
	}

	composition := make([]usecase.ProductOrderBlockComposition, len(in.GetItems()))

	for i, item := range in.GetItems() {
		composition[i] = usecase.ProductOrderBlockComposition{
			ProductID: item.GetProductId(),
			Quantity:  item.GetQuantity(),
		}
	}

	err := s.productUC.SetOrderBlock(ctx, in.GetOrderId(), composition)
	if err != nil {
		if isAppErr, appErr := e.IsAppError(err); isAppErr {
			return nil, appErr.AsGRPCError()
		}
		return nil, err
	}

	return nil, nil
}
