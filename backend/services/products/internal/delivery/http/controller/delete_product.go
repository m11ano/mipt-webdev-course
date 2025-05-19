package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m11ano/e"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/delivery/http/middleware"
)

// @Summary Удалить продукт
// @Security BearerAuth
// @Tags products
// @Param id path int true "Product ID"
// @Success 200 {string} string "OK"
// @Failure 400 {object} middleware.ErrorJSON
// @Router /products/{id} [delete]
func (ctrl *Controller) DeleteProductHandler(c *fiber.Ctx) error {

	authData := middleware.ExtractAuthData(c)

	if !authData.IsAuth {
		return e.ErrUnauthorized
	}

	productID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	err = ctrl.productUC.Delete(c.Context(), int64(productID))
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
