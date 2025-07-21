package logger

import (
	"log/slog"
	"os"
)

func New() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if os.Getenv("GIN_MODE") == "debug" {
		handler := slog.NewTextHandler(os.Stdout, opts)
		return slog.New(handler)
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler)
}
