package domain

import (
	"time"

	"github.com/google/uuid"
)

type ProductSliderImage struct {
	ProductID int64
	FileID    uuid.UUID
	Sort      int32

	CreatedAt time.Time
}

func NewProductSliderImage(productID int64, fileID uuid.UUID) *ProductSliderImage {
	return &ProductSliderImage{
		ProductID: productID,
		FileID:    fileID,
		CreatedAt: time.Now(),
	}
}
