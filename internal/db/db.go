package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New() (*pgxpool.Pool, error) {
	return pgxpool.New(
		context.Background(),
		"postgres://vinyl:vinyl@localhost:5432/vinyl?sslmode=disable",
	)
}
