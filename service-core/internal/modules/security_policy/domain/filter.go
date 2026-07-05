package domain

// FilterConfig holds the two classes of payload filters:
//   - Keywords: raw strings whose presence in a request blocks it.
//   - WhitelistedURLs: URL prefixes that bypass WAF rule evaluation.
type FilterConfig struct {
	Keywords        []string
	WhitelistedURLs []string
}
