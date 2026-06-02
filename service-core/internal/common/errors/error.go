package apperrors

import (
	"errors"

	applogger "service-core/internal/common/logger"
)

type ErrorType string

const (
	ErrTypeBadRequest       ErrorType = "BAD_REQUEST"
	ErrTypeInvalidInput     ErrorType = "INVALID_INPUT"
	ErrTypeNotFound         ErrorType = "NOT_FOUND"
	ErrTypeUnauthorized     ErrorType = "UNAUTHORIZED"
	ErrTypeForbidden        ErrorType = "FORBIDDEN"
	ErrTypeConflict         ErrorType = "CONFLICT"
	ErrTypeMethodNotAllowed ErrorType = "METHOD_NOT_ALLOWED"
	ErrTypeInternal         ErrorType = "INTERNAL_ERROR"
)

var (
	ErrBadRequest       = errors.New("bad request")
	ErrInvalidInput     = errors.New("invalid input")
	ErrNotFound         = errors.New("resource not found")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrConflict         = errors.New("conflict")
	ErrMethodNotAllowed = errors.New("method not allowed")
	ErrInternal         = errors.New("internal server error")
)

type AppError struct {
	Type       ErrorType          `json:"type"`
	Message    string             `json:"message"`
	StatusCode int                `json:"status_code"`
	LogLevel   applogger.LogLevel `json:"-"`
	Err        error              `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}
