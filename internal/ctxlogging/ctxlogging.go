package ctxlogging

import (
	"context"
	"log/slog"
)

const logger_key = "log"

func Add(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, logger_key, log)
}

func Get(ctx context.Context) *slog.Logger {
	logger := ctx.Value(logger_key)
	if logger == nil {
		slog.Debug("logger is not found in ctx. using default logger")
		return slog.Default()
	}

	log, ok := logger.(*slog.Logger)
	if !ok {
		slog.Debug("there is not *slog.Logger instance under ctx logger key. using default logger")
		return slog.Default()
	}

	return log
}
