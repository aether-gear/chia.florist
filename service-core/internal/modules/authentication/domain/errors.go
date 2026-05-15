package domain

import "errors"

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrAccountAlreadyExists = errors.New("account already exists")
