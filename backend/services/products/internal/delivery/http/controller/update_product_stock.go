package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/delivery/http/middleware"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/delivery/http/validation"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/e"
)

type UpdateProductStockIn struct {
	Operation string `json:"operation" validate:"required,oneof=increase decrease"`
	Value     int32  `json:"value" validate:"gte=1"`
}

func (ctrl *Controller) UpdateProductStockHandlerValidate(in *UpdateProductStockIn) (isOk bool, errMsg []string) {
	if err := ctrl.vldtr.Struct(in); err != nil {
		return validation.FormatErrors(err)
	}
	return true, []string{}
}

// @Summary Изменить остаток товара на складе
// @Security BearerAuth
// @Tags products
// @Accept  json
// @Param request body UpdateProductStockIn true "JSON"
// @Param id path int true "Product ID"
// @Success 200 {string} string "OK"
// @Failure 400 {object} middleware.ErrorJSON
// @Router /products/{id}/stock [post]
func (ctrl *Controller) UpdateProductStockHandler(c *fiber.Ctx) error {

	authData := middleware.ExtractAuthData(c)

	if !authData.IsAuth {
		return e.ErrUnauthorized
	}

	productID, err := c.ParamsInt("id", 0)
	if err != nil {
		return err
	}

	in := &UpdateProductStockIn{}

	if err := c.BodyParser(in); err != nil {
		return e.NewErrorFrom(e.ErrBadRequest).Wrap(err).SetMessage("cannot parse request body")
	}

	ok, errMsg := ctrl.UpdateProductStockHandlerValidate(in)
	if !ok {
		return e.NewErrorFrom(e.ErrBadRequest).AddDetails(errMsg)
	}

	isIncrease := true
	if in.Operation == "decrease" {
		isIncrease = false
	}

	err = ctrl.productUC.ChangeStock(c.Context(), int64(productID), in.Value, isIncrease)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
