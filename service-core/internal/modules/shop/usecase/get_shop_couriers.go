package usecase

import (
	"context"
	"fmt"

	courierDomain "service-core/internal/modules/courier/domain"
	courierRepo "service-core/internal/modules/courier/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type GetShopCouriersUsecase struct {
	courierRepo courierRepo.ShopCourierRepository
	executor    transaction.Executor
}

func NewGetShopCouriersUsecase(
	courierRepo courierRepo.ShopCourierRepository,
	executor transaction.Executor,
) *GetShopCouriersUsecase {
	return &GetShopCouriersUsecase{
		courierRepo: courierRepo,
		executor:    executor,
	}
}

func (u *GetShopCouriersUsecase) Execute(
	ctx context.Context,
	shopID uuid.UUID,
) ([]courierDomain.ShopCourier, error) {
	couriers, err := u.courierRepo.ListByShopID(ctx, u.executor, shopID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve shop couriers: %w", err)
	}

	return couriers, nil
}
