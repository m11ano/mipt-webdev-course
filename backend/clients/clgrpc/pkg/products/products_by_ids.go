package productscl

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/m11ano/avito-pvz/pkg/e"
	productsv1 "github.com/m11ano/mipt-webdev-course/backend/protos/gen/go/products"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

func (c *ClientImpl) GetProductsByIds(ctx context.Context, ids []int64) ([]ProductListItem, error) {
	items, err := c.api.ProductsByIDs(ctx, &productsv1.ProductsByIDsRequest{Ids: ids})
	if err != nil {
		return nil, err
	}

	result := make([]ProductListItem, len(items.GetItems()))
	for i, item := range items.GetItems() {
		price, err := decimal.NewFromString(item.GetPrice())
		if err != nil {
			c.logger.ErrorContext(ctx, "converting price", slog.Any("price", item.GetPrice()), slog.Any("error", err))
			return nil, e.ErrInternal.Wrap(err)
		}

		if item.GetCreatedAt() == nil {
			c.logger.ErrorContext(ctx, "createdAt in nil", slog.Any("error", err))
			return nil, e.ErrInternal.Wrap(err)
		}

		result[i] = ProductListItem{
			ID:                  item.GetId(),
			IsPublished:         item.GetIsPublished(),
			Name:                item.GetName(),
			FullDescription:     item.GetFullDescription(),
			Price:               price,
			StockAvailable:      item.GetStockAvailable(),
			ImagePreviewFileURL: item.GetImagePreviewFileUrl(),
			CreatedAt:           item.GetCreatedAt().AsTime(),
		}

		if item.GetUpdatedAt() != nil {
			result[i].UpdatedAt = lo.ToPtr(item.GetUpdatedAt().AsTime())
		}

		if item.GetDeletedAt() != nil {
			result[i].DeletedAt = lo.ToPtr(item.GetDeletedAt().AsTime())
		}

		if item.GetImagePreviewFileId() != nil {
			uuid, err := uuid.Parse(item.GetImagePreviewFileId().GetValue())
			if err != nil {
				c.logger.ErrorContext(ctx, "converting uuid", slog.Any("uuid", item.GetImagePreviewFileId().GetValue()), slog.Any("error", err))
				return nil, e.ErrInternal.Wrap(err)
			}

			result[i].ImagePreviewFileID = &uuid
		}
	}

	return result, nil
}
