package application

import (
	"context"
	"service-core/internal/modules/transaction/domain"
)

type TransactionUsecase struct {
	repo domain.TransactionPort
}

func NewTransactionUsecase(repo domain.TransactionPort) *TransactionUsecase {
	return &TransactionUsecase{repo: repo}
}

func (u *TransactionUsecase) GetAll(ctx context.Context) ([]domain.Transaction, error) {
	return u.repo.FindAll(ctx)
}
