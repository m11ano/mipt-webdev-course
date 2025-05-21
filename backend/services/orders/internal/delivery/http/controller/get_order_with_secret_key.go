package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/m11ano/e"
)

// @Summary Получить заказ по ID + secret_key
// @Tags orders
// @Produce  json
// @Param id path int true "Order ID"
// @Success 200 {object} GetOrderOut
// @Failure 404 {object} middleware.ErrorJSON
// @Router /orders/{id}/{secret_key} [get]
func (ctrl *Controller) GetOrderWithSecretKeyHandler(c *fiber.Ctx) error {

	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	secretKeyStr := c.Params("secret_key")
	secretKey, err := uuid.Parse(secretKeyStr)
	if err != nil {
		return err
	}

	data, err := ctrl.orderUC.FindOneFullByID(c.Context(), int64(id), nil)
	if err != nil {
		return err
	}

	if data.Order.SecretKey != secretKey {
		return e.ErrForbidden
	}

	orderSum, _ := data.Order.OrderSum.Float64()

	out := &GetOrderOut{
		ID:        data.Order.ID,
		OrderSum:  orderSum,
		Status:    data.Order.Status.String(),
		SecretKey: data.Order.SecretKey,
		Details: GetOrderOutDetails{
			ClientName:      data.Order.ClientName,
			ClientSurname:   data.Order.ClientSurname,
			ClientEmail:     data.Order.ClientEmail,
			ClientPhone:     data.Order.ClientPhone,
			DeliveryAddress: data.Order.DeliveryAddress,
		},
		Products: make([]GetOrderOutProduct, len(data.Products)),
	}

	for i, product := range data.Products {
		price, _ := product.Price.Float64()

		out.Products[i] = GetOrderOutProduct{
			ID:       product.ID,
			Quantity: product.Quantity,
			Price:    price,
		}
	}

	return c.JSON(out)
}
