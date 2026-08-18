package manual

import (
	"context"
	"errors"
	"testing"
	"time"

	shipping "service-core/internal/infra/shipping"
)

type mockTracker struct {
	events []shipping.TrackingEvent
	err    error
}

func (m *mockTracker) CreateOrder(_ context.Context, _ shipping.CreateOrderInput) (*shipping.CreateOrderResult, error) {
	return nil, nil
}

func (m *mockTracker) TrackShipment(_ context.Context, _ shipping.TrackShipmentInput) ([]shipping.TrackingEvent, error) {
	return m.events, m.err
}

func TestManualTrackedShippingProvider_CreateOrder(t *testing.T) {
	provider := NewManualTrackedShippingProvider(&mockTracker{})

	manualNo := "MANUAL-AWB-999"
	res, err := provider.CreateOrder(context.Background(), shipping.CreateOrderInput{
		ManualTrackingNumber: &manualNo,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.TrackingNumber != manualNo {
		t.Errorf("expected tracking number %s, got %s", manualNo, res.TrackingNumber)
	}
	if res.KomerceOrderNo != "" {
		t.Errorf("expected empty komerce order no, got %s", res.KomerceOrderNo)
	}
}

func TestManualTrackedShippingProvider_TrackShipment(t *testing.T) {
	now := time.Now()
	expectedEvents := []shipping.TrackingEvent{
		{
			Status:      "MANIFESTED",
			Description: "Package received at hub",
			Location:    "JAKARTA",
			Timestamp:   now,
		},
	}

	mock := &mockTracker{events: expectedEvents}
	provider := NewManualTrackedShippingProvider(mock)

	events, err := provider.TrackShipment(context.Background(), shipping.TrackShipmentInput{
		Courier:        "jne",
		TrackingNumber: "MANUAL-AWB-999",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Status != "MANIFESTED" {
		t.Errorf("expected status MANIFESTED, got %s", events[0].Status)
	}
}

func TestManualTrackedShippingProvider_TrackShipment_NilTracker(t *testing.T) {
	provider := NewManualTrackedShippingProvider(nil)

	events, err := provider.TrackShipment(context.Background(), shipping.TrackShipmentInput{
		Courier:        "jne",
		TrackingNumber: "MANUAL-AWB-999",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("expected 0 events for nil tracker, got %d", len(events))
	}
}

func TestManualTrackedShippingProvider_TrackShipment_Error(t *testing.T) {
	expectedErr := errors.New("external provider down")
	mock := &mockTracker{err: expectedErr}
	provider := NewManualTrackedShippingProvider(mock)

	_, err := provider.TrackShipment(context.Background(), shipping.TrackShipmentInput{
		Courier:        "jne",
		TrackingNumber: "MANUAL-AWB-999",
	})
	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}
