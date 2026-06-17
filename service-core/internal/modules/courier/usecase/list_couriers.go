package usecase

import (
	"context"
	"fmt"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/courier/repository"
	transaction "service-core/internal/shared/transaction"
)

type ListCouriersUsecase struct {
	exec        transaction.Executor
	courierRepo repository.CourierRepository
}

func NewListCouriersUsecase(
	exec transaction.Executor,
	courierRepo repository.CourierRepository,
) *ListCouriersUsecase {
	return &ListCouriersUsecase{
		exec:        exec,
		courierRepo: courierRepo,
	}
}

func (u *ListCouriersUsecase) Execute(
	ctx context.Context,
) ([]string, error) {
	codes, err := u.courierRepo.ListAll(ctx, u.exec)
	if err != nil {
		return nil, fmt.Errorf("failed to load couriers: %w", err)
	}
	if len(codes) == 0 {
		return nil, apperrors.NewNotFound("no courier service available at the moment")
	}

	return codes, nil
}
