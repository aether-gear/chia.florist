package usecase

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"service-core/internal/modules/security_policy/domain"
	"service-core/internal/modules/security_policy/repository"
	transaction "service-core/internal/shared/transaction"
)

// commentRegexp normalises SQL block comments
// (e.g. UNION/**/SELECT → UNION SELECT)
// to prevent comment-obfuscation evasion techniques.
var commentRegexp = regexp.MustCompile(`(?s)/\*.*?\*/`)

type InspectPayloadUsecase struct {
	executor     transaction.Executor
	securityRepo repository.SecurityPolicyRepository
}

func NewInspectPayloadUsecase(
	executor transaction.Executor,
	securityRepo repository.SecurityPolicyRepository,
) *InspectPayloadUsecase {
	return &InspectPayloadUsecase{
		executor:     executor,
		securityRepo: securityRepo,
	}
}

type InspectPayloadInput struct {
	ClientIP string
	Path     string
	RawQuery string
	Headers  map[string]string
	Body     []byte
}

type InspectionResult struct {
	Blocked        bool
	Reason         string
	RuleID         string
	MatchedPayload string
	Silent         bool
}

func (u *InspectPayloadUsecase) Execute(
	ctx context.Context,
	input InspectPayloadInput,
) (*InspectionResult, error) {
	ipConfig, err := u.securityRepo.GetIPConfig(ctx, u.executor)
	if err != nil {
		return nil, fmt.Errorf("inspect payload: load ip config: %w", err)
	}

	// Check whether the client's IP address has an
	// explicit access policy before performing any
	// payload inspection.
	for _, rec := range ipConfig.Records {
		if rec.IP != input.ClientIP {
			continue
		}
		switch rec.Status {
		case domain.IPStatusBanned:
			reason := "IP is manually banned"
			if rec.Reason != "" {
				reason += " (" + rec.Reason + ")"
			}
			return &InspectionResult{
				Blocked: true,
				Reason:  reason,
				RuleID:  "IP-BAN",
			}, nil
		case domain.IPStatusBannedMuted:
			reason := "IP is manually banned (Muted)"
			if rec.Reason != "" {
				reason += " (" + rec.Reason + ")"
			}
			return &InspectionResult{
				Blocked: true,
				Reason:  reason,
				RuleID:  "IP-BAN-MUTED",
				Silent:  true,
			}, nil
		case domain.IPStatusWhitelisted:
			return &InspectionResult{
				Blocked: false,
			}, nil
		case domain.IPStatusWhitelistedMuted:
			return &InspectionResult{
				Blocked: false,
				Silent:  true,
			}, nil
		case domain.IPStatusIgnored:
			return &InspectionResult{
				Blocked: false,
				Silent:  true,
			}, nil
		}
		break
	}

	filters, err := u.securityRepo.GetFilters(ctx, u.executor)
	if err != nil {
		return nil, fmt.Errorf("inspect payload: load filters: %w", err)
	}

	// Skip payload inspection for requests targeting
	// trusted endpoints.
	//
	// This allows health checks, internal callbacks,
	// or other approved routes to bypass WAF evaluation
	// and avoids unnecessary processing.
	lowerPath := strings.ToLower(input.Path)
	for _, wURL := range filters.WhitelistedURLs {
		if strings.HasPrefix(lowerPath, strings.ToLower(wURL)) {
			return &InspectionResult{
				Blocked: false,
			}, nil
		}
	}

	// Normalize the query string by repeatedly URL-decoding
	// to uncover nested encodings, then decode IIS-style Unicode
	// sequences and strip null bytes commonly used in evasion attempts.
	decodedQuery := input.RawQuery
	for range 3 {
		unescaped, err := url.QueryUnescape(decodedQuery)
		if err != nil || unescaped == decodedQuery {
			break
		}
		decodedQuery = unescaped
	}
	decodedQuery = decodeIISUnicode(decodedQuery)
	decodedQuery = strings.ReplaceAll(decodedQuery, "\x00", "")

	// Normalize the request path using the same decoding
	// strategy as the query string to ensure encoded attack
	// patterns are inspected.
	decodedPath := input.Path
	for range 3 {
		unescaped, err := url.PathUnescape(decodedPath)
		if err != nil || unescaped == decodedPath {
			break
		}
		decodedPath = unescaped
	}
	decodedPath = decodeIISUnicode(decodedPath)
	decodedPath = strings.ReplaceAll(decodedPath, "\x00", "")

	// Build a canonical header representation using only
	// headers that commonly carry user-controlled input.
	var headersBuf bytes.Buffer
	for _, name := range []string{"User-Agent", "Referer", "Cookie"} {
		if value, ok := input.Headers[name]; ok && value != "" {
			headersBuf.WriteString(name)
			headersBuf.WriteString(": ")
			headersBuf.WriteString(value)
			headersBuf.WriteByte('\n')
		}
	}
	headersStr := headersBuf.String()

	payload := fmt.Sprintf("%s %s\n%s\n%s",
		decodedPath,
		decodedQuery,
		headersStr,
		string(input.Body),
	)
	payload = commentRegexp.ReplaceAllString(payload, " ")
	lowerPayload := strings.ToLower(payload)

	// Compare the normalized request against configured
	// keyword signatures.
	//
	// Keyword matching provides a lightweight first-pass
	// inspection for well-known attack patterns before
	// executing more expensive regular expression rules.
	for _, kw := range filters.Keywords {
		if strings.Contains(lowerPayload, strings.ToLower(kw)) {
			return &InspectionResult{
				Blocked:        true,
				Reason:         "Blocked by keyword: " + kw,
				RuleID:         "KW-BLOCK",
				MatchedPayload: kw,
			}, nil
		}
	}

	rules, err := u.securityRepo.GetRules(ctx, u.executor)
	if err != nil {
		return nil, fmt.Errorf("inspect payload: load waf rules: %w", err)
	}

	// Evaluate enabled regular expression rules against
	// the normalized request.
	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}

		matched, _ := regexp.MatchString(rule.Pattern, lowerPayload)
		if matched {
			var matchedStr string
			re, err := regexp.Compile(rule.Pattern)
			if err == nil {
				matchedStr = re.FindString(lowerPayload)
			}
			if matchedStr == "" {
				matchedStr = lowerPayload
			}
			if len(matchedStr) > 500 {
				matchedStr = matchedStr[:500] + "..."
			}

			return &InspectionResult{
				Blocked:        true,
				Reason:         "Matched rule: " + rule.Description,
				RuleID:         rule.ID,
				MatchedPayload: matchedStr,
			}, nil
		}
	}

	return &InspectionResult{
		Blocked: false,
	}, nil
}

// decodeIISUnicode decodes IIS-style %uXXXX unicode escape sequences
// used to bypass pattern matching via alternate encoding.
func decodeIISUnicode(s string) string {
	for {
		idx := strings.Index(s, "%u")
		if idx == -1 || idx+6 > len(s) {
			break
		}

		hexStr := s[idx+2 : idx+6]
		val, err := strconv.ParseUint(hexStr, 16, 16)
		if err != nil {
			s = s[:idx] + "%%u" + s[idx+2:]
			continue
		}

		s = s[:idx] + string(rune(val)) + s[idx+6:]
	}

	return s
}
