package usecase

import (
	"context"
	"fmt"

	"service-core/internal/modules/security_policy/repository"
	transaction "service-core/internal/shared/transaction"
)

type DeleteRuleUsecase struct {
	executor     transaction.Executor
	securityRepo repository.SecurityPolicyRepository
}

func NewDeleteRuleUsecase(
	executor transaction.Executor,
	securityRepo repository.SecurityPolicyRepository,
) *DeleteRuleUsecase {
	return &DeleteRuleUsecase{
		executor:     executor,
		securityRepo: securityRepo,
	}
}

func (u *DeleteRuleUsecase) Execute(
	ctx context.Context,
	id string,
) error {
	if err := u.securityRepo.
		DeleteRule(ctx, u.executor, id); err != nil {
		return fmt.Errorf("delete waf rule: %w", err)
	}

	return nil
}
