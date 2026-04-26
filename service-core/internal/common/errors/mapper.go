package errors

import (
	"errors"
	"net/http"
	"service-core/internal/common/logger"
	"strings"
)

func Map(err error) *AppError {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, ErrBadRequest):
		return &AppError{
			Type:       ErrTypeBadRequest,
			Message:    "bad request",
			StatusCode: http.StatusBadRequest,
			LogLevel:   logger.LogLevelWarn,
			Err:        err,
		}

	case errors.Is(err, ErrInvalidInput):
		return &AppError{
			Type:       ErrTypeInvalidInput,
			Message:    "invalid input",
			StatusCode: http.StatusBadRequest,
			LogLevel:   logger.LogLevelWarn,
			Err:        err,
		}

	case errors.Is(err, ErrNotFound):
		return &AppError{
			Type:       ErrTypeNotFound,
			Message:    "resource not found",
			StatusCode: http.StatusNotFound,
			LogLevel:   logger.LogLevelInfo,
			Err:        err,
		}

	case errors.Is(err, ErrUnauthorized):
		return &AppError{
			Type:       ErrTypeUnauthorized,
			Message:    "unauthorized",
			StatusCode: http.StatusUnauthorized,
			LogLevel:   logger.LogLevelWarn,
			Err:        err,
		}

	case errors.Is(err, ErrForbidden):
		return &AppError{
			Type:       ErrTypeForbidden,
			Message:    "forbidden",
			StatusCode: http.StatusForbidden,
			LogLevel:   logger.LogLevelWarn,
			Err:        err,
		}

	case strings.Contains(err.Error(), "duplicate key"):
		return &AppError{
			Type:       ErrTypeConflict,
			Message:    "conflict",
			StatusCode: http.StatusConflict,
			LogLevel:   logger.LogLevelWarn,
			Err:        err,
		}

	case errors.Is(err, ErrMethodNotAllowed):
		return &AppError{
			Type:       ErrTypeMethodNotAllowed,
			Message:    "method not allowed",
			StatusCode: http.StatusMethodNotAllowed,
			LogLevel:   logger.LogLevelWarn,
			Err:        err,
		}
	}

	return &AppError{
		Type:       ErrTypeInternal,
		Message:    "internal server error",
		StatusCode: http.StatusInternalServerError,
		LogLevel:   logger.LogLevelError,
		Err:        err,
	}
}
