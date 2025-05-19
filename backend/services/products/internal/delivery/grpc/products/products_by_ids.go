package productsgrpc

import (
	"context"

	"github.com/m11ano/e"
	productsv1 "github.com/m11ano/mipt-webdev-course/backend/protos/gen/go/products"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/usecase"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (s *serverAPI) ProductsByIDs(ctx context.Context, in *productsv1.ProductsByIDsRequest) (*productsv1.ProductsByIDsResponse, error) {
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

	out := &productsv1.ProductsByIDsResponse{
		Items: make([]*productsv1.ProductListItem, len(items)),
	}

	for i, item := range items {
		out.Items[i] = &productsv1.ProductListItem{
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

	return out, nil
}
