package appcookie

import "errors"

var (
	ErrNoCookie = errors.New("http: named cookie not present")
)
