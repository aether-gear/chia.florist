package repository

import (
	"context"

	"service-core/internal/modules/courier/domain"

	"github.com/google/uuid"
)

type CourierRepository interface {
	GetActiveCodes(
		ctx context.Context,
		codes []string,
	) ([]string, error)

	ValidateCouriers(
		ctx context.Context,
		codes []string,
	) ([]string, error)
}

type ShopCourierRepository interface {
	GetByShopID(
		ctx context.Context,
		shopID uuid.UUID,
	) ([]domain.ShopCourier, error)

	SaveShopCouriers(
		ctx context.Context,
		shopCouriers []domain.ShopCourier,
	) error
}
