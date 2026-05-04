package bootstrap

import (
	database "service-core/internal/infra/db"
	"service-core/internal/shared/config"
)

type Config struct {
	App      config.AppConfig
	Shipping config.ShippingConfig
	DB       database.DatabaseConfig
	JWT      config.JWTConfig
}

func LoadConfig() Config {
	return Config{
		App:      config.LoadAppConfig(),
		Shipping: config.LoadShippingConfig(),
		DB:       database.LoadConfig(),
		JWT:      config.LoadJWTConfig(),
	}
}
