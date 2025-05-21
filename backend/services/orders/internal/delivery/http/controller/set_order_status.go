package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m11ano/e"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/delivery/http/middleware"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/delivery/http/validation"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/domain"
)

type SetOrderStatusIn struct {
	Status string `json:"status" validate:"required"`
}

func (ctrl *Controller) SetOrderStatusHandlerValidate(in *SetOrderStatusIn) (isOk bool, errMsg []string) {
	if err := ctrl.vldtr.Struct(in); err != nil {
		return validation.FormatErrors(err)
	}
	return true, []string{}
}

// @Summary Поменять статус заказу
// @Security BearerAuth
// @Tags orders
// @Accept  json
// @Param request body SetOrderStatusIn true "JSON"
// @Param id path int true "Order ID"
// @Success 200 {string} string "OK"
// @Failure 400 {object} middleware.ErrorJSON
// @Router /orders/{id}/status [put]
func (ctrl *Controller) SetOrderStatusHandler(c *fiber.Ctx) error {

	authData := middleware.ExtractAuthData(c)

	if !authData.IsAuth {
		return e.ErrUnauthorized
	}

	in := &SetOrderStatusIn{}

	if err := c.BodyParser(in); err != nil {
		return e.NewErrorFrom(e.ErrBadRequest).Wrap(err).SetMessage("cannot parse request body")
	}

	ok, errMsg := ctrl.SetOrderStatusHandlerValidate(in)
	if !ok {
		return e.NewErrorFrom(e.ErrBadRequest).AddDetails(errMsg)
	}

	orderID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	status, ok := domain.OrderStatusMap[in.Status]
	if !ok {
		return e.ErrBadRequest
	}

	err = ctrl.orderUC.SetStatus(c.Context(), int64(orderID), status)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
