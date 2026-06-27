package repository

import (
	"context"

	"service-core/internal/modules/staff/domain"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type StaffRepository interface {
	Create(
		ctx context.Context,
		exec transaction.Executor,
		staff domain.Staff,
	) error

	GetByID(
		ctx context.Context,
		exec transaction.Executor,
		id uuid.UUID,
	) (*domain.Staff, error)

	GetProfileByUserID(
		ctx context.Context,
		exec transaction.Executor,
		userID uuid.UUID,
	) (*domain.StaffProfile, error)

	FindStaff(
		ctx context.Context,
		exec transaction.Executor,
		params FindStaffParams,
	) ([]domain.StaffProfile, int, error)
}
