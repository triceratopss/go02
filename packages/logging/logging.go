package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"
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
		return slog.Attr{
			Key:   "logging.googleapis.com/sourceLocation",
			Value: a.Value,
		}
	}

	return a
}

func log(ctx context.Context, level slog.Level, msg string, args ...any) {
	logger := slog.Default()
	if !logger.Enabled(ctx, level) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])

	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(args...)

	_ = logger.Handler().Handle(ctx, r)
}

func Info(msg string, args ...any) {
	log(context.Background(), slog.LevelInfo, msg, args...)
}

func Infof(format string, args ...any) {
	log(context.Background(), slog.LevelInfo, fmt.Sprintf(format, args...))
}

func Error(msg string, args ...interface{}) {
	log(context.Background(), slog.LevelError, msg, args...)
}

func Errorf(err error, format string, args ...interface{}) {
	log(context.Background(), slog.LevelError, fmt.Sprintf(format, args...))
}
