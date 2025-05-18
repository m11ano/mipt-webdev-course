package controller

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/internal/delivery/http/middleware"
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/internal/e"
)

type AuthCheckHandlerOut struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Surname string    `json:"surname"`
	Email   string    `json:"email"`
}

// @Summary Проверить сессию и получить информацию о пользователе
// @Security BearerAuth
// @Tags auth
// @Accept  json
// @Success 200 {object} AuthCheckHandlerOut
// @Failure 400 {object} middleware.ErrorJSON
// @Router /auth [post]
func (ctrl *Controller) AuthCheckHandler(c *fiber.Ctx) error {

	authData := middleware.ExtractAuthData(c)

	if !authData.IsAuth {
		return e.ErrUnauthorized
	}

	account, err := ctrl.accountUC.FindOneByID(c.Context(), authData.AccountID, nil)
	if err != nil {
		if errors.Is(err, e.ErrNotFound) {
			return e.ErrUnauthorized
		}
		return err
	}

	out := AuthCheckHandlerOut{
		ID:      account.ID,
		Name:    account.Name,
		Surname: account.Surname,
		Email:   account.Email,
	}

	return c.JSON(out)
}
