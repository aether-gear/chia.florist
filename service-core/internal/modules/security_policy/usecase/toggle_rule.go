package usecase

import (
	"context"
	"fmt"

	"service-core/internal/modules/security_policy/repository"
	transaction "service-core/internal/shared/transaction"
)

type ToggleRuleUsecase struct {
	executor     transaction.Executor
	securityRepo repository.SecurityPolicyRepository
}

func NewToggleRuleUsecase(
	executor transaction.Executor,
	securityRepo repository.SecurityPolicyRepository,
) *ToggleRuleUsecase {
	return &ToggleRuleUsecase{
		executor:     executor,
		securityRepo: securityRepo,
	}
}

type ToggleRuleInput struct {
	ID      string
	Enabled bool
}

func (u *ToggleRuleUsecase) Execute(
	ctx context.Context,
	input ToggleRuleInput,
) error {
	if err := u.securityRepo.
		UpdateRuleStatus(
			ctx,
			u.executor,
			input.ID,
			input.Enabled,
		); err != nil {
		return fmt.Errorf("toggle waf rule: %w", err)
	}

	return nil
}
