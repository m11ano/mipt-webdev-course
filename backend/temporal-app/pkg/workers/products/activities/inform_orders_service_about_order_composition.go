package activities

import (
	"context"

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
	OrderProducts []InformOrdersServiceAboutOrderCompositionItem
}

func (c *Controller) InformOrdersServiceAboutOrderComposition(ctx context.Context, input InformOrdersServiceAboutOrderCompositionIn) error {

	err := c.ordersGRPC.Client.SetOrderComposition(ctx, orderscl.SetOrderCompositionIn{
		OrderID: input.OrderID,
		IsOk:    input.IsOk,
		OrderProducts: lo.Map(input.OrderProducts, func(item InformOrdersServiceAboutOrderCompositionItem, _ int) orderscl.OrderCompositionItem {
			return orderscl.OrderCompositionItem{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Price,
			}
		}),
	})
	if err != nil {
		return e2temperr.ErrToTempErr(err)
	}

	return nil
}
