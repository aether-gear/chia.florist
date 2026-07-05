package persistence

import "service-core/internal/modules/security_policy/domain"

func ruleRowToDomain(row wafRuleRow) domain.WAFRule {
	tags := row.Tags
	if tags == nil {
		tags = []string{}
	}
	return domain.WAFRule{
		ID:          row.ID,
		Description: row.Description,
		Pattern:     row.Pattern,
		Tags:        tags,
		Impact:      row.Impact,
		Enabled:     row.Enabled,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func ipRowToDomain(row ipAccessControlRow) domain.IPRecord {
	return domain.IPRecord{
		IP:     row.IP,
		Status: domain.IPStatus(row.Status),
		Reason: row.Reason,
	}
}
