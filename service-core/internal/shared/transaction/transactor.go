package transaction

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Executor abstracts database query execution.
//
// Both connection pools and database transactions implement this
// contract, allowing repositories to execute queries without being
// coupled to a specific execution context.
type Executor interface {
	Exec(
		ctx context.Context,
		sql string,
		args ...any,
	) (pgconn.CommandTag, error)

	Query(
		ctx context.Context,
		sql string,
		args ...any,
	) (pgx.Rows, error)

	QueryRow(
		ctx context.Context,
		sql string,
		args ...any,
	) pgx.Row
}

// Transactor provides transactional execution for application services.
//
// It enables use cases to define transaction boundaries while delegating
// transaction lifecycle management to the infrastructure layer.
type Transactor interface {
	WithinTransaction(
		ctx context.Context,
		fn func(Executor) error,
	) error
}
