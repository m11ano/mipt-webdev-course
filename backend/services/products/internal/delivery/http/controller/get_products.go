package controller

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/m11ano/e"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/delivery/http/middleware"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/usecase"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/usecase/uctypes"
	"github.com/samber/lo"
)

type GetProductsOutItem struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	IsPublished    bool    `json:"is_published"`
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
// @Param ids query string false "IDs of products, separated by comma. If not empty, then limit and offset will be ignored"
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

	IDsStr := c.Query("ids")
	IDs := make([]int64, 0)
	if IDsStr != "" {
		IDsSeparated := strings.Split(IDsStr, ",")
		for _, idStr := range IDsSeparated {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				return e.NewErrorFrom(e.ErrBadRequest).Wrap(err).SetMessage("invalid ids")
			}
			IDs = append(IDs, id)
		}
	}

	authData := middleware.ExtractAuthData(c)

	listSort := usecase.ProductListOptions{
		Sort: &[]usecase.ProductListSort{
			{
				Field:  usecase.ProductListSortFieldCreatedAt,
				IsDesc: true,
			},
		},
	}

	if !authData.IsAuth && len(IDs) == 0 {
		listSort.IsPublished = lo.ToPtr(true)
	}

	if len(IDs) > 0 {
		listSort.IDs = &IDs
		limit = 100
	}

	queryParams := &uctypes.QueryGetListParams{
		Limit: uint64(limit),
	}

	if len(IDs) == 0 {
		queryParams.Offset = uint64(offset)
	}

	data, total, err := ctrl.productUC.FindFullPagedList(c.Context(), listSort, queryParams)
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
			IsPublished:    item.Product.IsPublished,
			Price:          price,
			StockAvailable: item.Product.StockAvailable,
		}

		if item.ProductPreviewFile != nil {
			result.Items[i].ImagePreview = item.ProductPreviewFile.GetURL(&ctrl.cfg)
		}
	}

	return c.JSON(result)
}
