package activities

import (
	"context"
	"log/slog"

	orderscl "github.com/m11ano/mipt-webdev-course/backend/clients/clgrpc/pkg/orders"
	"github.com/m11ano/mipt-webdev-course/backend/temporal-app/pkg/e2temperr"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

type InformOrdersServiceAboutOrderCompositionItem struct {
	ProductID int64
	Quantity  int32
	Price     decimal.Decimal
}

type InformOrdersServiceAboutOrderCompositionIn struct {
	OrderID       int64
	IsOk          bool
	OrderProducts *[]InformOrdersServiceAboutOrderCompositionItem
	OrderStatus   *string
}

func (c *Controller) InformOrdersServiceAboutOrderComposition(ctx context.Context, input InformOrdersServiceAboutOrderCompositionIn) error {

	req := orderscl.SetOrderCompositionIn{
		OrderID:     input.OrderID,
		IsOk:        input.IsOk,
		OrderStatus: input.OrderStatus,
	}

	if input.OrderProducts != nil {
		req.OrderProducts = lo.ToPtr(lo.Map(*input.OrderProducts, func(item InformOrdersServiceAboutOrderCompositionItem, _ int) orderscl.OrderCompositionItem {
			return orderscl.OrderCompositionItem{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Price,
			}
		}))
	}

	err := c.ordersGRPC.Client.SetOrderComposition(ctx, req)
	if err != nil {
		c.logger.Error("failed to inform orders service about order composition", slog.Any("error", err.Error()))

		return e2temperr.ErrToTempErr(err)
	}

	return nil
}
