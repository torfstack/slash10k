package db

import (
	"context"
	"github.com/jackc/pgx/v5"
)

func GetConnection(ctx context.Context, connectionString string) (*pgx.Conn, error) {
	return pgx.Connect(ctx, connectionString)
}
