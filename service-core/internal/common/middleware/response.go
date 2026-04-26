package middleware

import (
	"net/http"
	"strings"

	"service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
)

func Response() Middleware {
	return func(next apphttp.AppHandler) apphttp.AppHandler {
		return func(w http.ResponseWriter, r *http.Request) error {
			err := next(w, r)

			if err != nil {
				appErr := errors.Resolve(err)

				errors := []string{firstMessage(err), lastMessage(err)}
				messDebug := strings.Join(errors, ": ")
				errRes := ErrResponse{
					Type:       appErr.Type,
					Message:    firstMessage(err),
					Debug:      &messDebug,
					StatusCode: appErr.StatusCode,
				}

				apphttp.WriteJSON(w, appErr.StatusCode, errRes)

				return err
			}

			return nil
		}
	}
}

func firstMessage(err error) string {
	msg := err.Error()

	parts := strings.Split(msg, ":")
	if len(parts) > 0 {
		return strings.TrimSpace(parts[0])
	}

	return msg
}

func lastMessage(err error) string {
	msg := err.Error()

	parts := strings.Split(msg, ":")
	if len(parts) > 0 {
		return strings.TrimSpace(parts[len(parts)-1])
	}

	return msg
}

type ErrResponse struct {
	Type       errors.ErrorType `json:"type"`
	StatusCode int              `json:"status_code"`
	Message    string           `json:"message"`
	Debug      *string          `json:"debug"`
}
