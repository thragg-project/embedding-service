package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

func New(level, version string) (*slog.Logger, error) {
	return NewWithWriter(os.Stdout, level, version)
}

func NewWithWriter(writer io.Writer, level, version string) (*slog.Logger, error) {
	slogLevel, err := parseLevel(level)
	if err != nil {
		return nil, err
	}

	var handler slog.Handler = slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level: slogLevel,
	})
	if version != "" {
		handler = handler.WithAttrs([]slog.Attr{slog.String("version", version)})
	}
	return slog.New(handler), nil
}

func Fatal(ctx context.Context, log *slog.Logger, message string, err error) {
	log.ErrorContext(ctx, message, slog.Any("error", err))
	os.Exit(1)
}

func parseLevel(level string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "", "info":
		return slog.LevelInfo, nil
	case "debug":
		return slog.LevelDebug, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("unknown logger level %q", level)
	}
}
