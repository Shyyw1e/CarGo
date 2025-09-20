package middlewares

import (
	"context"
	"log/slog"
)

type ctxKey int

const loggerKey ctxKey = iota + 1

func WithLogger(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, log)
}

func FromContext(ctx context.Context) *slog.Logger {
	if v := ctx.Value(loggerKey); v != nil {
		if lg, ok := v.(*slog.Logger); ok {
			return lg
		}
	}
	return slog.Default()
}
