package domain

import (
	"time"

	"github.com/m11ano/e"
	"github.com/shopspring/decimal"
)

var ErrOrderProductQuantityLess1 = e.NewErrorFrom(e.ErrBadRequest).SetMessage("invalid quantity")
var ErrOrderProductPriceLess1 = e.NewErrorFrom(e.ErrBadRequest).SetMessage("invalid price")

type OrderProduct struct {
	ProductID int64
	OrderID   int64
	Quantity  int32
	Price     decimal.Decimal

	CreatedAt time.Time
}

func NewOrderProduct(orderID int64, productID int64, quantity int32, price decimal.Decimal) (*OrderProduct, error) {
	item := &OrderProduct{
		ProductID: productID,
		OrderID:   orderID,
		CreatedAt: time.Now(),
	}

	err := item.SetQuantity(quantity)
	if err != nil {
		return nil, err
	}

	err = item.SetPrice(price)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (op *OrderProduct) SetQuantity(quantity int32) error {
	if quantity < 1 {
		return ErrOrderProductQuantityLess1
	}
	op.Quantity = quantity

	return nil
}

func (op *OrderProduct) SetPrice(price decimal.Decimal) error {
	if price.LessThan(decimal.Zero) {
		return ErrOrderProductPriceLess1
	}
	op.Price = price

	return nil
}
