package bootstrap

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/m11ano/mipt-webdev-course/backend/services/auth/internal/delivery/http/middleware"
)

const defaultBodyLimit = 10 * 1024 * 1024

type HTTPConfig struct {
	UnderProxy bool
	UseTraceID bool
	UseLogger  bool
	BodyLimit  int
}

func NewHTTPFiber(httpCfg HTTPConfig, logger *slog.Logger) *fiber.App {
	if httpCfg.BodyLimit == -1 {
		httpCfg.BodyLimit = defaultBodyLimit
	}

	fiberCfg := fiber.Config{
		ErrorHandler: middleware.ErrorHandler(),
		BodyLimit:    httpCfg.BodyLimit,
	}

	if httpCfg.UnderProxy {
		fiberCfg.ProxyHeader = fiber.HeaderXForwardedFor
	}

	app := fiber.New(fiberCfg)

	app.Use(middleware.Recovery(logger))

	if httpCfg.UseTraceID {
		app.Use(middleware.TraceID())
	}

	if httpCfg.UseLogger {
		app.Use(middleware.Logger(logger))
	}

	return app
}
