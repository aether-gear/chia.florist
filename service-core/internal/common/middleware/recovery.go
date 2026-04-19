package middleware

import (
	"fmt"
	"net/http"

	apphttp "service-core/internal/common/http"
	"service-core/internal/common/logger"
)

func Recovery(log logger.Logger) Middleware {
	return func(next apphttp.AppHandler) apphttp.AppHandler {
		return func(w http.ResponseWriter, r *http.Request) (err error) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Error(r.Context(), "panic recovered",
						logger.Field{
							Key:   "panic",
							Value: rec,
						},
					)
					err = fmt.Errorf("internal server error")
				}
			}()

			return next(w, r)
		}
	}
}
