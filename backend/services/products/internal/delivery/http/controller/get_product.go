package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/m11ano/e"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/delivery/http/middleware"
)

type FileOut struct {
	ID  uuid.UUID `json:"id"`
	URL string    `json:"url"`
}

type GetProductOut struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	IsPublished     bool      `json:"is_published"`
	FullDescription string    `json:"full_description"`
	Price           float64   `json:"price"`
	StockAvailable  int32     `json:"stock_available"`
	ImagePreview    FileOut   `json:"image_preview"`
	Slider          []FileOut `json:"slider"`
}

// @Summary Получить продукт по ID
// @Tags products
// @Produce  json
// @Param id path int true "Product ID"
// @Success 200 {object} GetProductOut
// @Failure 404 {object} middleware.ErrorJSON
// @Router /products/{id} [get]
func (ctrl *Controller) GetProductHandler(c *fiber.Ctx) error {

	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	authData := middleware.ExtractAuthData(c)

	data, err := ctrl.productUC.FindOneFullByID(c.Context(), int64(id), nil)
	if err != nil {
		return err
	}

	if !authData.IsAuth && !data.Product.IsPublished {
		return e.NewErrorFrom(e.ErrNotFound)
	}

	price, _ := data.Product.Price.Float64()

	out := &GetProductOut{
		ID:              data.Product.ID,
		Name:            data.Product.Name,
		IsPublished:     data.Product.IsPublished,
		FullDescription: data.Product.FullDescription,
		Price:           price,
		StockAvailable:  data.Product.StockAvailable,
		Slider:          make([]FileOut, len(data.SliderFiles)),
	}

	if data.ProductPreviewFile != nil {
		out.ImagePreview = FileOut{
			ID:  data.ProductPreviewFile.ID,
			URL: data.ProductPreviewFile.GetURL(&ctrl.cfg),
		}
	}

	for i, item := range data.SliderFiles {
		out.Slider[i] = FileOut{
			ID:  item.ID,
			URL: item.GetURL(&ctrl.cfg),
		}
	}

	return c.JSON(out)
}
