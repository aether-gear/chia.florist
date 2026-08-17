package manual

import (
	"context"

	shipping "service-core/internal/infra/shipping"
)

// manualTrackedShippingProvider is a hybrid implementation of shipping.LogisticsProvider.
// It bypasses external order creation (echoing staff's manually entered tracking number),
// but delegates tracking queries (TrackShipment) to an external logistics provider (e.g. Komerce).
type manualTrackedShippingProvider struct {
	tracker shipping.LogisticsProvider
}

// NewManualTrackedShippingProvider returns a LogisticsProvider that accepts manual tracking
// numbers during dispatch while delegating live tracking queries to the provided tracker.
func NewManualTrackedShippingProvider(tracker shipping.LogisticsProvider) shipping.LogisticsProvider {
	return &manualTrackedShippingProvider{
		tracker: tracker,
	}
}

// CreateOrder echoes ManualTrackingNumber from input as the result without calling any
// external booking API so staff can enter tracking numbers manually.
func (p *manualTrackedShippingProvider) CreateOrder(
	_ context.Context,
	input shipping.CreateOrderInput,
) (*shipping.CreateOrderResult, error) {
	tracking := ""
	if input.ManualTrackingNumber != nil {
		tracking = *input.ManualTrackingNumber
	}

	return &shipping.CreateOrderResult{
		KomerceOrderNo: "",
		TrackingNumber: tracking,
	}, nil
}

// TrackShipment delegates tracking lookup to the external tracker (e.g. Komerce Waybill API).
func (p *manualTrackedShippingProvider) TrackShipment(
	ctx context.Context,
	input shipping.TrackShipmentInput,
) ([]shipping.TrackingEvent, error) {
	if p.tracker == nil {
		return []shipping.TrackingEvent{}, nil
	}
	return p.tracker.TrackShipment(ctx, input)
}
