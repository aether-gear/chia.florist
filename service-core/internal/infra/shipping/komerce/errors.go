package komerce

import (
	"errors"
	"fmt"
	"strings"

	apperrors "service-core/internal/common/errors"
)

// Sentinel errors for Komerce provider configuration and initialization.
var (
	ErrAPIKeyRequired      = errors.New("komerce: api key is required")
	ErrOrderBaseURLRequired = errors.New("komerce: order base URL is required")
	ErrTrackBaseURLRequired = errors.New("komerce: track base URL is required")
)

// ErrInvalidCourierCode returns a bad request error for unsupported/invalid courier code.
func ErrInvalidCourierCode(courier string) *apperrors.AppError {
	return apperrors.NewBadRequest(fmt.Sprintf("komerce courier code %q is invalid or unsupported", courier))
}

// ErrTrackingNotFound returns a not found error when a waybill is missing or not yet in courier system.
func ErrTrackingNotFound(trackingNumber string, detail string) *apperrors.AppError {
	if strings.TrimSpace(detail) != "" {
		return apperrors.NewNotFound(fmt.Sprintf("tracking number %q not found or not yet scanned by courier: %s", trackingNumber, detail))
	}
	return apperrors.NewNotFound(fmt.Sprintf("tracking number %q not found or not yet scanned by courier", trackingNumber))
}

// ErrAPIUnauthorized returns an unauthorized error when Komerce API key is invalid or rejected.
func ErrAPIUnauthorized(detail string) *apperrors.AppError {
	if strings.TrimSpace(detail) != "" {
		return apperrors.NewUnauthorized(fmt.Sprintf("invalid or missing Komerce API key: %s", detail))
	}
	return apperrors.NewUnauthorized("invalid or missing Komerce API key in backend config")
}

// ErrAPIRateLimit returns a rate limit exceeded error.
func ErrAPIRateLimit(statusCode int, detail string) *apperrors.AppError {
	if strings.TrimSpace(detail) != "" {
		return apperrors.NewTooManyRequests(fmt.Sprintf("Komerce API rate limit exceeded (%d): %s", statusCode, detail))
	}
	return apperrors.NewTooManyRequests("Komerce API rate limit exceeded")
}

// ErrHTTPStatus maps raw HTTP error status codes from Komerce API into structured application errors.
func ErrHTTPStatus(statusCode int, body string) error {
	switch statusCode {
	case 400, 422:
		return apperrors.NewBadRequest(fmt.Sprintf("komerce API error (%d): %s", statusCode, body))
	case 401:
		return ErrAPIUnauthorized(body)
	case 404:
		return apperrors.NewNotFound(fmt.Sprintf("komerce resource not found (%d): %s", statusCode, body))
	case 429:
		return ErrAPIRateLimit(statusCode, body)
	case 500, 502, 503, 504:
		return fmt.Errorf("komerce provider service unavailable (%d): %s", statusCode, body)
	default:
		return fmt.Errorf("unexpected komerce HTTP status %d: %s", statusCode, body)
	}
}

// mapKomerceError converts Komerce API meta status code and message into system *apperrors.AppError.
func mapKomerceError(code int, rawMsg string) *apperrors.AppError {
	msg := strings.ToLower(rawMsg)

	switch {
	case strings.Contains(msg, "api key") || strings.Contains(msg, "key not found") || strings.Contains(msg, "unauthorized"):
		return ErrAPIUnauthorized(rawMsg)

	case strings.Contains(msg, "waybill not found") || strings.Contains(msg, "resi tidak ditemukan") || strings.Contains(msg, "data waybill not found") || strings.Contains(msg, "not found"):
		return apperrors.NewNotFound("Tracking number not found or not yet scanned by courier: " + rawMsg)

	case code == 400 || (strings.Contains(msg, "courier") && (strings.Contains(msg, "not supported") || strings.Contains(msg, "invalid") || strings.Contains(msg, "tidak didukung"))):
		return apperrors.NewBadRequest(fmt.Sprintf("Komerce API error (%d): %s", code, rawMsg))

	case code == 429 || strings.Contains(msg, "too many requests") || strings.Contains(msg, "rate limit"):
		return ErrAPIRateLimit(code, rawMsg)

	case code >= 500:
		return apperrors.NewInternal(fmt.Errorf("komerce provider internal error (%d): %s", code, rawMsg))

	default:
		return apperrors.NewBadRequest(fmt.Sprintf("Komerce API error (%d): %s", code, rawMsg))
	}
}
