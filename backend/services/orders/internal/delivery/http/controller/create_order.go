package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/delivery/http/validation"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/e"
)

type CreateOrderIn struct {
	Name               string      `json:"name" validate:"required,min=1,max=150"`
	IsPublished        bool        `json:"is_published"`
	FullDescription    string      `json:"full_description"`
	Price              float64     `json:"price" validate:"gte=0"`
	StockAvailable     int32       `json:"stock_available" validate:"gte=0"`
	ImagePreviewFileID *uuid.UUID  `json:"image_preview_file_id" validate:"required,uuid"`
	SliderFilesIDs     []uuid.UUID `json:"slider_files_ids" validate:"min=1,dive,uuid"`
}

type CreateOrderOut struct {
	ID int64 `json:"id"`
}

func (ctrl *Controller) CreateOrderHandlerValidate(in *CreateOrderIn) (isOk bool, errMsg []string) {
	if err := ctrl.vldtr.Struct(in); err != nil {
		return validation.FormatErrors(err)
	}
	return true, []string{}
}

// @Summary Создать заказ
// @Security BearerAuth
// @Tags orders
// @Accept  json
// @Produce  json
// @Param request body CreateOrderIn true "JSON"
// @Success 201 {object} CreateOrderOut
// @Failure 400 {object} middleware.ErrorJSON
// @Router /orders [post]
func (ctrl *Controller) CreateOrderHandler(c *fiber.Ctx) error {

	in := &CreateOrderIn{}

	if err := c.BodyParser(in); err != nil {
		return e.NewErrorFrom(e.ErrBadRequest).Wrap(err).SetMessage("cannot parse request body")
	}

	ok, errMsg := ctrl.CreateOrderHandlerValidate(in)
	if !ok {
		return e.NewErrorFrom(e.ErrBadRequest).AddDetails(errMsg)
	}

	result := CreateOrderOut{
		ID: 1,
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}
