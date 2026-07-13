package config

// KomerceConfig holds credentials and base URLs for the Komerce Collaborator
// logistics platform, which is part of the RajaOngkir ecosystem.
//
// Production:  https://api.collaborator.komerce.id
// Sandbox:     https://api-sandbox.collaborator.komerce.id
type KomerceConfig struct {
	OrderBaseURL string
	TrackBaseURL string
	APIKey       string
	ShippingKey  string
}

func LoadKomerceConfig() KomerceConfig {
	return KomerceConfig{
		OrderBaseURL: GetEnv("KOMERCE_ORDER_URL"),
		TrackBaseURL: GetEnv("KOMERCE_TRACK_URL"),
		APIKey:       GetEnv("KOMERCE_API_KEY"),
		ShippingKey:  GetEnv("KOMERCE_SHIPPING_DELIVERY"),
	}
}
