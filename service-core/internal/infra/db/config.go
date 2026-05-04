package database

import "service-core/internal/shared/config"

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func LoadConfig() DatabaseConfig {
	return DatabaseConfig{
		Driver:   config.GetEnv("DB_DRIVER"),
		Host:     config.GetEnv("DB_HOST"),
		Port:     config.GetEnv("DB_PORT"),
		User:     config.GetEnv("DB_USER"),
		Password: config.GetEnv("DB_PASSWORD"),
		Name:     config.GetEnv("DB_NAME"),
		SSLMode:  config.GetEnv("DB_SSLMODE"),
	}
}

func (c DatabaseConfig) DSN() string {
	dsn := "postgresql://" +
		c.User + ":" +
		c.Password + "@" +
		c.Host + ":" +
		c.Port + "/" +
		c.Name

	if c.SSLMode != "" {
		dsn += "?sslmode=" + c.SSLMode
	}

	return dsn
}
