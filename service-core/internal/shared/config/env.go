package config

import (
	"fmt"
	"os"
)

func MustGetEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		panic(fmt.Sprintf("Missing required environment variable: %s", key))
	}
	return val
}
