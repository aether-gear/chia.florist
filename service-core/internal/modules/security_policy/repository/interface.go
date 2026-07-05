package repository

import (
	"context"

	"service-core/internal/modules/security_policy/domain"
	transaction "service-core/internal/shared/transaction"
)

// SecurityPolicyRepository is the persistence
// interface for the security_policy module.
//
// All three data classes:
// WAF rules, IP access control, and keyword/URL filters
// are managed through this single repository
// to keep the bootstrap wiring simple.
type SecurityPolicyRepository interface {
	// --- WAF Rules ---

	GetRules(
		ctx context.Context,
		exec transaction.Executor,
	) ([]domain.WAFRule, error)

	GetRuleByID(
		ctx context.Context,
		exec transaction.Executor,
		id string,
	) (*domain.WAFRule, error)

	SaveRule(
		ctx context.Context,
		exec transaction.Executor,
		rule domain.WAFRule,
	) error

	UpdateRuleStatus(
		ctx context.Context,
		exec transaction.Executor,
		id string,
		enabled bool,
	) error

	DeleteRule(
		ctx context.Context,
		exec transaction.Executor,
		id string,
	) error

	// --- IP Access Control ---

	GetIPConfig(
		ctx context.Context,
		exec transaction.Executor,
	) (*domain.IPConfig, error)

	// UpsertIPRecord inserts or replaces
	// the IP's record atomically.
	UpsertIPRecord(
		ctx context.Context,
		exec transaction.Executor,
		record domain.IPRecord,
	) error

	DeleteIPRecord(
		ctx context.Context,
		exec transaction.Executor,
		ip string,
	) error

	// --- Keyword / URL Filters ---

	GetFilters(
		ctx context.Context,
		exec transaction.Executor,
	) (*domain.FilterConfig, error)

	// UpsertFilterEntry adds a keyword or
	// URL whitelist entry (idempotent).
	//
	// entryType must be "keyword" or "url".
	UpsertFilterEntry(
		ctx context.Context,
		exec transaction.Executor,
		entryType string,
		value string,
	) error

	// DeleteFilterEntry removes a filter entry.
	//
	// entryType must be "keyword" or "url".
	DeleteFilterEntry(
		ctx context.Context,
		exec transaction.Executor,
		entryType string,
		value string,
	) error
}
