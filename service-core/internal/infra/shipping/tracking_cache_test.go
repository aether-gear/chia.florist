package shipping

import (
	"testing"
	"time"
)

func TestTrackingCache(t *testing.T) {
	ttl := 100 * time.Millisecond
	cache := NewTrackingCache(ttl)

	courier := "jne"
	trackingNo := "TRACK123"

	events := []TrackingEvent{
		{
			Status:      "MANIFESTED",
			Description: "Package received",
			Location:    "JAKARTA",
			Timestamp:   time.Now(),
		},
	}

	// 1. Initial get should miss
	if _, ok := cache.Get(courier, trackingNo); ok {
		t.Fatal("expected cache miss initially")
	}

	// 2. Set and get should hit
	cache.Set(courier, trackingNo, events)
	cachedEvents, ok := cache.Get(courier, trackingNo)
	if !ok {
		t.Fatal("expected cache hit after set")
	}
	if len(cachedEvents) != 1 || cachedEvents[0].Status != "MANIFESTED" {
		t.Errorf("unexpected cached events: %v", cachedEvents)
	}

	// 3. Wait for TTL expiration
	time.Sleep(150 * time.Millisecond)

	// 4. Get should miss after TTL
	if _, ok := cache.Get(courier, trackingNo); ok {
		t.Fatal("expected cache miss after TTL expiration")
	}

	// 5. GetStale should still return stale items if expired
	staleEvents, ok := cache.GetStale(courier, trackingNo)
	if !ok {
		t.Fatal("expected GetStale hit after TTL expiration")
	}
	if len(staleEvents) != 1 {
		t.Errorf("expected 1 stale event, got %d", len(staleEvents))
	}
}
