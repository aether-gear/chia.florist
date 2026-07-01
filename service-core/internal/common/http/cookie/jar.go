package appcookie

import (
	"errors"
	"net/http"
	"time"
)

var (
	// ErrNoCookie is returned when
	// a requested cookie is not present in the HTTP request
	ErrNoCookie = errors.New("http: named cookie not present")
)

type CookieName string

const (
	// CookieAccess is the cookie name used to
	// store the user's access token
	CookieAccess CookieName = "chast"

	// CookieStaff is the cookie name used for
	// staff authentication sessions
	CookieStaff CookieName = "hotpot"

	// CookieOAuthState is the cookie name used to
	// store OAuth state during authentication flow
	CookieOAuthState CookieName = "cinnamon-bun"
)

// Bind sets an HTTP cookie on the response writer
// with the given name, value, and expiration time
//
// The cookie is configured as HttpOnly, Secure,
// and uses SameSiteLaxMode by default
func Bind(
	w http.ResponseWriter,
	name CookieName,
	value string,
	exp time.Time,
) {
	http.SetCookie(w, &http.Cookie{
		Name:     string(name),
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  exp,
	})
}

// Extract retrieves a cookie value from the
// incoming HTTP request by name
//
// It returns ErrNoCookie if the cookie
// is not present
func Extract(
	r *http.Request,
	coukieName CookieName,
) (string, error) {
	cookie, err := r.Cookie(string(coukieName))
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

// Clear removes a cookie from the client by
// setting an expired cookie with the same name
//
// It effectively instructs the browser to
// delete the cookie
func Clear(
	w http.ResponseWriter,
	name CookieName,
) {
	http.SetCookie(w, &http.Cookie{
		Name:     string(name),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
}

// Exists checks whether a cookie with the
// given name exists in the HTTP request
//
// It returns true if the cookie is present,
// false otherwise
func Exists(
	r *http.Request,
	name CookieName,
) bool {
	_, err := r.Cookie(string(name))
	return err == nil
}
