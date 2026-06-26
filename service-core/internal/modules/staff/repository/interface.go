package repository

import (
	"context"

	"service-core/internal/modules/staff/domain"
	transaction "service-core/internal/shared/transaction"
)

type StaffRepository interface {
	Create(
		ctx context.Context,
		exec transaction.Executor,
		staff domain.Staff,
	) error

	FindStaff(
		ctx context.Context,
		exec transaction.Executor,
		params FindStaffParams,
	) ([]domain.Staff, int, error)
}
