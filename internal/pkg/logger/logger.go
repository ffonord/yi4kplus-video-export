package logger

import (
	"io"
	"log/slog"
	"os"
)

const (
	EnvTest  = "test"
	EnvLocal = "local"
	EnvProd  = "prod"
	EnvDev   = "dev"
)

type Logger struct {
	*slog.Logger
}

func New(env string) *Logger {
	var log *slog.Logger

	switch env {
	case EnvTest:
		log = slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case EnvLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case EnvDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case EnvProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return &Logger{log}
}
