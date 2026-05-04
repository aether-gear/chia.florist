package config

import (
	"fmt"
	"os"
)

func GetEnv(key string, defaultKey ...string) string {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		if defaultKey != nil {
			return defaultKey[0]
		}

		panic(fmt.Sprintf("Missing required environment variable: %s", key))
	}
	return val
}
