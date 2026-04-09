package database

import (
	"fmt"
	"os"
)

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
		Driver:   mustGetEnv("DB_DRIVER"),
		Host:     mustGetEnv("DB_HOST"),
		Port:     mustGetEnv("DB_PORT"),
		User:     mustGetEnv("DB_USER"),
		Password: mustGetEnv("DB_PASSWORD"),
		Name:     mustGetEnv("DB_NAME"),
		SSLMode:  mustGetEnv("DB_SSLMODE"),
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

func mustGetEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		panic(fmt.Sprintf("Missing required environment variable: %s", key))
	}
	return val
}
