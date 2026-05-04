package middleware

import (
	"net/http"
	"time"

	"service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	"service-core/internal/common/logger"
)

func Logging(log logger.Logger) Middleware {
	return func(next apphttp.AppHandler) apphttp.AppHandler {
		return func(w http.ResponseWriter, r *http.Request) error {
			start := time.Now()

			err := next(w, r)

			fields := []logger.Field{
				{Key: "method", Value: r.Method},
				{Key: "path", Value: r.URL.Path},
				{Key: "duration_ms", Value: time.Since(start).Milliseconds()},
			}

			if err != nil {
				appErr := errors.Resolve(err)

				fields = append(fields,
					logger.Field{Key: "error", Value: err.Error()},
					logger.Field{Key: "type", Value: appErr.Type},
					logger.Field{Key: "status_code", Value: appErr.StatusCode},
				)

				switch appErr.LogLevel {
				case logger.LogLevelWarn:
					log.Warn(r.Context(), "request warning", fields...)

				case logger.LogLevelInfo:
					log.Info(r.Context(), "request info", fields...)

				default:
					log.Error(r.Context(), "request failed", fields...)
				}

			} else {
				log.Info(r.Context(), "request success", fields...)
			}

			return err
		}
	}
}
