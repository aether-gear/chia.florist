package config

import "time"

type ShippingConfig struct {
	BaseURL        string
	DestinationURL string
	Timeout        time.Duration
}

func LoadShippingConfig() ShippingConfig {
	return ShippingConfig{
		BaseURL:        GetEnv("KOMERCE_SHIPPING"),
		DestinationURL: GetEnv("KOMERCE_DESTINATION_URL"),
		Timeout:        10 * time.Second,
	}
}
