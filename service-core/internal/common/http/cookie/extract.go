package appcookie

import (
	"errors"
	"net/http"
)

func CookieValue(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		if errors.Is(err, ErrNoCookie) {
			return "", ErrNoCookie
		}

		return "", err
	}

	if cookie.Value == "" {
		return "", ErrNoCookie
	}

	return cookie.Value, nil
}
