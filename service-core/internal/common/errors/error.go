package errors

import (
	"errors"

	"service-core/internal/common/logger"
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
	Type       ErrorType       `json:"type"`
	Message    string          `json:"message"`
	StatusCode int             `json:"status_code"`
	LogLevel   logger.LogLevel `json:"-"`
	Err        error           `json:"-"`
}

func newAppError(
	t ErrorType,
	msg string,
	code int,
	level logger.LogLevel,
	err error,
) *AppError {
	return &AppError{
		Type:       t,
		Message:    msg,
		StatusCode: code,
		LogLevel:   level,
		Err:        err,
	}
}

func NewBadRequest(msg string) *AppError {
	return newAppError(ErrTypeBadRequest, msg, 400, logger.LogLevelWarn, nil)
}

func NewInvalidInput(msg string) *AppError {
	return newAppError(ErrTypeInvalidInput, msg, 400, logger.LogLevelWarn, nil)
}

func NewUnauthorized(msg string) *AppError {
	return newAppError(ErrTypeUnauthorized, msg, 401, logger.LogLevelWarn, nil)
}

func NewForbidden(msg string) *AppError {
	return newAppError(ErrTypeForbidden, msg, 403, logger.LogLevelWarn, nil)
}

func NewNotFound(msg string) *AppError {
	return newAppError(ErrTypeNotFound, msg, 404, logger.LogLevelWarn, nil)
}

func NewConflict(msg string) *AppError {
	return newAppError(ErrTypeConflict, msg, 409, logger.LogLevelWarn, nil)
}

func NewMethodNotAllowed(msg string) *AppError {
	return newAppError(ErrTypeMethodNotAllowed, msg, 405, logger.LogLevelWarn, nil)
}

func NewInternal(err error) *AppError {
	return newAppError(ErrTypeInternal, "internal server error", 500, logger.LogLevelError, err)
}
