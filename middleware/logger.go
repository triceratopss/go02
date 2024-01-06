package middleware

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func Logger() echo.MiddlewareFunc {
	return echomiddleware.RequestLoggerWithConfig(echomiddleware.RequestLoggerConfig{
		LogMethod:       true,
		LogURI:          true,
		LogStatus:       true,
		LogResponseSize: true,
		LogUserAgent:    true,
		LogRemoteIP:     true,
		LogReferer:      true,
		LogLatency:      true,
		LogValuesFunc: func(c echo.Context, v echomiddleware.RequestLoggerValues) error {
			slog.Info("http request received",
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

			return nil
		},
	})
}
