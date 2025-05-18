package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/internal/e"
)

type ErrorJSON struct {
	Code    int      `json:"code"`
	Error   string   `json:"error"`
	Details []string `json:"details"`
}

func ErrorHandler() func(*fiber.Ctx, error) error {
	return func(c *fiber.Ctx, err error) error {
		json := ErrorJSON{
			Code:    500,
			Details: []string{},
		}

		switch errTyped := err.(type) {
		case *e.LogicError:
			json.Code = errTyped.Code()
			json.Error = errTyped.Error()
			json.Details = errTyped.Details()
		case *fiber.Error:
			json.Code = errTyped.Code
			json.Error = err.Error()
		default:
			json.Error = "internal error"
		}

		return c.Status(json.Code).JSON(json)
	}
}
