package pgsql

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type db struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, connString string) (*db, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	return &db{
		pool: pool,
	}, nil
}
