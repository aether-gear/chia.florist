package database

import (
	"context"
	"fmt"
	"time"

	"service-core/internal/shared/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Connection struct {
	Pool *pgxpool.Pool
}

func NewConnection(cfg config.DatabaseConfig) (*Connection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(*cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("parse pool config: %w", err)
	}

	poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeExec
	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		_, err := conn.Exec(ctx, "SET timezone TO 'Asia/Jakarta'")
		return err
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Connection{
		Pool: pool,
	}, nil
}

func (c *Connection) Close() {
	c.Pool.Close()
}
