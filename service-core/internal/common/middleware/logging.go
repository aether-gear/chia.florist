package middleware

import (
	"net/http"
	apphttp "service-core/internal/common/http"
	"service-core/internal/common/logger"
	"time"
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
				fields = append(fields,
					logger.Field{Key: "error", Value: err.Error()},
				)
				log.Error(r.Context(), "request failed", fields...)
			} else {
				log.Info(r.Context(), "request success", fields...)
			}

			return err
		}
	}
}
