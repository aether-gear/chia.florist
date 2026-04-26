package errors

import "service-core/internal/common/logger"

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
