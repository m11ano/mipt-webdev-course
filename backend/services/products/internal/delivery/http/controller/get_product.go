package controller

import (
	"github.com/gofiber/fiber/v2"
	_ "github.com/m11ano/mipt-webdev-course/backend/services/products/docs"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/delivery/http/middleware"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/e"
)

type GetProductOut struct {
	ID              int64    `json:"id"`
	Name            string   `json:"name"`
	IsPublished     bool     `json:"is_published"`
	FullDescription string   `json:"full_description"`
	Price           float64  `json:"price"`
	StockAvailable  int32    `json:"stock_available"`
	ImagePreview    string   `json:"image_preview"`
	Slider          []string `json:"slider"`
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
		Slider:          make([]string, len(data.SliderFiles)),
	}

	if data.ProductPreviewFile != nil {
		out.ImagePreview = data.ProductPreviewFile.GetURL(&ctrl.cfg)
	}

	for i, item := range data.SliderFiles {
		out.Slider[i] = item.GetURL(&ctrl.cfg)
	}

	return c.JSON(out)
}
