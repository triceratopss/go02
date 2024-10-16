package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

type HTTPRequest struct {
	RequestMethod string   `json:"requestMethod"`
	RequestURL    string   `json:"requestUrl"`
	RequestSize   string   `json:"requestSize,omitempty"`
	Status        int      `json:"status"`
	ResponseSize  string   `json:"responseSize,omitempty"`
	UserAgent     string   `json:"userAgent,omitempty"`
	RemoteIP      string   `json:"remoteIp,omitempty"`
	ServerIP      string   `json:"serverIp,omitempty"`
	Referer       string   `json:"referer,omitempty"`
	Latency       Duration `json:"latency"`
	Protocol      string   `json:"protocol"`
}

type Duration struct {
	Nanos   int32 `json:"nanos"`
	Seconds int64 `json:"seconds"`
}

func MakeDuration(d time.Duration) Duration {
	nanos := d.Nanoseconds()
	secs := nanos / 1e9
	nanos -= secs * 1e9
	return Duration{
		Nanos:   int32(nanos),
		Seconds: secs,
	}
}

func Logger() echo.MiddlewareFunc {
	logger := slog.Default()
	return echomiddleware.RequestLoggerWithConfig(echomiddleware.RequestLoggerConfig{
		LogMethod:        true,
		LogURI:           true,
		LogContentLength: true,
		LogStatus:        true,
		LogResponseSize:  true,
		LogUserAgent:     true,
		LogRemoteIP:      true,
		LogReferer:       true,
		LogLatency:       true,
		LogProtocol:      true,
		LogError:         true,
		HandleError:      true,
		LogValuesFunc: func(c echo.Context, v echomiddleware.RequestLoggerValues) error {
			if v.Error == nil {
				httpRequest := HTTPRequest{
					RequestMethod: v.Method,
					RequestURL:    v.URI,
					RequestSize:   v.ContentLength,
					Status:        v.Status,
					ResponseSize:  fmt.Sprintf("%d", v.ResponseSize),
					UserAgent:     v.UserAgent,
					RemoteIP:      v.RemoteIP,
					Referer:       v.Referer,
					Latency:       MakeDuration(v.Latency),
					Protocol:      v.Protocol,
				}
				logger.LogAttrs(context.Background(), slog.LevelInfo, "request",
					slog.Any("httpRequest", httpRequest),
				)
			} else {
				httpRequest := HTTPRequest{
					RequestMethod: v.Method,
					RequestURL:    v.URI,
					RequestSize:   v.ContentLength,
					Status:        v.Status,
					ResponseSize:  fmt.Sprintf("%d", v.ResponseSize),
					UserAgent:     v.UserAgent,
					RemoteIP:      v.RemoteIP,
					Referer:       v.Referer,
					Latency:       MakeDuration(v.Latency),
					Protocol:      v.Protocol,
				}
				logger.LogAttrs(context.Background(), slog.LevelError, "request error",
					slog.Any("httpRequest", httpRequest),
					slog.Any("err", v.Error.Error()),
				)
			}

			return nil
		},
	})
}
