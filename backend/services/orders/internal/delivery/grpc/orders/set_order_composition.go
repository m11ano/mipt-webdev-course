package ordersgrpc

import (
	"context"
	"fmt"

	ordersv1 "github.com/m11ano/mipt-webdev-course/backend/protos/gen/go/orders"
)

func (s *serverAPI) SetOrderComposition(ctx context.Context, in *ordersv1.SetOrderCompositionRequest) (*ordersv1.SetOrderCompositionResponse, error) {
	fmt.Println("запрос пришел")

	fmt.Println(in.GetOrderId(), in.GetIsOk())

	for _, item := range in.GetItems() {
		fmt.Println(item.GetProductId(), item.GetQuantity(), item.GetPrice())
	}

	//return nil, e.ErrInternal.AsGRPCError()

	return nil, nil

	/*
		if len(in.GetIds()) == 0 {
			return nil, status.Error(codes.InvalidArgument, "empty ids")
		}

		items, err := s.productUC.FindFullList(ctx, usecase.ProductListOptions{
			IDs: lo.ToPtr(in.GetIds()),
		}, nil)
		if err != nil {
			if isAppErr, appErr := e.IsAppError(err); isAppErr {
				return nil, status.Error(appErr.GetGRPCCode(), err.Error())
			}
			return nil, err
		}

		out := &ordersv1.GetProductsByIDsResponse{
			Items: make([]*ordersv1.ProductListItem, len(items)),
		}

		for i, item := range items {
			out.Items[i] = &ordersv1.ProductListItem{
				Id:              item.Product.ID,
				Name:            item.Product.Name,
				FullDescription: item.Product.FullDescription,
				Price:           item.Product.Price.String(),
				IsPublished:     item.Product.IsPublished,
				StockAvailable:  item.Product.StockAvailable,
				CreatedAt:       timestamppb.New(item.Product.CreatedAt),
				UpdatedAt:       toProtoTimestamp(item.Product.UpdatedAt),
				DeletedAt:       toProtoTimestamp(item.Product.DeletedAt),
			}

			if item.ProductPreviewFile != nil {
				out.Items[i].ImagePreviewFileId = wrapperspb.String(item.ProductPreviewFile.ID.String())
				out.Items[i].ImagePreviewFileUrl = item.ProductPreviewFile.GetURL(&s.cfg)
			}
		}
	*/

	return nil, nil
}
