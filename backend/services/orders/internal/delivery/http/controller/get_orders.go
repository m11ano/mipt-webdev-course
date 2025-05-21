package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/m11ano/e"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/delivery/http/middleware"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/usecase"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/usecase/uctypes"
	"github.com/samber/lo"
)

type GetOrdersOutItem struct {
	ID        int64              `json:"id"`
	SecretKey uuid.UUID          `json:"secret_key"`
	OrderSum  float64            `json:"order_sum"`
	Status    string             `json:"status"`
	Details   GetOrderOutDetails `json:"details"`
}

type GetOrdersOut struct {
	Items []GetOrdersOutItem `json:"items"`
	Total int64              `json:"total"`
}

// @Summary Получить список заказов
// @Security BearerAuth
// @Tags orders
// @Produce  json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} GetOrdersOut
// @Failure 400 {object} middleware.ErrorJSON
// @Router /orders [get]
func (ctrl *Controller) GetOrdersHandler(c *fiber.Ctx) error {

	authData := middleware.ExtractAuthData(c)

	if !authData.IsAuth {
		return e.ErrUnauthorized
	}

	limit := c.QueryInt("limit", 20)
	if limit > 100 {
		limit = 100
	}
	if limit < 1 {
		limit = 1
	}

	offset := c.QueryInt("offset", 0)
	if offset < 0 {
		offset = 0
	}

	data, total, err := ctrl.orderUC.FindPagedList(c.Context(), usecase.OrderListOptions{
		OnlyCreated: lo.ToPtr(true),
		Sort: &[]usecase.OrderListSort{
			{
				Field:  usecase.OrderListSortFieldID,
				IsDesc: true,
			},
		},
	}, &uctypes.QueryGetListParams{
		Limit:  uint64(limit),
		Offset: uint64(offset),
	})
	if err != nil {
		return err
	}

	result := GetOrdersOut{
		Items: make([]GetOrdersOutItem, len(data)),
		Total: total,
	}

	for i, item := range data {
		orderSum, _ := item.OrderSum.Float64()

		result.Items[i] = GetOrdersOutItem{
			ID:        item.ID,
			SecretKey: item.SecretKey,
			OrderSum:  orderSum,
			Status:    item.Status.String(),
			Details: GetOrderOutDetails{
				ClientName:      item.ClientName,
				ClientSurname:   item.ClientSurname,
				ClientEmail:     item.ClientEmail,
				ClientPhone:     item.ClientPhone,
				DeliveryAddress: item.DeliveryAddress,
			},
		}
	}

	return c.JSON(result)
}
