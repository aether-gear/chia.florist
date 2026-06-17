package repository

import (
	"context"

	"service-core/internal/modules/user/domain"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type UserRepository interface {
	FindCustomers(
		ctx context.Context,
		exec transaction.Executor,
		params FindUserParams,
	) ([]domain.User, int, error)

	GetByID(
		ctx context.Context,
		exec transaction.Executor,
		id uuid.UUID,
	) (*domain.User, error)
	GetByUsername(
		ctx context.Context,
		exec transaction.Executor,
		username string,
	) (*domain.User, error)

	CreateUser(
		ctx context.Context,
		exec transaction.Executor,
		props CreateUserProps,
	) error
}
