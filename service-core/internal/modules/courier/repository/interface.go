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
	ListByShopID(
		ctx context.Context,
		exec transaction.Executor,
		shopID uuid.UUID,
	) ([]domain.ShopCourier, error)

	GetByShopIDAndCode(
		ctx context.Context,
		exec transaction.Executor,
		shopID uuid.UUID,
		code string,
	) (*domain.ShopCourier, error)

	ListsByShopIDs(
		ctx context.Context,
		exec transaction.Executor,
		shopIDs []uuid.UUID,
	) (map[uuid.UUID][]domain.ShopCourier, error)

	SaveShopCourier(
		ctx context.Context,
		exec transaction.Executor,
		shopCourier domain.ShopCourier,
	) error

	VerifyShopCourier(
		ctx context.Context,
		exec transaction.Executor,
		shopID uuid.UUID,
		code string,
		status domain.CourierVerificationStatus,
		active bool,
		verifiedBy uuid.UUID,
		rejectionReason *string,
	) error
}
