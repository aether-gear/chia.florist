package appmiddleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	apphttp "service-core/internal/common/http"
	applogger "service-core/internal/common/logger"
	secPolicyUsecase "service-core/internal/modules/security_policy/usecase"
)

var (
	requestMu  sync.Mutex
	ipRequests = make(map[string][]time.Time)
)

const (
	rateLimitWindow = 10 * time.Second
	rateLimitMax    = 30 // Max 30 requests per 10 seconds
)

func isRateLimited(ip string) bool {
	requestMu.Lock()
	defer requestMu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rateLimitWindow)

	// Clean up old timestamps
	times := ipRequests[ip]
	var validTimes []time.Time
	for _, t := range times {
		if t.After(cutoff) {
			validTimes = append(validTimes, t)
		}
	}

	if len(validTimes) >= rateLimitMax {
		ipRequests[ip] = validTimes
		return true
	}

	validTimes = append(validTimes, now)
	ipRequests[ip] = validTimes
	return false
}

// WAF returns a middleware that inspects every incoming
// request against the configured WAF rules,
// IP access control list, and keyword filters.
//
// Requests to the /api/ prefix are bypassed (admin panel endpoints)
// to avoid the WAF evaluating its own management traffic
// — mirroring the original prototype.
func WAF(
	inspectPayload *secPolicyUsecase.InspectPayloadUsecase,
	updateIPAction *secPolicyUsecase.UpdateIPActionUsecase,
	auditLogger applogger.AuditLogger,
	autoBanEnabled bool,
) Middleware {
	return func(next apphttp.AppHandler) apphttp.AppHandler {
		return func(w http.ResponseWriter, r *http.Request) error {
			if strings.HasPrefix(r.URL.Path, "/api/") {
				return next(w, r)
			}

			ip := apphttp.ClientIP(r)

			if isRateLimited(ip) {
				log.Printf("⚠️  [WAF Rate Limit Blocked] IP: %s | Path: %s\n", ip, r.URL.Path)
				if autoBanEnabled && !isLocalhost(ip) {
					_ = updateIPAction.Execute(
						r.Context(),
						secPolicyUsecase.UpdateIPActionInput{
							IP:     ip,
							Action: "ban",
							Reason: "Auto-Banned: Rate limit exceeded (Spam)",
						},
					)
				}

				auditLogger.Log(r.Context(), applogger.AuditEvent{
					Category: "waf_event",
					Action:   "rate_limit_exceeded",
					Resource: "request",
					Outcome:  applogger.OutcomeBlocked,
					Metadata: map[string]any{
						"path":   r.URL.Path,
						"ip":     ip,
						"reason": "Rate limit exceeded (max 30 reqs/10s)",
					},
				})

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = w.Write([]byte(`{"error":"429 Too Many Requests - WAF Rate Limited"}`))

				return nil
			}

			// Read and restore the request body
			// so downstream handlers can still use it.
			var bodyBytes []byte
			if r.Body != nil && r.Body != http.NoBody {
				var err error
				bodyBytes, err = io.ReadAll(r.Body)
				if err != nil {
					bodyBytes = []byte{}
				}

				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}

			// Collect only the
			// three security-relevant headers.
			headers := map[string]string{
				"User-Agent": r.Header.Get("User-Agent"),
				"Referer":    r.Header.Get("Referer"),
				"Cookie":     r.Header.Get("Cookie"),
			}

			result, err := inspectPayload.
				Execute(
					r.Context(),
					secPolicyUsecase.InspectPayloadInput{
						ClientIP: ip,
						Path:     r.URL.Path,
						RawQuery: r.URL.RawQuery,
						Headers:  headers,
						Body:     bodyBytes,
					},
				)
			if err != nil {
				// On an inspection error, system will allow the request
				// to pass through rather than blocking legitimate traffic
				// due to a DB hiccup. The error is visible via the audit log.
				auditLogger.Log(
					r.Context(),
					applogger.AuditEvent{
						Category: "waf_event",
						Action:   "inspect_error",
						Resource: "request",
						Outcome:  applogger.OutcomeFailure,
						Metadata: map[string]any{
							"path":  r.URL.Path,
							"ip":    ip,
							"error": err.Error(),
						},
					},
				)

				return next(w, r)
			}

			if result.Blocked {
				log.Printf("🛡️  [WAF Payload Blocked] IP: %s | Rule: %s | Reason: %s | Matched: %s\n", ip, result.RuleID, result.Reason, result.MatchedPayload)
				if autoBanEnabled &&
					!isLocalhost(ip) {
					_ = updateIPAction.Execute(
						r.Context(),
						secPolicyUsecase.UpdateIPActionInput{
							IP:     ip,
							Action: "ban",
							Reason: "Auto-Banned: " + result.Reason,
						},
					)
				}

				auditLogger.Log(r.Context(), applogger.AuditEvent{
					Category: "waf_event",
					Action:   "request_blocked",
					Resource: "request",
					Outcome:  applogger.OutcomeBlocked,
					Metadata: map[string]any{
						"path":       r.URL.RequestURI(),
						"ip":         ip,
						"reason":     result.Reason,
						"rule_id":    result.RuleID,
						"method":     r.Method,
						"payload":    result.MatchedPayload,
						"user_agent": r.UserAgent(),
					},
				})

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				_, _ = w.Write([]byte(`{"error":"403 Forbidden - WAF Blocked Request"}`))

				return nil
			}

			auditLogger.Log(
				r.Context(),
				applogger.AuditEvent{
					Category: "waf_event",
					Action:   "request_allowed",
					Resource: "request",
					Outcome:  applogger.OutcomeSuccess,
					Metadata: map[string]any{
						"path":       r.URL.RequestURI(),
						"ip":         ip,
						"method":     r.Method,
						"user_agent": r.UserAgent(),
					},
				},
			)

			return next(w, r)
		}
	}
}

func isLocalhost(ip string) bool {
	switch ip {
	case "127.0.0.1", "::1", "localhost":
		return true
	}

	return false
}
