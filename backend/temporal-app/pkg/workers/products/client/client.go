package productstc

import (
	"context"

	"github.com/shopspring/decimal"
)

type SetOrderProductsAndStatusIn struct {
	NotWait       bool
	OrderID       int64
	OrderProducts *[]OrderProductsItem
	OrderStatus   *string
}

type OrderProductsItem struct {
	ProductID int64
	Quantity  int32
	Price     decimal.Decimal
}

type Client interface {
	SetOrderProductsAndStatus(ctx context.Context, input SetOrderProductsAndStatusIn) error
}

var WorkflowOrderProductsPrefx = "order_products"
