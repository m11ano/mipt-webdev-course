package domain

import (
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/e"
	"github.com/shopspring/decimal"
)

var ErrOrderInvalidPrice = e.NewErrorFrom(e.ErrBadRequest).SetMessage("invalid price")
var ErrOrderStockLowerZero = e.NewErrorFrom(e.ErrBadRequest).SetMessage("stock available must be greater than zero")
var ErrOrderStockMoreMax = e.NewErrorFrom(e.ErrBadRequest).SetMessage(fmt.Sprintf("total stock available must be lower than %d", math.MaxInt32))

type Order struct {
	ID                 int64
	IsPublished        bool
	Name               string
	FullDescription    string
	Price              decimal.Decimal
	StockAvailable     int32
	ImagePreviewFileID *uuid.UUID

	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

func NewOrder(id int64) *Order {
	return &Order{
		ID:        id,
		CreatedAt: time.Now(),
	}
}

func (p *Order) SetPrice(price decimal.Decimal) error {
	if price.LessThan(decimal.Zero) {
		return ErrOrderInvalidPrice
	}

	p.Price = price

	return nil
}

func (p *Order) SetStockAvailable(value int32) error {
	if value < 0 {
		return ErrOrderStockLowerZero
	}
	p.StockAvailable = value

	return nil
}

func (p *Order) IncreaseStock(value int64) error {
	newValue := int64(p.StockAvailable) + value
	if newValue > math.MaxInt32 {
		return ErrOrderStockMoreMax
	}

	return p.SetStockAvailable(int32(newValue))
}

func (p *Order) DecreaseStock(value int64) error {
	newValue := int64(p.StockAvailable) - value
	if newValue < 0 {
		return ErrOrderStockLowerZero
	}

	return p.SetStockAvailable(int32(newValue))
}
