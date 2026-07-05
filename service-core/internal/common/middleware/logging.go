package appmiddleware

import (
	"net/http"
	"time"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	applogger "service-core/internal/common/logger"

	"github.com/google/uuid"
)

func Logging(log applogger.Logger) Middleware {
	return func(next apphttp.AppHandler) apphttp.AppHandler {
		return func(w http.ResponseWriter, r *http.Request) error {
			// Stamp a unique request ID and client IP onto the context.
			// All downstream log calls (system + future audit) will pick
			// these up automatically via FieldsFromContext.
			requestID := uuid.NewString()
			ctx := applogger.WithRequestID(r.Context(), requestID)
			ctx = applogger.WithClientIP(ctx, r.RemoteAddr)
			r = r.WithContext(ctx)

			// Expose the request ID in the response header for client-side tracing.
			w.Header().Set("X-Request-ID", requestID)

			start := time.Now()

			err := next(w, r)

			// request_id and client_ip are already in ctx and will be prepended
			// automatically by the Logger implementations — no need to add them here.
			fields := []applogger.Field{
				{Key: "method",      Value: r.Method},
				{Key: "path",        Value: r.URL.Path},
				{Key: "user_agent",  Value: r.UserAgent()},
				{Key: "duration_ms", Value: time.Since(start).Milliseconds()},
				{Key: "category",    Value: applogger.CategorySystem},
			}

			if err != nil {
				appErr := apperrors.Resolve(err)

				fields = append(fields,
					applogger.Field{Key: "error",       Value: err.Error()},
					applogger.Field{Key: "type",        Value: appErr.Type},
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
