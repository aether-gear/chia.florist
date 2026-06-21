package config

import (
	"log"
	"strconv"
)

type MidTransConfig struct {
	IsProduction bool
	URL          string
	MerchantID   string
	ClientKey    string
	ServerKey    string
}

func LoadMidTransConfig() MidTransConfig {
	isProduction, err := strconv.
		ParseBool(GetEnv("MIDTRANS_IS_PRODUCTION"))
	if err != nil {
		log.Fatalf("failed to parse: %v", err)
	}

	return MidTransConfig{
		IsProduction: isProduction,
		URL:          GetEnv("MIDTRANS_URL"),
		MerchantID:   GetEnv("MIDTRANS_MERCHANT_ID"),
		ClientKey:    GetEnv("MIDTRANS_CLIENT_KEY"),
		ServerKey:    GetEnv("MIDTRANS_SERVER_KEY"),
	}
}
