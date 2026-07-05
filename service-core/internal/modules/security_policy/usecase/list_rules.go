package usecase

import (
	"context"
	"fmt"

	"service-core/internal/modules/security_policy/domain"
	"service-core/internal/modules/security_policy/repository"
	transaction "service-core/internal/shared/transaction"
)

type ListRulesUsecase struct {
	executor     transaction.Executor
	securityRepo repository.SecurityPolicyRepository
}

func NewListRulesUsecase(
	executor transaction.Executor,
	securityRepo repository.SecurityPolicyRepository,
) *ListRulesUsecase {
	return &ListRulesUsecase{
		executor:     executor,
		securityRepo: securityRepo,
	}
}

func (u *ListRulesUsecase) Execute(
	ctx context.Context,
) ([]domain.WAFRule, error) {
	rules, err := u.securityRepo.GetRules(ctx, u.executor)
	if err != nil {
		return nil, fmt.Errorf("list waf rules: %w", err)
	}
	if rules == nil {
		rules = []domain.WAFRule{}
	}

	return rules, nil
}
