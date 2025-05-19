package bootstrap

import (
	"io"
	"log/slog"
	"os"

	"github.com/m11ano/mipt-webdev-course/backend/temporal-app/internal/infra/config"

	"github.com/imperatorofdwelling/Website-backend/pkg/logger/slogpretty"
)

func NewLogger(config config.Config) *slog.Logger {
	var logger *slog.Logger

	switch {
	case !config.App.UseLogger:
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	case config.App.IsProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	default:
		logger = setupPrettySlog()
	}

	return logger
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
