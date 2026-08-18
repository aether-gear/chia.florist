package komerce

import (
	"fmt"
	"strings"

	apperrors "service-core/internal/common/errors"
)

// mapKomerceError converts Komerce API meta status code and message into system *apperrors.AppError.
func mapKomerceError(code int, rawMsg string) *apperrors.AppError {
	msg := strings.ToLower(rawMsg)

	switch {
	case strings.Contains(msg, "api key") || strings.Contains(msg, "key not found") || strings.Contains(msg, "unauthorized"):
		return apperrors.NewUnauthorized("Invalid or missing Komerce API key in backend config: " + rawMsg)

	case strings.Contains(msg, "waybill not found") || strings.Contains(msg, "resi tidak ditemukan") || strings.Contains(msg, "data waybill not found") || strings.Contains(msg, "not found"):
		return apperrors.NewNotFound("Tracking number not found or not yet scanned by courier: " + rawMsg)

	case strings.Contains(msg, "courier") && (strings.Contains(msg, "not supported") || strings.Contains(msg, "invalid") || strings.Contains(msg, "tidak didukung")):
		return apperrors.NewInvalidInput("Courier code is not supported by Komerce tracking: " + rawMsg)

	case code == 429 || strings.Contains(msg, "too many requests") || strings.Contains(msg, "rate limit"):
		return apperrors.NewTooManyRequests("Komerce API rate limit exceeded: " + rawMsg)

	case code >= 500:
		return apperrors.NewInternal(fmt.Errorf("komerce provider internal error (%d): %s", code, rawMsg))

	default:
		return apperrors.NewBadRequest("Komerce API request failed: " + rawMsg)
	}
}
