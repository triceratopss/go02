package logger

import (
	"log/slog"
	"os"
)

func Init() {
	opt := slog.HandlerOptions{
		AddSource:   true,
		ReplaceAttr: replaceAttr,
	}

	if os.Getenv("ENV") == "local" {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &opt)))
		return
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &opt)))
}

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.LevelKey:
		return slog.String("severity", a.Value.String())
	case slog.MessageKey:
		return slog.String("message", a.Value.String())
	case slog.SourceKey:
		return slog.String("sourceLocation", a.Value.String())
	}

	return a
}
