package errors

import "net/http"

type AppError struct {
	Message    string
	StatusCode int
}

func Map(err error) AppError {
	switch err {
	case ErrNotFound:
		return AppError{"not found", http.StatusNotFound}
	case ErrUnauthorized:
		return AppError{"unauthorized", http.StatusUnauthorized}
	default:
		return AppError{"internal error", http.StatusInternalServerError}
	}
}
