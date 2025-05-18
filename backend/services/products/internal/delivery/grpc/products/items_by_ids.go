package productsgrpc

import (
	"context"
	"errors"

	productsv1 "github.com/m11ano/mipt-webdev-course/backend/protos/gen/go/products"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/e"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/usecase"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (s *serverAPI) ItemsByIDs(ctx context.Context, in *productsv1.ItemsByIDsRequest) (*productsv1.ItemsByIDsResponse, error) {
	if len(in.GetIds()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty ids")
	}

	items, err := s.productUC.FindFullList(ctx, usecase.ProductListOptions{
		IDs: lo.ToPtr(in.GetIds()),
	}, nil)
	if err != nil {
		if errors.Is(err, e.ErrInternal) {
			return nil, status.Error(codes.Internal, err.Error())
		}

		return nil, err
	}

	out := &productsv1.ItemsByIDsResponse{
		Items: make([]*productsv1.ProductListItem, len(items)),
	}

	for i, item := range items {
		out.Items[i] = &productsv1.ProductListItem{
			Id:             item.Product.ID,
			Name:           item.Product.Name,
			Price:          item.Product.Price.String(),
			IsPublished:    item.Product.IsPublished,
			StockAvailable: item.Product.StockAvailable,
			CreatedAt:      timestamppb.New(item.Product.CreatedAt),
			UpdatedAt:      toProtoTimestamp(item.Product.UpdatedAt),
			DeletedAt:      toProtoTimestamp(item.Product.DeletedAt),
		}

		if item.ProductPreviewFile != nil {
			out.Items[i].ImagePreviewFileId = wrapperspb.String(item.ProductPreviewFile.ID.String())
			out.Items[i].ImagePreviewFileUrl = item.ProductPreviewFile.GetURL(&s.cfg)
		}
	}

	return out, nil
}
