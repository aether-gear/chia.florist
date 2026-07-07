package http

import "time"

// --- WAF Rules ---

type ruleResponse struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Pattern     string    `json:"pattern"`
	Tags        []string  `json:"tags"`
	Impact      string    `json:"impact,omitempty"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type createRuleRequest struct {
	Description string   `json:"description"`
	Pattern     string   `json:"pattern"`
	Tags        []string `json:"tags,omitempty"`
	Impact      string   `json:"impact,omitempty"`
}

type toggleRuleRequest struct {
	Enabled bool `json:"enabled"`
}

type updateRuleRequest struct {
	Description *string  `json:"description,omitempty"`
	Pattern     *string  `json:"pattern,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Impact      *string  `json:"impact,omitempty"`
	Enabled     *bool    `json:"enabled,omitempty"`
}

// --- IP Access Control ---

type ipEntryResponse struct {
	IP     string `json:"ip"`
	Status string `json:"status"` // "banned" | "whitelisted" | "ignored"
	Reason string `json:"reason,omitempty"`
}

type updateIPActionRequest struct {
	IP string `json:"ip"`
	// Action must be one of: "ban", "whitelist", "ignore", "reset".
	Action string `json:"action"`
	Reason string `json:"reason,omitempty"`
}

// --- Filters ---

type filterConfigResponse struct {
	Keywords        []string `json:"keywords"`
	WhitelistedURLs []string `json:"whitelisted_urls"`
}

type updateFilterRequest struct {
	// Type must be "keyword" or "url".
	Type string `json:"type"`
	// Action must be "add" or "remove".
	Action string `json:"action"`
	Value  string `json:"value"`
}
