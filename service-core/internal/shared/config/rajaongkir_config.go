package config

type RajaOngkirConfig struct {
	URL         string
	ShippingKey string
	PaymentKey  string
	QRISLYKey   string
}

func LoadRajaOngkirConfig() RajaOngkirConfig {
	return RajaOngkirConfig{
		URL:         GetEnv("RAJAONGKIR_URL"),
		ShippingKey: GetEnv("RAJAONGKIR_SHIPPING"),
		PaymentKey:  GetEnv("RAJAONGKIR_PAYMENT"),
		QRISLYKey:   GetEnv("RAJAONGKIR_QRISLY"),
	}
}
