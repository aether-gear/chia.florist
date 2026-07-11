package domain

type ReputationReport struct {
	IP        string         `json:"ip"`
	RawReport map[string]any `json:"raw_report"`
}

type GeoReport struct {
	IP        string         `json:"ip"`
	RawReport map[string]any `json:"raw_report"`
}
