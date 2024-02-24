package utils

import (
	"context"
	"github.com/jackc/pgx/v5"
	"scurvy10k/src/config"
)

func GetConnection(config config.Config) (*pgx.Conn, error) {
	s := config.ConnectionString
	return pgx.Connect(context.Background(), s)
}
