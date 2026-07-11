package config

type WAFConfig struct {
	// AutoBanEnabled controls whether the WAF automatically bans an IP
	// after a rule or keyword match. Default: true.
	// Set WAF_AUTO_BAN_ENABLED=false to disable auto-ban (e.g. in staging).
	AutoBanEnabled    bool
	VirusTotalAPIKey  string
	IP2LocationAPIKey string
}

func LoadWAFConfig() WAFConfig {
	autoBan := GetEnv("WAF_AUTO_BAN_ENABLED", "true")
	vtKey := GetEnv("VIRUSTOTAL_API_KEY", "")
	ip2locKey := GetEnv("IP2LOCATION_API_KEY", "")

	return WAFConfig{
		AutoBanEnabled:    autoBan != "false",
		VirusTotalAPIKey:  vtKey,
		IP2LocationAPIKey: ip2locKey,
	}
}
