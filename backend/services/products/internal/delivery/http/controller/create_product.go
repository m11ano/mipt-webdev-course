package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/m11ano/e"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/delivery/http/middleware"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/delivery/http/validation"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/domain"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/usecase"
	"github.com/shopspring/decimal"
)

type CreateProductIn struct {
	Name               string      `json:"name" validate:"required,min=1,max=150"`
	IsPublished        bool        `json:"is_published"`
	FullDescription    string      `json:"full_description"`
	Price              float64     `json:"price" validate:"gte=0"`
	StockAvailable     int32       `json:"stock_available" validate:"gte=0"`
	ImagePreviewFileID *uuid.UUID  `json:"image_preview_file_id" validate:"required,uuid"`
	SliderFilesIDs     []uuid.UUID `json:"slider_files_ids" validate:"min=1,dive,uuid"`
}

type CreateProductOut struct {
	ID int64 `json:"id"`
}

func (ctrl *Controller) CreateProductHandlerValidate(in *CreateProductIn) (isOk bool, errMsg []string) {
	if err := ctrl.vldtr.Struct(in); err != nil {
		return validation.FormatErrors(err)
	}
	return true, []string{}
}

// @Summary Создать продукт
// @Security BearerAuth
// @Tags products
// @Accept  json
// @Produce  json
// @Param request body CreateProductIn true "JSON"
// @Success 201 {object} CreateProductOut
// @Failure 400 {object} middleware.ErrorJSON
// @Router /products [post]
func (ctrl *Controller) CreateProductHandler(c *fiber.Ctx) error {

	authData := middleware.ExtractAuthData(c)

	if !authData.IsAuth {
		return e.ErrUnauthorized
	}

	in := &CreateProductIn{}

	if err := c.BodyParser(in); err != nil {
		return e.NewErrorFrom(e.ErrBadRequest).Wrap(err).SetMessage("cannot parse request body")
	}

	ok, errMsg := ctrl.CreateProductHandlerValidate(in)
	if !ok {
		return e.NewErrorFrom(e.ErrBadRequest).AddDetails(errMsg)
	}

	product := domain.NewProduct(0)
	err := product.SetPrice(decimal.NewFromFloat(in.Price))
	if err != nil {
		return err
	}
	product.Name = in.Name
	product.IsPublished = in.IsPublished
	product.FullDescription = in.FullDescription
	product.StockAvailable = in.StockAvailable
	product.ImagePreviewFileID = in.ImagePreviewFileID

	createIn := usecase.ProductCreateIn{
		Product:        product,
		SliderFilesIDs: in.SliderFilesIDs,
	}

	data, _, err := ctrl.productUC.Create(c.Context(), createIn)
	if err != nil {
		return err
	}

	result := CreateProductOut{
		ID: data.ID,
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}
