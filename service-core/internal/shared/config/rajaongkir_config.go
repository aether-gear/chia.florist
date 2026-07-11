package config

type RajaOngkirConfig struct {
	URL                 string
	ShippingCostKey     string
	ShippingDeliveryKey string
	PaymentKey          string
	QRISLYKey           string
}

func LoadRajaOngkirConfig() RajaOngkirConfig {
	return RajaOngkirConfig{
		URL:                 GetEnv("RAJAONGKIR_URL"),
		ShippingCostKey:     GetEnv("RAJAONGKIR_SHIPPING_COST"),
		ShippingDeliveryKey: GetEnv("RAJAONGKIR_SHIPPING_DELIVERY"),
		PaymentKey:          GetEnv("RAJAONGKIR_PAYMENT"),
		QRISLYKey:           GetEnv("RAJAONGKIR_QRISLY"),
	}
}
