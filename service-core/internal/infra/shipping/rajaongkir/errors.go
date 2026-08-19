package rajaongkir

import (
	"errors"
	"fmt"

	apperrors "service-core/internal/common/errors"
)

var (
	ErrServerKeyRequired    = errors.New("raja ongkir: server key is required")
	ErrInvalidCourierOption = apperrors.NewInvalidInput("The selected courier option is invalid or no longer supported")
	ErrShippingUnavailable  = apperrors.NewBadRequest("Shipping is currently unavailable for this destination using the selected courier")
	ErrMissingParams        = apperrors.NewInternal(errors.New("Missing params"))
)

func mapRajaOngkirError(code int, rawMsg string) error {
	switch {
	case code == 422:
		return ErrInvalidCourierOption
	case code == 404:
		return ErrShippingUnavailable
	case code == 400 || rawMsg == "Missing Params":
		return ErrMissingParams
	default:
		return fmt.Errorf("rajaongkir error (%d): %s", code, rawMsg)
	}
}
