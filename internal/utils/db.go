package utils

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"os"
	"scurvy10k/internal/config"
)

func GetConnection(config config.Config) (*pgx.Conn, error) {
	s := config.ConnectionString
	return pgx.Connect(context.Background(), s)
}

func DefaultConfig() config.Config {
	host := os.Getenv("DATABASE_CONNECTION_HOST")
	port := os.Getenv("DATABASE_CONNECTION_PORT")
	user := os.Getenv("DATABASE_CONNECTION_USER")
	password := os.Getenv("DATABASE_CONNECTION_PASSWORD")
	dbname := os.Getenv("DATABASE_CONNECTION_DBNAME")

	atLeastOneIsEmpty := host == "" || port == "" || user == "" || dbname == "" || password == ""

	if atLeastOneIsEmpty {
		log.Debug().Msg("using default config")
		return config.NewConfig()
	} else {
		log.Debug().Msg("using env config")
		connString := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
		return config.NewConfig(
			config.WithConnectionString(connString),
		)
	}
}
