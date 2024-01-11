package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

type AccessLog struct {
	Severity    string      `json:"severity"`
	Time        time.Time   `json:"time"`
	Trace       string      `json:"logging.googleapis.com/trace"`
	HTTPRequest HTTPRequest `json:"httpRequest"`
}

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

func Logger() echo.MiddlewareFunc {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.LevelKey:
				return slog.String("severity", a.Value.String())
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
		LogProtocol:     true,
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
				AccessLog := AccessLog{
					Severity: "INFO",
					Time:     time.Now(),
					Trace:    "",
					HTTPRequest: HTTPRequest{
						RequestMethod: v.Method,
						RequestURL:    v.URI,
						RequestSize:   v.ContentLength,
						Status:        v.Status,
						ResponseSize:  fmt.Sprint(v.ResponseSize),
						UserAgent:     v.UserAgent,
						RemoteIP:      v.RemoteIP,
						ServerIP:      v.RemoteIP,
						Referer:       v.Referer,
						Latency:       Duration{},
						Protocol:      v.Protocol,
					},
				}
				log, _ := json.Marshal(AccessLog)
				fmt.Println(string(log))
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
