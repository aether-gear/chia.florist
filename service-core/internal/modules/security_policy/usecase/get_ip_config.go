package usecase

import (
	"context"
	"fmt"

	"service-core/internal/modules/security_policy/domain"
	"service-core/internal/modules/security_policy/repository"
	transaction "service-core/internal/shared/transaction"
)

type GetIPConfigUsecase struct {
	executor     transaction.Executor
	securityRepo repository.SecurityPolicyRepository
}

func NewGetIPConfigUsecase(
	executor transaction.Executor,
	securityRepo repository.SecurityPolicyRepository,
) *GetIPConfigUsecase {
	return &GetIPConfigUsecase{
		executor:     executor,
		securityRepo: securityRepo,
	}
}

func (u *GetIPConfigUsecase) Execute(
	ctx context.Context,
) (*domain.IPConfig, error) {
	config, err := u.securityRepo.GetIPConfig(ctx, u.executor)
	if err != nil {
		return nil, fmt.Errorf("get ip config: %w", err)
	}

	return config, nil
}
