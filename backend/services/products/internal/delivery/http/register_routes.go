package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/pkg/auth"
	_ "github.com/m11ano/mipt-webdev-course/backend/services/products/docs"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/delivery/http/controller"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/delivery/http/middleware"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/infra/config"
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

	serviceGroup.Get("/", ctrl.GetProductsHandler)
	serviceGroup.Get("/:id<min(1)>", ctrl.GetProductHandler)
	serviceGroup.Post("/", ctrl.CreateProductHandler)
	serviceGroup.Put("/:id<min(1)>", ctrl.UpdateProductHandler)
	serviceGroup.Delete("/:id<min(1)>", ctrl.DeleteProductHandler)
	serviceGroup.Post("/:id<min(1)>/stock", ctrl.UpdateProductStockHandler)

	serviceGroup.Post("/image", ctrl.UploadImageHandler)
}
