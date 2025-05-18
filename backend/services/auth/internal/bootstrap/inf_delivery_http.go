package bootstrap

import (
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/internal/delivery/http"
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/internal/delivery/http/controller"
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/internal/delivery/http/validation"
	"go.uber.org/fx"
)

var DeliveryHTTP = fx.Options(
	fx.Provide(validation.NewValidator),
	fx.Provide(controller.New),
	fx.Invoke(http.RegisterRoutes),
)
