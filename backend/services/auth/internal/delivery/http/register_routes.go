package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "github.com/m11ano/mipt-webdev-course/backend/services/auth/docs"
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/internal/delivery/http/controller"
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/internal/delivery/http/middleware"
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/internal/infra/config"
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/pkg/auth"
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

	serviceGroup.Post("/", ctrl.AuthCheckHandler)
	serviceGroup.Post("/login", ctrl.AuthLoginHandler)
}
