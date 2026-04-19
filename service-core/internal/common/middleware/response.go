package middleware

import (
	"net/http"

	"service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
)

func Response() Middleware {
	return func(next apphttp.AppHandler) apphttp.AppHandler {
		return func(w http.ResponseWriter, r *http.Request) error {
			err := next(w, r)

			if err != nil {
				appErr := errors.Map(err)

				apphttp.WriteJSON(w, appErr.StatusCode, map[string]any{
					"error": appErr.Message,
				})

				return nil
			}

			return nil
		}
	}
}
