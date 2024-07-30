package slogutils

import (
	"fmt"
	"log/slog"
	"message-processor/internal/config"
	"os"
)

func MustNewLogger(env config.Env) (logger *slog.Logger) {
	switch env {
	case config.EnvLocal:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case config.EnvDev, config.EnvProd:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		panic(fmt.Errorf("unknown env: %v", env))
	}

	return logger
}
