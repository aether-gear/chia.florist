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

	// ListsByShopIDs retrieves couriers grouped by shop IDs.
	// The returned map uses the shop ID as the key and all associated
	// courier records as the value.
	//
	// Example:
	//
	//	shopIDs := []uuid.UUID{shopA, shopB}
	//
	//	result := map[uuid.UUID][]domain.Courier{
	//		shopA: {
	//			courierA1,
	//			courierA2,
	//		},
	//		shopB: {
	//			courierB1,
	//		},
	//	}
	//
	// This allows callers to efficiently look up couriers belonging to a
	// specific shop without additional filtering.
	ListsByShopIDs(
		ctx context.Context,
		exec transaction.Executor,
		shopIDs []uuid.UUID,
	) (map[uuid.UUID][]domain.ShopCourier, error)

	SaveShopCourier(
		ctx context.Context,
		exec transaction.Executor,
		shopCouriers domain.ShopCourier,
	) error
}
