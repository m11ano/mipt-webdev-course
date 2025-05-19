package productscl

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ProductListItem struct {
	ID                 int64
	IsPublished        bool
	Name               string
	FullDescription    string
	Price              decimal.Decimal
	StockAvailable     int32
	ImagePreviewFileID *uuid.UUID

	ImagePreviewFileURL string

	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
