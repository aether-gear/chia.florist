package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"service-core/internal/modules/security_policy/domain"
	"service-core/internal/modules/security_policy/repository"
	"service-core/internal/shared/transaction"

	"github.com/jackc/pgx/v5"
)

type securityPolicyRepositoryImpl struct{}

func NewSecurityPolicyRepository() repository.SecurityPolicyRepository {
	return &securityPolicyRepositoryImpl{}
}

// --- WAF Rules ---

func (r *securityPolicyRepositoryImpl) GetRules(
	ctx context.Context,
	exec transaction.Executor,
) ([]domain.WAFRule, error) {
	queryStr := `
		SELECT
			id,
			description,
			pattern,
			tags,
			impact,
			enabled,
			created_at,
			updated_at
		FROM waf_rules
		ORDER BY created_at DESC
	`

	rows, err := exec.Query(ctx, queryStr)
	if err != nil {
		return nil, fmt.Errorf("security_policy: query waf rules failed: %w", err)
	}
	defer rows.Close()

	var results []domain.WAFRule
	for rows.Next() {
		var row wafRuleRow
		if err := rows.Scan(
			&row.ID,
			&row.Description,
			&row.Pattern,
			&row.Tags,
			&row.Impact,
			&row.Enabled,
			&row.CreatedAt,
			&row.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("security_policy: scan waf rule failed: %w", err)
		}
		results = append(results, ruleRowToDomain(row))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("security_policy: iterate waf rules failed: %w", err)
	}

	return results, nil
}

func (r *securityPolicyRepositoryImpl) GetRuleByID(
	ctx context.Context,
	exec transaction.Executor,
	id string,
) (*domain.WAFRule, error) {
	queryStr := `
		SELECT
			id,
			description,
			pattern,
			tags,
			impact,
			enabled,
			created_at,
			updated_at
		FROM waf_rules
		WHERE id = $1
		LIMIT 1
	`

	var row wafRuleRow
	err := exec.QueryRow(ctx, queryStr, id).Scan(
		&row.ID,
		&row.Description,
		&row.Pattern,
		&row.Tags,
		&row.Impact,
		&row.Enabled,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("security_policy: query waf rule by id failed: %w", err)
	}

	result := ruleRowToDomain(row)
	return &result, nil
}

func (r *securityPolicyRepositoryImpl) SaveRule(
	ctx context.Context,
	exec transaction.Executor,
	rule domain.WAFRule,
) error {
	now := time.Now().UTC()
	if rule.CreatedAt.IsZero() {
		rule.CreatedAt = now
	}
	rule.UpdatedAt = now

	tags := rule.Tags
	if tags == nil {
		tags = []string{}
	}

	queryStr := `
		INSERT INTO waf_rules (
			id,
			description,
			pattern,
			tags,
			impact,
			enabled,
			created_at,
			updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		ON CONFLICT (id) DO UPDATE SET
			description = EXCLUDED.description,
			pattern = EXCLUDED.pattern,
			tags = EXCLUDED.tags,
			impact = EXCLUDED.impact,
			enabled = EXCLUDED.enabled,
			updated_at = EXCLUDED.updated_at
	`

	_, err := exec.Exec(ctx, queryStr,
		rule.ID,
		rule.Description,
		rule.Pattern,
		tags,
		rule.Impact,
		rule.Enabled,
		rule.CreatedAt,
		rule.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("security_policy: save waf rule failed: %w", err)
	}

	return nil
}

func (r *securityPolicyRepositoryImpl) UpdateRuleStatus(
	ctx context.Context,
	exec transaction.Executor,
	id string,
	enabled bool,
) error {
	queryStr := `
		UPDATE waf_rules
		SET
			enabled = $1,
			updated_at = $2
		WHERE
			id = $3
	`

	tag, err := exec.Exec(ctx, queryStr, enabled, time.Now().UTC(), id)
	if err != nil {
		return fmt.Errorf("security_policy: update waf rule status failed: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrRuleNotFound
	}

	return nil
}

func (r *securityPolicyRepositoryImpl) DeleteRule(
	ctx context.Context,
	exec transaction.Executor,
	id string,
) error {
	queryStr := `
		DELETE FROM waf_rules
		WHERE id = $1
		`

	tag, err := exec.Exec(ctx, queryStr, id)
	if err != nil {
		return fmt.Errorf("security_policy: delete waf rule failed: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrRuleNotFound
	}

	return nil
}

// --- IP Access Control ---

func (r *securityPolicyRepositoryImpl) GetIPConfig(
	ctx context.Context,
	exec transaction.Executor,
) (*domain.IPConfig, error) {
	queryStr := `
		SELECT
			ip,
			status,
			reason
		FROM ip_access_control
		ORDER BY ip ASC
	`

	rows, err := exec.Query(ctx, queryStr)
	if err != nil {
		return nil, fmt.Errorf("security_policy: query ip config failed: %w", err)
	}
	defer rows.Close()

	config := &domain.IPConfig{Records: []domain.IPRecord{}}
	for rows.Next() {
		var row ipAccessControlRow
		if err := rows.Scan(&row.IP, &row.Status, &row.Reason); err != nil {
			return nil, fmt.Errorf("security_policy: scan ip record failed: %w", err)
		}
		config.Records = append(config.Records, ipRowToDomain(row))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("security_policy: iterate ip records failed: %w", err)
	}

	return config, nil
}

func (r *securityPolicyRepositoryImpl) UpsertIPRecord(
	ctx context.Context,
	exec transaction.Executor,
	record domain.IPRecord,
) error {
	queryStr := `
		INSERT INTO ip_access_control (
			ip,
			status,
			reason,
			updated_at
		)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (ip) DO UPDATE
		SET
			status = EXCLUDED.status,
		    reason = EXCLUDED.reason,
		    updated_at = EXCLUDED.updated_at
	`

	_, err := exec.Exec(ctx, queryStr,
		record.IP,
		string(record.Status),
		record.Reason,
		time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("security_policy: upsert ip record failed: %w", err)
	}

	return nil
}

func (r *securityPolicyRepositoryImpl) DeleteIPRecord(
	ctx context.Context,
	exec transaction.Executor,
	ip string,
) error {
	queryStr := `
		DELETE FROM ip_access_control
		WHERE ip = $1
	`

	_, err := exec.Exec(ctx, queryStr, ip)
	if err != nil {
		return fmt.Errorf("security_policy: delete ip record failed: %w", err)
	}

	return nil
}

// --- Keyword / URL Filters ---

func (r *securityPolicyRepositoryImpl) GetFilters(
	ctx context.Context,
	exec transaction.Executor,
) (*domain.FilterConfig, error) {
	queryStr := `
		SELECT
			value,
			type
		FROM filter_config
		ORDER BY type, value`

	rows, err := exec.Query(ctx, queryStr)
	if err != nil {
		return nil, fmt.Errorf("security_policy: query filter config failed: %w", err)
	}
	defer rows.Close()

	config := &domain.FilterConfig{
		Keywords:        []string{},
		WhitelistedURLs: []string{},
	}

	for rows.Next() {
		var row filterConfigRow
		if err := rows.Scan(&row.Value, &row.Type); err != nil {
			return nil, fmt.Errorf("security_policy: scan filter entry failed: %w", err)
		}
		switch row.Type {
		case "keyword":
			config.Keywords = append(config.Keywords, row.Value)
		case "url":
			config.WhitelistedURLs = append(config.WhitelistedURLs, row.Value)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("security_policy: iterate filter entries failed: %w", err)
	}

	return config, nil
}

func (r *securityPolicyRepositoryImpl) UpsertFilterEntry(
	ctx context.Context,
	exec transaction.Executor,
	entryType string,
	value string,
) error {
	queryStr := `
		INSERT INTO filter_config (
			value,
			type
		)
		VALUES ($1,$2)
		ON CONFLICT (value, type) DO NOTHING
	`

	_, err := exec.Exec(ctx, queryStr, value, entryType)
	if err != nil {
		return fmt.Errorf("security_policy: upsert filter entry failed: %w", err)
	}

	return nil
}

func (r *securityPolicyRepositoryImpl) DeleteFilterEntry(
	ctx context.Context,
	exec transaction.Executor,
	entryType string,
	value string,
) error {
	queryStr := `
		DELETE FROM filter_config
		WHERE
			value = $1
			AND type = $2
	`

	_, err := exec.Exec(ctx, queryStr, value, entryType)
	if err != nil {
		return fmt.Errorf("security_policy: delete filter entry failed: %w", err)
	}

	return nil
}
