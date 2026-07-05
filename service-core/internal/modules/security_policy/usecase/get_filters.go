package usecase

import (
	"context"
	"fmt"

	"service-core/internal/modules/security_policy/domain"
	"service-core/internal/modules/security_policy/repository"
	transaction "service-core/internal/shared/transaction"
)

type GetFiltersUsecase struct {
	executor     transaction.Executor
	securityRepo repository.SecurityPolicyRepository
}

func NewGetFiltersUsecase(
	executor transaction.Executor,
	securityRepo repository.SecurityPolicyRepository,
) *GetFiltersUsecase {
	return &GetFiltersUsecase{
		executor:     executor,
		securityRepo: securityRepo,
	}
}

func (u *GetFiltersUsecase) Execute(
	ctx context.Context,
) (*domain.FilterConfig, error) {
	config, err := u.securityRepo.GetFilters(ctx, u.executor)
	if err != nil {
		return nil, fmt.Errorf("get filters: %w", err)
	}

	return config, nil
}
