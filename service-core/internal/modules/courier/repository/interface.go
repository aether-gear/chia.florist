package repository

import (
	"service-core/internal/modules/courier/domain"

	"github.com/google/uuid"
)

type CourierRepository interface {
	GetActiveCodes(codes []string) ([]string, error)
	ValidateCouriers(codes []string) ([]string, error)
}

type ShopCourierRepository interface {
	GetByShopID(shopID uuid.UUID) ([]domain.ShopCourier, error)
	SaveShopCouriers(shopCouriers []domain.ShopCourier) error
}
