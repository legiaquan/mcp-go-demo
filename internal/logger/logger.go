package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func Init() {
	Log = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}
