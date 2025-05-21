package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/m11ano/e"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/delivery/http/middleware"
)

type GetOrderOut struct {
	ID        int64                `json:"id"`
	SecretKey uuid.UUID            `json:"secret_key"`
	OrderSum  float64              `json:"order_sum"`
	Status    string               `json:"status"`
	Details   GetOrderOutDetails   `json:"details"`
	Products  []GetOrderOutProduct `json:"products"`
}

type GetOrderOutDetails struct {
	ClientName      string `json:"client_name"`
	ClientSurname   string `json:"client_surname"`
	ClientEmail     string `json:"client_email"`
	ClientPhone     string `json:"client_phone"`
	DeliveryAddress string `json:"delivery_address"`
}

type GetOrderOutProduct struct {
	ID       int64   `json:"id"`
	Quantity int32   `json:"quantity"`
	Price    float64 `json:"price"`
}

// @Summary Получить заказ по ID
// @Security BearerAuth
// @Tags orders
// @Produce  json
// @Param id path int true "Order ID"
// @Success 200 {object} GetOrderOut
// @Failure 404 {object} middleware.ErrorJSON
// @Router /orders/{id} [get]
func (ctrl *Controller) GetOrderHandler(c *fiber.Ctx) error {

	authData := middleware.ExtractAuthData(c)

	if !authData.IsAuth {
		return e.ErrUnauthorized
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	data, err := ctrl.orderUC.FindOneFullByID(c.Context(), int64(id), nil)
	if err != nil {
		return err
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
