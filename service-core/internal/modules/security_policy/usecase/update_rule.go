package usecase

import (
	"context"
	"fmt"
	"regexp"

	"service-core/internal/modules/security_policy/domain"
	"service-core/internal/modules/security_policy/repository"
	transaction "service-core/internal/shared/transaction"
)

type UpdateRuleUsecase struct {
	executor     transaction.Executor
	securityRepo repository.SecurityPolicyRepository
}

func NewUpdateRuleUsecase(
	executor transaction.Executor,
	securityRepo repository.SecurityPolicyRepository,
) *UpdateRuleUsecase {
	return &UpdateRuleUsecase{
		executor:     executor,
		securityRepo: securityRepo,
	}
}

type UpdateRuleInput struct {
	ID          string
	Description *string
	Pattern     *string
	Tags        []string
	Impact      *string
	Enabled     *bool
}

func (u *UpdateRuleUsecase) Execute(
	ctx context.Context,
	input UpdateRuleInput,
) (*domain.WAFRule, error) {
	rules, err := u.securityRepo.GetRules(ctx, u.executor)
	if err != nil {
		return nil, fmt.Errorf("update waf rule: load rules: %w", err)
	}

	var existing *domain.WAFRule
	for _, r := range rules {
		if r.ID == input.ID {
			existing = &r
			break
		}
	}

	if existing == nil {
		return nil, fmt.Errorf("update waf rule: rule not found")
	}

	if input.Description != nil {
		existing.Description = *input.Description
	}
	if input.Pattern != nil {
		if _, err := regexp.Compile(*input.Pattern); err != nil {
			return nil, domain.ErrInvalidPattern
		}
		existing.Pattern = *input.Pattern
	}
	if input.Tags != nil {
		existing.Tags = input.Tags
	}
	if input.Impact != nil {
		existing.Impact = *input.Impact
	}
	if input.Enabled != nil {
		existing.Enabled = *input.Enabled
	}

	if err := u.securityRepo.SaveRule(ctx, u.executor, *existing); err != nil {
		return nil, fmt.Errorf("update waf rule: save failed: %w", err)
	}

	return existing, nil
}
