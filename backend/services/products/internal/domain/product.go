package domain

import (
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/m11ano/e"
	"github.com/shopspring/decimal"
)

var ErrProductInvalidPrice = e.NewErrorFrom(e.ErrBadRequest).SetMessage("invalid price")
var ErrProductStockLowerZero = e.NewErrorFrom(e.ErrBadRequest).SetMessage("stock available must be greater than zero")
var ErrProductStockMoreMax = e.NewErrorFrom(e.ErrBadRequest).SetMessage(fmt.Sprintf("total stock available must be lower than %d", math.MaxInt32))

type Product struct {
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

func NewProduct(id int64) *Product {
	return &Product{
		ID:        id,
		CreatedAt: time.Now(),
	}
}

func (p *Product) SetPrice(price decimal.Decimal) error {
	if price.LessThan(decimal.Zero) {
		return ErrProductInvalidPrice
	}

	p.Price = price

	return nil
}

func (p *Product) SetStockAvailable(value int32) error {
	if value < 0 {
		return ErrProductStockLowerZero
	}
	p.StockAvailable = value

	return nil
}

func (p *Product) IncreaseStock(value int64) error {
	newValue := int64(p.StockAvailable) + value
	if newValue > math.MaxInt32 {
		return ErrProductStockMoreMax
	}

	return p.SetStockAvailable(int32(newValue))
}

func (p *Product) DecreaseStock(value int64) error {
	newValue := int64(p.StockAvailable) - value
	if newValue < 0 {
		return ErrProductStockLowerZero
	}

	return p.SetStockAvailable(int32(newValue))
}
