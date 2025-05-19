package bootstrap

import (
	"os"

	"github.com/m11ano/mipt-webdev-course/backend/temporal-app/internal/infra/config"
	"go.uber.org/fx/fxevent"
)

func ProvideFXLogger(config config.Config) fxevent.Logger {
	if !config.App.UseFxLogger {
		return fxevent.NopLogger
	}
	return &fxevent.ConsoleLogger{
		W: os.Stdout,
	}
}
