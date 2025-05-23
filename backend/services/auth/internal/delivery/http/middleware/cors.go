package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Cors() func(*fiber.Ctx) error {
	return cors.New(cors.Config{
		AllowOrigins: "*",
	})
}
