package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
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

func Error(err error, msg string, args ...any) {
	stackTrace := getFormattedStackTrace(err)
	args = append(args, slog.Any("stack_trace", stackTrace))
	log(context.Background(), slog.LevelError, msg, args...)
}

func Errorf(err error, format string, args ...any) {
	formattedMsg := fmt.Sprintf(format, args...)
	stackTrace := getFormattedStackTrace(err)
	log(context.Background(), slog.LevelError, formattedMsg, slog.Any("stack_trace", stackTrace))
}

// スタックトレースを整形する関数
func getFormattedStackTrace(err error) []map[string]string {
	stackTrace := fmt.Sprintf("%+v", err)
	lines := strings.Split(stackTrace, "\n")
	var formattedStackTrace []map[string]string

	for i := 0; i < len(lines); i++ {
		trimmedLine := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmedLine, "|") {
			// 関数名
			function := strings.TrimSpace(strings.TrimPrefix(trimmedLine, "|"))
			// 次の行にファイル名と行番号がある
			if i+1 < len(lines) {
				i++
				fileLine := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(lines[i]), "|"))
				formattedStackTrace = append(formattedStackTrace, map[string]string{
					"function": function,
					"location": fileLine,
				})
			}
		}
	}

	return formattedStackTrace
}
