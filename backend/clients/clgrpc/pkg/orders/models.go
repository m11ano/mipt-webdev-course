package orderscl

import (
	"github.com/shopspring/decimal"
)

type OrderCompositionItem struct {
	ProductID int64
	Quantity  int32
	Price     decimal.Decimal
}

type SetOrderCompositionIn struct {
	OrderID       int64
	IsOk          bool
	OrderProducts []OrderCompositionItem
}
