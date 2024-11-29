package logging

import (
	"context"
	"fmt"
	"go02/internal/package/apperrors"
	"log/slog"
	"os"
	"runtime"
	"time"

	"go.opentelemetry.io/otel/trace"
)

func Init() {
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		ReplaceAttr: replacer,
	})
	instrumentedHandler := handlerWithSpanContext(jsonHandler)
	slog.SetDefault(slog.New(instrumentedHandler))
}

func handlerWithSpanContext(handler slog.Handler) *spanContextLogHandler {
	return &spanContextLogHandler{Handler: handler}
}

type spanContextLogHandler struct {
	slog.Handler
}

func (t *spanContextLogHandler) Handle(ctx context.Context, record slog.Record) error {
	if s := trace.SpanContextFromContext(ctx); s.IsValid() {
		record.AddAttrs(
			slog.Any("logging.googleapis.com/trace", s.TraceID()),
		)
		record.AddAttrs(
			slog.Any("logging.googleapis.com/spanId", s.SpanID()),
		)
		record.AddAttrs(
			slog.Bool("logging.googleapis.com/trace_sampled", s.TraceFlags().IsSampled()),
		)
	}
	return t.Handler.Handle(ctx, record)
}

func replacer(groups []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.LevelKey:
		a.Key = "severity"
		if level := a.Value.Any().(slog.Level); level == slog.LevelWarn {
			a.Value = slog.StringValue("WARNING")
		}
	case slog.TimeKey:
		a.Key = "timestamp"
	case slog.MessageKey:
		a.Key = "message"
	case slog.SourceKey:
		a.Key = "logging.googleapis.com/sourceLocation"
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

func Info(ctx context.Context, msg string, args ...any) {
	log(ctx, slog.LevelInfo, msg, args...)
}

func Infof(ctx context.Context, format string, args ...any) {
	log(ctx, slog.LevelInfo, fmt.Sprintf(format, args...))
}

func Error(ctx context.Context, err error, msg string, args ...any) {
	args = append(args, apperrors.LogStackTrace(err))
	log(ctx, slog.LevelError, msg, args...)
}

func Errorf(ctx context.Context, err error, format string, args ...any) {
	log(ctx, slog.LevelError, fmt.Sprintf(format, args...), apperrors.LogStackTrace(err))
}
