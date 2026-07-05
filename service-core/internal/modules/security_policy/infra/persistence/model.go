package persistence

import "time"

type wafRuleRow struct {
	ID          string
	Description string
	Pattern     string
	Tags        []string
	Impact      string
	Enabled     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ipAccessControlRow struct {
	IP        string
	Status    string // "banned" | "whitelisted" | "ignored"
	Reason    string
	UpdatedAt time.Time
}

type filterConfigRow struct {
	Value string
	Type  string // "keyword" | "url"
}
