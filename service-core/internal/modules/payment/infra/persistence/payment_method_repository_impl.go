package persistence

import (
	database "service-core/internal/infra/db"
	"service-core/internal/modules/payment/domain"
	"service-core/internal/modules/payment/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentMethodRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewPaymentMethodRepository(conn *database.Connection) repository.PaymentMethodRepository {
	return &PaymentMethodRepositoryImpl{
		db: conn.Pool,
	}
}

func (r *PaymentMethodRepositoryImpl) Save(paymentMethod domain.PaymentMethod) error {
	return nil
}

func (r *PaymentMethodRepositoryImpl) FindByName(name string) (*domain.PaymentMethod, error) {

}

func (r *PaymentMethodRepositoryImpl) GetByID(paymentID uuid.UUID) (*domain.PaymentMethod, error) {

}
