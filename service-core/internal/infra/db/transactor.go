package database

import (
	"context"

	transaction "service-core/internal/shared/transaction"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresTransactor struct {
	db *pgxpool.Pool
}

func NewPostgresTransactor(db *pgxpool.Pool) transaction.Transactor {
	return &PostgresTransactor{
		db: db,
	}
}

func (t *PostgresTransactor) WithinTransaction(
	ctx context.Context,
	fn func(transaction.Executor) error,
) error {
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
