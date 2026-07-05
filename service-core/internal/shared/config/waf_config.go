package config

type WAFConfig struct {
	// AutoBanEnabled controls whether the WAF automatically bans an IP
	// after a rule or keyword match. Default: true.
	// Set WAF_AUTO_BAN_ENABLED=false to disable auto-ban (e.g. in staging).
	AutoBanEnabled bool
}

func LoadWAFConfig() WAFConfig {
	autoBan := GetEnv("WAF_AUTO_BAN_ENABLED", "true")
	return WAFConfig{
		AutoBanEnabled: autoBan != "false",
	}
}
