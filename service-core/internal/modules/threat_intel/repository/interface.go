package repository

import "context"

type ThreatIntelProvider interface {
	GetReputation(
		ctx context.Context,
		ip string,
		apiKey string,
	) (map[string]any, error)

	GetGeolocation(
		ctx context.Context,
		ip string,
	) (map[string]any, error)
}
