package komerce

import (
	"testing"

	apperrors "service-core/internal/common/errors"
)

func TestMapKomerceError(t *testing.T) {
	tests := []struct {
		code         int
		rawMsg       string
		expectedType apperrors.ErrorType
		expectedCode int
	}{
		{400, "Invalid Api key, key not found", apperrors.ErrTypeUnauthorized, 401},
		{400, "API key is required", apperrors.ErrTypeUnauthorized, 401},
		{400, "Data waybill not found", apperrors.ErrTypeNotFound, 404},
		{400, "Resi tidak ditemukan", apperrors.ErrTypeNotFound, 404},
		{400, "Courier is not supported", apperrors.ErrTypeInvalidInput, 400},
		{429, "Too Many Requests", apperrors.ErrTypeTooManyRequests, 429},
		{500, "Internal Server Error", apperrors.ErrTypeInternal, 500},
		{400, "General invalid param", apperrors.ErrTypeBadRequest, 400},
	}

	for _, tt := range tests {
		appErr := mapKomerceError(tt.code, tt.rawMsg)
		if appErr == nil {
			t.Fatalf("mapKomerceError(%d, %q) returned nil", tt.code, tt.rawMsg)
		}
		if appErr.Type != tt.expectedType {
			t.Errorf("mapKomerceError(%d, %q).Type = %v; want %v", tt.code, tt.rawMsg, appErr.Type, tt.expectedType)
		}
		if appErr.StatusCode != tt.expectedCode {
			t.Errorf("mapKomerceError(%d, %q).StatusCode = %d; want %d", tt.code, tt.rawMsg, appErr.StatusCode, tt.expectedCode)
		}
	}
}
