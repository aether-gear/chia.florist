package database

import "service-core/internal/shared/config"

type Config struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func LoadConfig() Config {
	return Config{
		Driver:   config.MustGetEnv("DB_DRIVER"),
		Host:     config.MustGetEnv("DB_HOST"),
		Port:     config.MustGetEnv("DB_PORT"),
		User:     config.MustGetEnv("DB_USER"),
		Password: config.MustGetEnv("DB_PASSWORD"),
		Name:     config.MustGetEnv("DB_NAME"),
		SSLMode:  config.MustGetEnv("DB_SSLMODE"),
	}
}

func (c Config) DSN() string {
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
