package logger

import (
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	Env     string // dev|prod
	Service string // api|worker|...
	Level   string // debug|info|warn|error
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "info", "":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func New(cfg Config) *slog.Logger {
	lvl := parseLevel(cfg.Level)

	var h slog.Handler
	if strings.ToLower(cfg.Env) == "dev" {
		// человекочитаемый вывод в dev
		h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     lvl,
			AddSource: true,
		})
	} else {
		// JSON для продакшена / агрегаторов логов
		h = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     lvl,
			AddSource: true,
		})
	}

	return slog.New(h).With(
		"service", cfg.Service,
		"env", strings.ToLower(cfg.Env),
	)
}
