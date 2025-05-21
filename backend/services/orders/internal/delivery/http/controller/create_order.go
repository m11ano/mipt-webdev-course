package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/m11ano/e"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/delivery/http/validation"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/usecase"
)

type CreateOrderIn struct {
	Details  CreateOrderInDetails   `json:"details" validate:"required"`
	Products []CreateOrderInProduct `json:"products" validate:"required,min=1"`
}

type CreateOrderInDetails struct {
	ClientName      string `json:"client_name" validate:"required,min=1,max=150"`
	ClientSurname   string `json:"client_surname" validate:"required,min=1,max=150"`
	ClientEmail     string `json:"client_email" validate:"required,email"`
	ClientPhone     string `json:"client_phone" validate:"required,min=1,max=20"`
	DeliveryAddress string `json:"delivery_address" validate:"required,min=1,max=150"`
}

type CreateOrderInProduct struct {
	ID       int64 `json:"id" validate:"gte=0"`
	Quantity int32 `json:"quantity" validate:"gte=1"`
}

type CreateOrderOut struct {
	ID        int64     `json:"id"`
	SecretKey uuid.UUID `json:"secret_key"`
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

	createIn := usecase.OrderCreateIn{
		Details: usecase.OrderCreateInDetails{
			ClientName:      in.Details.ClientName,
			ClientSurname:   in.Details.ClientSurname,
			ClientEmail:     in.Details.ClientEmail,
			ClientPhone:     in.Details.ClientPhone,
			DeliveryAddress: in.Details.DeliveryAddress,
		},
		Products: make([]usecase.OrderProductIn, len(in.Products)),
	}

	for i, item := range in.Products {
		createIn.Products[i] = usecase.OrderProductIn{
			ID:       item.ID,
			Quantity: item.Quantity,
		}
	}

	order, err := ctrl.orderUC.Create(c.Context(), createIn)
	if err != nil {
		return err
	}

	result := CreateOrderOut{
		ID:        order.ID,
		SecretKey: order.SecretKey,
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}
