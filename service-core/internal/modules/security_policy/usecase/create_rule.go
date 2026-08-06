package usecase

import (
	"context"
	"fmt"
	"regexp"

	appclock "service-core/internal/common/clock"
	"service-core/internal/modules/security_policy/domain"
	"service-core/internal/modules/security_policy/repository"
	transaction "service-core/internal/shared/transaction"
)

type CreateRuleUsecase struct {
	executor     transaction.Executor
	securityRepo repository.SecurityPolicyRepository
}

func NewCreateRuleUsecase(
	executor transaction.Executor,
	securityRepo repository.SecurityPolicyRepository,
) *CreateRuleUsecase {
	return &CreateRuleUsecase{
		executor:     executor,
		securityRepo: securityRepo,
	}
}

type CreateRuleInput struct {
	Description string
	Pattern     string
	Tags        []string
	Impact      string
}

func (u *CreateRuleUsecase) Execute(
	ctx context.Context,
	input CreateRuleInput,
) (*domain.WAFRule, error) {
	if _, err := regexp.Compile(input.Pattern); err != nil {
		return nil, domain.ErrInvalidPattern
	}

	rules, err := u.securityRepo.GetRules(ctx, u.executor)
	if err != nil {
		return nil, fmt.Errorf("create waf rule: load existing rules: %w", err)
	}

	tags := input.Tags
	if tags == nil {
		tags = []string{}
	}

	now := appclock.Now()
	rule := domain.WAFRule{
		ID:          fmt.Sprintf("%d", len(rules)+1000),
		Description: input.Description,
		Pattern:     input.Pattern,
		Tags:        tags,
		Impact:      input.Impact,
		Enabled:     true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := u.securityRepo.
		SaveRule(ctx, u.executor, rule); err != nil {
		return nil, fmt.Errorf("create waf rule: %w", err)
	}

	return &rule, nil
}
