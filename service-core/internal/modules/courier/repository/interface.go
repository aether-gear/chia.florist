package repository

import (
	"context"

	"service-core/internal/modules/courier/domain"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type CourierRepository interface {
	ListAll(
		ctx context.Context,
		exec transaction.Executor,
	) ([]string, error)

	GetActiveCodes(
		ctx context.Context,
		exec transaction.Executor,
		codes []string,
	) ([]string, error)

	ValidateCouriers(
		ctx context.Context,
		exec transaction.Executor,
		codes []string,
	) ([]string, error)
}

type ShopCourierRepository interface {
	GetByShopID(
		ctx context.Context,
		exec transaction.Executor,
		shopID uuid.UUID,
	) ([]domain.ShopCourier, error)

	SaveShopCouriers(
		ctx context.Context,
		exec transaction.Executor,
		shopCouriers []domain.ShopCourier,
	) error
}
