package productstc

import (
	"context"

	"github.com/shopspring/decimal"
)

type SetOrderProductsAndBlockIn struct {
	OrderID       int64
	OrderProducts []OrderProductsItem
}

type OrderProductsItem struct {
	ProductID int64
	Quantity  int32
	Price     decimal.Decimal
}

type Client interface {
	SetOrderProductsAndBlock(ctx context.Context, input SetOrderProductsAndBlockIn) error
}

var WorkflowOrderProductsPrefx = "order_products"
