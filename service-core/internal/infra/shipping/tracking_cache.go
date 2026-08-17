package shipping

import (
	"sync"
	"time"
)

const DefaultTrackingCacheTTL = 5 * time.Minute

type cachedTrackingEntry struct {
	events    []TrackingEvent
	fetchedAt time.Time
}

// TrackingCache provides an in-memory thread-safe TTL cache for tracking events.
type TrackingCache struct {
	mu    sync.RWMutex
	ttl   time.Duration
	items map[string]cachedTrackingEntry
}

func NewTrackingCache(ttl time.Duration) *TrackingCache {
	if ttl <= 0 {
		ttl = DefaultTrackingCacheTTL
	}
	return &TrackingCache{
		ttl:   ttl,
		items: make(map[string]cachedTrackingEntry),
	}
}

func (c *TrackingCache) cacheKey(courier, trackingNo string) string {
	return courier + ":" + trackingNo
}

func (c *TrackingCache) Get(courier, trackingNo string) ([]TrackingEvent, bool) {
	if c == nil {
		return nil, false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	key := c.cacheKey(courier, trackingNo)
	entry, ok := c.items[key]
	if !ok {
		return nil, false
	}

	if time.Since(entry.fetchedAt) > c.ttl {
		return nil, false
	}

	return entry.events, true
}

func (c *TrackingCache) GetStale(courier, trackingNo string) ([]TrackingEvent, bool) {
	if c == nil {
		return nil, false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	key := c.cacheKey(courier, trackingNo)
	entry, ok := c.items[key]
	if !ok {
		return nil, false
	}

	return entry.events, true
}

func (c *TrackingCache) Set(courier, trackingNo string, events []TrackingEvent) {
	if c == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.cacheKey(courier, trackingNo)
	c.items[key] = cachedTrackingEntry{
		events:    events,
		fetchedAt: time.Now(),
	}
}
