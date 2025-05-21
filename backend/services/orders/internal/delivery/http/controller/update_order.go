package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m11ano/e"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/delivery/http/middleware"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/delivery/http/validation"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/usecase"
	"github.com/shopspring/decimal"
)

type UpdateOrderIn struct {
	Details  UpdateOrderInDetails   `json:"details" validate:"required"`
	Products []UpdateOrderInProduct `json:"products" validate:"required,min=1"`
}

type UpdateOrderInDetails struct {
	ClientName      string `json:"client_name" validate:"required,min=1,max=150"`
	ClientSurname   string `json:"client_surname" validate:"required,min=1,max=150"`
	ClientEmail     string `json:"client_email" validate:"required,email"`
	ClientPhone     string `json:"client_phone" validate:"required,min=1,max=20"`
	DeliveryAddress string `json:"delivery_address" validate:"required,min=1,max=150"`
}

type UpdateOrderInProduct struct {
	ID       int64   `json:"id" validate:"gte=0"`
	Quantity int32   `json:"quantity" validate:"gte=1"`
	Price    float64 `json:"price" validate:"gte=0"`
}

func (ctrl *Controller) UpdateOrderHandlerValidate(in *UpdateOrderIn) (isOk bool, errMsg []string) {
	if err := ctrl.vldtr.Struct(in); err != nil {
		return validation.FormatErrors(err)
	}
	return true, []string{}
}

// @Summary Обновить заказ
// @Security BearerAuth
// @Tags orders
// @Accept  json
// @Param request body UpdateOrderIn true "JSON"
// @Param id path int true "Order ID"
// @Success 200 {string} string "OK"
// @Failure 400 {object} middleware.ErrorJSON
// @Router /orders/{id} [put]
func (ctrl *Controller) UpdateOrderHandler(c *fiber.Ctx) error {

	authData := middleware.ExtractAuthData(c)

	if !authData.IsAuth {
		return e.ErrUnauthorized
	}

	in := &UpdateOrderIn{}

	if err := c.BodyParser(in); err != nil {
		return e.NewErrorFrom(e.ErrBadRequest).Wrap(err).SetMessage("cannot parse request body")
	}

	ok, errMsg := ctrl.UpdateOrderHandlerValidate(in)
	if !ok {
		return e.NewErrorFrom(e.ErrBadRequest).AddDetails(errMsg)
	}

	orderID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	updateIn := usecase.OrderUpdateIn{
		Details: usecase.OrderDataDetailsIn{
			ClientName:      in.Details.ClientName,
			ClientSurname:   in.Details.ClientSurname,
			ClientEmail:     in.Details.ClientEmail,
			ClientPhone:     in.Details.ClientPhone,
			DeliveryAddress: in.Details.DeliveryAddress,
		},
		Products: make([]usecase.OrderProductWithPrice, len(in.Products)),
	}

	for i, item := range in.Products {
		updateIn.Products[i] = usecase.OrderProductWithPrice{
			ID:       item.ID,
			Quantity: item.Quantity,
			Price:    decimal.NewFromFloat(item.Price),
		}
	}

	err = ctrl.orderUC.Update(c.Context(), int64(orderID), updateIn)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
