package logger

import (
	"io"
	"log/slog"
)

func LoadLogger(writer io.Writer) *slog.Logger {
	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})
	return slog.New(handler)
}
