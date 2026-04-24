package middleware

import (
	"net/http"
	"runtime/debug"

	"service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	"service-core/internal/common/logger"
)

func Recovery(log logger.Logger) Middleware {
	return func(next apphttp.AppHandler) apphttp.AppHandler {
		return func(w http.ResponseWriter, r *http.Request) (err error) {
			defer func() {
				if rec := recover(); rec != nil {
					stack := debug.Stack()

					log.Error(r.Context(), "panic recovered",
						logger.Field{Key: "panic", Value: rec},
						logger.Field{Key: "stack", Value: string(stack)},
						logger.Field{Key: "path", Value: r.URL.Path},
						logger.Field{Key: "method", Value: r.Method},
					)

					err = errors.ErrInternal
				}
			}()

			return next(w, r)
		}
	}
}
