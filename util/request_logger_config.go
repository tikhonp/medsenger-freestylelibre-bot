package util

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func GetRequestLoggerConfig(enableJSONLog bool) middleware.RequestLoggerConfig {
	return middleware.RequestLoggerConfig{
		LogURI:           true,
		LogStatus:        true,
		LogMethod:        true,
		LogHost:          true,
		LogRemoteIP:      true,
		LogUserAgent:     true,
		LogContentLength: true,
		LogError:         true,
		LogResponseSize:  true,
		LogLatency:       true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if enableJSONLog {
				fmt.Printf(
					`{"time":"%s","remote_ip":"%s","host":"%s","method":"%s","uri":"%s","user_agent":"%s",`+
						`"status":%d,"error":"%e","latency":%d,"latency_human":"%s","bytes_out": %d}`+"\n",
					v.StartTime.Format("2006-01-02 15:04:05"),
					v.RemoteIP,
					v.Host,
					v.Method,
					v.URI,
					v.UserAgent,
					v.Status,
					v.Error,
					v.Latency.Nanoseconds(),
					v.Latency.String(),
					v.ResponseSize,
				)
			} else {
				fmt.Printf(
					"[%s] %d %s %s (%s) %s\n",
					v.StartTime.Format("2006-01-02 15:04:05"),
					v.Status,
					v.Method,
					v.URI,
					v.RemoteIP,
					v.Latency.String(),
				)
			}
			return nil
		},
	}
}
