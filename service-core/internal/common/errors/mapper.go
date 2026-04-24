package errors

import (
	"errors"
	"net/http"
)

type AppError struct {
	Message    string
	StatusCode int
	Err        error
}

func Map(err error) AppError {
	if err != nil {
		return AppError{}
	}

	switch {
	case errors.Is(err, ErrBadRequest):
		return AppError{
			Message:    "bad request",
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}

	case errors.Is(err, ErrNotFound):
		return AppError{
			Message:    "resource not found",
			StatusCode: http.StatusNotFound,
			Err:        err,
		}

	case errors.Is(err, ErrUnauthorized):
		return AppError{
			Message:    "unauthorized",
			StatusCode: http.StatusUnauthorized,
			Err:        err,
		}

	case errors.Is(err, ErrForbidden):
		return AppError{
			Message:    "forbidden",
			StatusCode: http.StatusForbidden,
			Err:        err,
		}

	case errors.Is(err, ErrInvalidInput):
		return AppError{
			Message:    "invalid input",
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}

	case errors.Is(err, ErrMethodNotAllowed):
		return AppError{
			Message:    "method not allowed",
			StatusCode: http.StatusMethodNotAllowed,
			Err:        err,
		}

	default:
		return AppError{
			Message:    "internal server error",
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
}
