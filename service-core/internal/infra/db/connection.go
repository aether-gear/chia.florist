package database

import (
	"context"
	"fmt"
	"time"

	"service-core/internal/shared/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Connection struct {
	Pool *pgxpool.Pool
}

func NewConnection(cfg config.DatabaseConfig) (*Connection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, *cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &Connection{
		Pool: pool,
	}, nil
}

func (c *Connection) Close() {
	c.Pool.Close()
}
