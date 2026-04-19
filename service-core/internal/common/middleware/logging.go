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

			log.Info(r.Context(), "request",
				logger.Field{
					Key:   "method",
					Value: r.Method,
				},
				logger.Field{
					Key:   "path",
					Value: r.URL.Path,
				},
				logger.Field{
					Key:   "duration_ms",
					Value: time.Since(start).Milliseconds(),
				},
				logger.Field{
					Key:   "error",
					Value: err,
				},
			)

			return err
		}
	}
}
