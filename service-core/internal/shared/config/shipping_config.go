package config

import "time"

type ShippingConfig struct {
	DestinationURL string
	CalculateURL   string
	DestinationKey string
	CalculateKEY   string
	Timeout        time.Duration
}

func LoadShippingConfig() ShippingConfig {
	return ShippingConfig{
		DestinationURL: GetEnv("RAJAONGKIR_DESTINATION_URL"),
		CalculateURL:   GetEnv("RAJAONGKIR_CALCULATE_URL"),
		DestinationKey: GetEnv("RAJAONGKIR_SHIPPING"),
		CalculateKEY:   GetEnv("RAJAONGKIR_SHIPPING"),
		Timeout:        10 * time.Second,
	}
}
