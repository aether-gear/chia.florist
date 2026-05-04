package http

import (
	"net/http"
	"service-core/internal/common/errors"
)

type MethodHandler map[string]AppHandler

func HandleMethods(m MethodHandler) AppHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if h, ok := m[r.Method]; ok {
			return h(w, r)
		}
		return errors.ErrMethodNotAllowed
	}
}
