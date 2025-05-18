package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/m11ano/mipt-webdev-course/backend/clients/auth"
)

func Auth(authClient auth.Client) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		isAuth := false

		AuthorizationHeader := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")

		if len(AuthorizationHeader) > 0 {
			claims, err := authClient.ParseJWT(AuthorizationHeader)
			if err != nil && !errors.Is(err, auth.ErrInvalidToken) {
				return err
			}
			if claims != nil {
				isAuth = true
				c.Locals("authAccountID", claims.AccountID)
			}
		}

		c.Locals("isAuth", isAuth)

		return c.Next()
	}
}

type ExtractAuthDataOut struct {
	IsAuth    bool
	AccountID uuid.UUID
}

func ExtractAuthData(c *fiber.Ctx) ExtractAuthDataOut {
	out := ExtractAuthDataOut{}

	if isAuth, ok := c.Locals("isAuth").(bool); !ok || !isAuth {
		return out
	}

	out.IsAuth = true

	if accountID, ok := c.Locals("authAccountID").(uuid.UUID); ok {
		out.AccountID = accountID
	}

	return out
}
