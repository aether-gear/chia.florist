package appmiddleware

import (
	"net/http"
	"strings"
	"time"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	applogger "service-core/internal/common/logger"

	"github.com/google/uuid"
)

func Logging(log applogger.Logger) Middleware {
	return func(next apphttp.AppHandler) apphttp.AppHandler {
		return func(w http.ResponseWriter, r *http.Request) error {
			requestID := uuid.NewString()
			ctx := applogger.WithRequestID(r.Context(), requestID)
			ctx = applogger.WithClientIP(ctx, apphttp.ClientIP(r))
			r = r.WithContext(ctx)

			w.Header().Set("X-Request-ID", requestID)

			rw := newResponseRecorder(w)

			start := time.Now()
			err := next(rw, r)

			// request_id and client_ip are already in ctx and prepended
			// automatically by the Logger implementations.
			fields := []applogger.Field{
				{Key: "layer", Value: applogger.LayerMiddleware},
				{Key: "method", Value: r.Method},
				{Key: "path", Value: r.URL.Path},
				{Key: "status_code", Value: rw.statusCode},
				{Key: "user_agent", Value: r.UserAgent()},
				{Key: "duration_ms", Value: time.Since(start).Milliseconds()},
				{Key: "category", Value: applogger.CategorySystem},
			}

			if err != nil {
				appErr := apperrors.Resolve(err)

				fields = append(fields,
					applogger.Field{Key: "error", Value: err.Error()},
					applogger.Field{Key: "type", Value: appErr.Type},
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
				// Mute successful logs for admin polling
				// api audit and health routes to keep the console clean
				if !strings.HasPrefix(r.URL.Path, "/api/") && r.URL.Path != "/health" {
					log.Info(r.Context(), "request success", fields...)
				}
			}

			return err
		}
	}
}
