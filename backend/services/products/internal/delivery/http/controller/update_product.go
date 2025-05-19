package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/m11ano/e"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/delivery/http/middleware"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/delivery/http/validation"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/usecase"
	"github.com/shopspring/decimal"
)

type UpdateProductIn struct {
	Name               string      `json:"name" validate:"required,min=1,max=150"`
	IsPublished        bool        `json:"is_published"`
	FullDescription    string      `json:"full_description"`
	Price              float64     `json:"price" validate:"gte=0"`
	ImagePreviewFileID *uuid.UUID  `json:"image_preview_file_id" validate:"required,uuid"`
	SliderFilesIDs     []uuid.UUID `json:"slider_files_ids" validate:"min=1,dive,uuid"`
}

func (ctrl *Controller) UpdateProductHandlerValidate(in *UpdateProductIn) (isOk bool, errMsg []string) {
	if err := ctrl.vldtr.Struct(in); err != nil {
		return validation.FormatErrors(err)
	}
	return true, []string{}
}

// @Summary Редактировать продукт
// @Security BearerAuth
// @Tags products
// @Accept  json
// @Param request body UpdateProductIn true "JSON"
// @Param id path int true "Product ID"
// @Success 200 {string} string "OK"
// @Failure 400 {object} middleware.ErrorJSON
// @Router /products/{id} [put]
func (ctrl *Controller) UpdateProductHandler(c *fiber.Ctx) error {

	authData := middleware.ExtractAuthData(c)

	if !authData.IsAuth {
		return e.ErrUnauthorized
	}

	productID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	in := &UpdateProductIn{}

	if err := c.BodyParser(in); err != nil {
		return e.NewErrorFrom(e.ErrBadRequest).Wrap(err).SetMessage("cannot parse request body")
	}

	ok, errMsg := ctrl.UpdateProductHandlerValidate(in)
	if !ok {
		return e.NewErrorFrom(e.ErrBadRequest).AddDetails(errMsg)
	}

	_, _, err = ctrl.productUC.Update(c.Context(), int64(productID), usecase.ProductUpdateIn{
		Name:               in.Name,
		IsPublished:        in.IsPublished,
		FullDescription:    in.FullDescription,
		Price:              decimal.NewFromFloat(in.Price),
		ImagePreviewFileID: in.ImagePreviewFileID,
		SliderFilesIDs:     in.SliderFilesIDs,
	})
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
