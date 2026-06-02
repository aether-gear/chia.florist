package appmiddleware

import (
	"net/http"
	"runtime/debug"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	applogger "service-core/internal/common/logger"
)

func Recovery(log applogger.Logger) Middleware {
	return func(next apphttp.AppHandler) apphttp.AppHandler {
		return func(w http.ResponseWriter, r *http.Request) (err error) {
			defer func() {
				if rec := recover(); rec != nil {
					stack := debug.Stack()

					log.Error(r.Context(), "panic recovered",
						applogger.Field{Key: "panic", Value: rec},
						applogger.Field{Key: "stack", Value: string(stack)},
						applogger.Field{Key: "path", Value: r.URL.Path},
						applogger.Field{Key: "method", Value: r.Method},
					)

					err = apperrors.ErrInternal
				}
			}()

			return next(w, r)
		}
	}
}
