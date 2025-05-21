package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/delivery/http/middleware"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/usecase"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/usecase/uctypes"
	"github.com/samber/lo"
)

type GetProductsOutItem struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	Price          float64 `json:"price"`
	StockAvailable int32   `json:"stock_available"`
	ImagePreview   string  `json:"image_preview"`
}

type GetProductsOut struct {
	Items []GetProductsOutItem `json:"items"`
	Total int64                `json:"total"`
}

// @Summary Получить список продуктов
// @Tags products
// @Produce  json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} GetProductsOut
// @Failure 400 {object} middleware.ErrorJSON
// @Router /products [get]
func (ctrl *Controller) GetProductsHandler(c *fiber.Ctx) error {

	limit := c.QueryInt("limit", 20)
	if limit > 100 {
		limit = 100
	}
	if limit < 1 {
		limit = 1
	}

	offset := c.QueryInt("offset", 0)
	if offset < 0 {
		offset = 0
	}

	authData := middleware.ExtractAuthData(c)

	var onlyIsPublished *bool

	if !authData.IsAuth {
		onlyIsPublished = lo.ToPtr(true)
	}

	data, total, err := ctrl.productUC.FindFullPagedList(c.Context(), usecase.ProductListOptions{
		IsPublished: onlyIsPublished,
		Sort: &[]usecase.ProductListSort{
			{
				Field:  usecase.ProductListSortFieldCreatedAt,
				IsDesc: true,
			},
		},
	}, &uctypes.QueryGetListParams{
		Limit:  uint64(limit),
		Offset: uint64(offset),
	})
	if err != nil {
		return err
	}

	result := GetProductsOut{
		Items: make([]GetProductsOutItem, len(data)),
		Total: total,
	}

	for i, item := range data {
		price, _ := item.Product.Price.Float64()

		result.Items[i] = GetProductsOutItem{
			ID:             item.Product.ID,
			Name:           item.Product.Name,
			Price:          price,
			StockAvailable: item.Product.StockAvailable,
		}

		if item.ProductPreviewFile != nil {
			result.Items[i].ImagePreview = item.ProductPreviewFile.GetURL(&ctrl.cfg)
		}
	}

	return c.JSON(result)
}
