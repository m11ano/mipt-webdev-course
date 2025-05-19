package domain

import (
	"time"

	"github.com/m11ano/e"
)

var ErrProductOrderBlockQuantityLess1 = e.NewErrorFrom(e.ErrBadRequest).SetMessage("invalid quantity")

type ProductOrderBlock struct {
	ProductID int64
	OrderID   int64
	Quantity  int32

	CreatedAt time.Time
}

func NewProductOrderBlock(productID int64, OrderID int64, Quantity int32) (*ProductOrderBlock, error) {
	item := &ProductOrderBlock{
		ProductID: productID,
		OrderID:   OrderID,
		CreatedAt: time.Now(),
	}

	err := item.SetQuantity(Quantity)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (pb *ProductOrderBlock) SetQuantity(quantity int32) error {
	if quantity < 1 {
		return ErrProductOrderBlockQuantityLess1
	}
	return nil
}
