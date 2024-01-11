package middleware

import (
	"context"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func Logger() echo.MiddlewareFunc {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.LevelKey:
				return slog.String("severity", a.Value.String())
			case slog.MessageKey:
				return slog.String("message", a.Value.String())
			}
			return a
		},
	}))
	return echomiddleware.RequestLoggerWithConfig(echomiddleware.RequestLoggerConfig{
		LogMethod:       true,
		LogURI:          true,
		LogStatus:       true,
		LogResponseSize: true,
		LogUserAgent:    true,
		LogRemoteIP:     true,
		LogReferer:      true,
		LogLatency:      true,
		LogError:        true,
		HandleError:     true,
		LogValuesFunc: func(c echo.Context, v echomiddleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "request",
					slog.Group("httpRequest",
						slog.String("requestMethod", v.Method),
						slog.String("requestUrl", v.URI),
						slog.Int("status", v.Status),
						slog.Int64("responseSize", v.ResponseSize),
						slog.String("userAgent", v.UserAgent),
						slog.String("remoteIp", v.RemoteIP),
						slog.String("referer", v.Referer),
						slog.Float64("latency", v.Latency.Seconds()),
					),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "request error",
					slog.Group("httpRequest",
						slog.String("requestMethod", v.Method),
						slog.String("requestUrl", v.URI),
						slog.Int("status", v.Status),
						slog.Int64("responseSize", v.ResponseSize),
						slog.String("userAgent", v.UserAgent),
						slog.String("remoteIp", v.RemoteIP),
						slog.String("referer", v.Referer),
						slog.Float64("latency", v.Latency.Seconds()),
						slog.String("err", v.Error.Error()),
					),
				)
			}

			return nil
		},
	})
}
