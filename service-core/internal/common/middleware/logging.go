package appmiddleware

import (
	"net/http"
	"time"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	applogger "service-core/internal/common/logger"
)

func Logging(log applogger.Logger) Middleware {
	return func(next apphttp.AppHandler) apphttp.AppHandler {
		return func(w http.ResponseWriter, r *http.Request) error {
			start := time.Now()

			err := next(w, r)

			fields := []applogger.Field{
				{Key: "method", Value: r.Method},
				{Key: "path", Value: r.URL.Path},
				{Key: "duration_ms", Value: time.Since(start).Milliseconds()},
			}

			if err != nil {
				appErr := apperrors.Resolve(err)

				fields = append(fields,
					applogger.Field{Key: "error", Value: err.Error()},
					applogger.Field{Key: "type", Value: appErr.Type},
					applogger.Field{Key: "status_code", Value: appErr.StatusCode},
				)

				switch appErr.LogLevel {
				case applogger.LogLevelWarn:
					log.Warn(r.Context(), "request warning", fields...)

				case applogger.LogLevelInfo:
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
