package domain

import "context"

type TransactionPort interface {
	FindAll(ctx context.Context) ([]Transaction, error)
}
