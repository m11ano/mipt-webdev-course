package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m11ano/e"
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/internal/delivery/http/validation"
)

type AuthLoginHandlerIn struct {
	Email    string `json:"email" validate:"required,email,max=150"`
	Password string `json:"password" validate:"required,max=150"`
}

func (ctrl *Controller) AuthLoginHandlerValidate(in *AuthLoginHandlerIn) (isOk bool, errMsg []string) {
	if err := ctrl.vldtr.Struct(in); err != nil {
		return validation.FormatErrors(err)
	}
	return true, []string{}
}

// @Summary Аутентификация пользователя
// @Tags auth
// @Accept  json
// @Param request body AuthLoginHandlerIn true "JSON"
// @Success 200 {string} string "JWT token"
// @Failure 400 {object} middleware.ErrorJSON
// @Router /auth/login [post]
func (ctrl *Controller) AuthLoginHandler(c *fiber.Ctx) error {
	in := &AuthLoginHandlerIn{}

	if err := c.BodyParser(in); err != nil {
		return e.NewErrorFrom(e.ErrBadRequest).Wrap(err).SetMessage("cannot parse request body")
	}

	ok, errMsg := ctrl.AuthLoginHandlerValidate(in)
	if !ok {
		return e.NewErrorFrom(e.ErrBadRequest).AddDetails(errMsg)
	}

	jwtToken, err := ctrl.authUC.Login(c.Context(), in.Email, in.Password)
	if err != nil {
		if isAppErr, appErr := e.IsAppError(err); isAppErr {
			return appErr
		}
		return e.ErrInternal
	}

	return c.JSON(jwtToken)
}
