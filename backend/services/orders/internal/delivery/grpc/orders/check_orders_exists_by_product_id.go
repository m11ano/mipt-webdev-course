package ordersgrpc

import (
	"context"

	ordersv1 "github.com/m11ano/mipt-webdev-course/backend/protos/gen/go/orders"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/usecase"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/usecase/uctypes"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverAPI) CheckOrdersExistsByProductID(ctx context.Context, in *ordersv1.CheckOrdersExistsByProductIDRequest) (*ordersv1.CheckOrdersExistsByProductIDResponse, error) {

	if in.GetProductId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty product_id")
	}

	check, err := s.orderProductUC.FindList(ctx, usecase.OrderProductListOptions{
		ProductID: lo.ToPtr(in.GetProductId()),
	}, &uctypes.QueryGetListParams{
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}

	out := &ordersv1.CheckOrdersExistsByProductIDResponse{
		Exists: len(check) > 0,
	}

	return out, nil
}
