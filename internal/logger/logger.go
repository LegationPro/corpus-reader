package logger

import (
	"log/slog"
	"os"
)

// Creates a new logger with JSON encoding
func New() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
