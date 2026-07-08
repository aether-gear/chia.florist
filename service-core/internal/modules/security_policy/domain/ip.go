package domain

// IPStatus represents the access
// control status of an IP address.
type IPStatus string

const (
	IPStatusBanned      IPStatus = "banned"
	IPStatusWhitelisted IPStatus = "whitelisted"
	IPStatusIgnored     IPStatus = "ignored"
	IPStatusBannedMuted      IPStatus = "banned_muted"
	IPStatusWhitelistedMuted IPStatus = "whitelisted_muted"
)

// IPRecord associates a single IP address
// with an access control status.
type IPRecord struct {
	IP     string
	Status IPStatus
	Reason string
}

// IPConfig is the aggregate view of
// all IP access control entries.
type IPConfig struct {
	Records []IPRecord
}
