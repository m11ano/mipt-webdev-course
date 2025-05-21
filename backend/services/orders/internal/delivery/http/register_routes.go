package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/m11ano/mipt-webdev-course/backend/clients/auth"
	_ "github.com/m11ano/mipt-webdev-course/backend/services/orders/docs"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/delivery/http/controller"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/delivery/http/middleware"
	"github.com/m11ano/mipt-webdev-course/backend/services/orders/internal/infra/config"
)

func RegisterRoutes(app *fiber.App, config config.Config, ctrl *controller.Controller) {
	authClient := auth.NewClient(config.Secrets.JWT)
	authMiddleware := middleware.Auth(authClient)

	rootGroup := app.Group(config.HTTP.Prefix)
	v1Group := rootGroup.Group("v1", authMiddleware)
	serviceGroup := v1Group.Group(config.HTTP.ServiceName)

	if config.HTTP.StartSwagger {
		serviceGroup.Get("/swagger/*", swagger.HandlerDefault)
	}

	serviceGroup.Post("/", ctrl.CreateOrderHandler)
	serviceGroup.Put("/:id<min(1)>", ctrl.UpdateOrderHandler)
	serviceGroup.Put("/:id<min(1)>/status", ctrl.SetOrderStatusHandler)
	serviceGroup.Get("/", ctrl.GetOrdersHandler)
	serviceGroup.Get("/:id<min(1)>", ctrl.GetOrderHandler)
	serviceGroup.Get("/:id<min(1)>/:secret_key<guid>", ctrl.GetOrderWithSecretKeyHandler)
}
