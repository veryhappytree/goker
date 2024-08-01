package logger

import (
	"log/slog"
	"os"
)

func Setup() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	slog.SetDefault(logger)
}
