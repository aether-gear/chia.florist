package infra

import (
	"context"
	"service-core/internal/features/transaction/domain"
)

type TransactionRepositoryImpl struct {
	data []domain.Transaction
}

func NewTransactionRepository() *TransactionRepositoryImpl {
	return &TransactionRepositoryImpl{
		data: []domain.Transaction{
			{
				ID:     "trx-001",
				Amount: 150000,
				Status: "SUCCESS",
			},
			{
				ID:     "trx-002",
				Amount: 275000,
				Status: "PENDING",
			},
		},
	}
}

func (r *TransactionRepositoryImpl) FindAll(ctx context.Context) ([]domain.Transaction, error) {
	return r.data, nil
}
