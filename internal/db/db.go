package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"scurvy10k/internal/utils"
	sqlc "scurvy10k/sql/gen"
)

//go:generate mockgen -destination=../mocks/db_mocks.go -package=mock_db scurvy10k/internal/db Database,Connection,Queries

type Database interface {
	Connect(ctx context.Context) (Connection, error)
}

type Connection interface {
	Close(ctx context.Context) error
	Queries() Queries
}

type Queries interface {
	NumberOfPlayers(ctx context.Context) (int64, error)
	AddPlayer(ctx context.Context, name string) (sqlc.Player, error)
	GetIdOfPlayer(ctx context.Context, name string) (int32, error)

	GetAllDebts(ctx context.Context) ([]sqlc.GetAllDebtsRow, error)
	GetDebt(ctx context.Context, id pgtype.Int4) (sqlc.Debt, error)
	SetDebt(ctx context.Context, params sqlc.SetDebtParams) (sqlc.Debt, error)
	UpdateDebt(ctx context.Context, params sqlc.UpdateDebtParams) (sqlc.Debt, error)

	GetBotSetup(ctx context.Context) (sqlc.BotSetup, error)
	PutBotSetup(ctx context.Context, params sqlc.PutBotSetupParams) (sqlc.BotSetup, error)
}

type database struct {
}

func NewDatabase() Database {
	return &database{}
}

func (d database) Connect(ctx context.Context) (Connection, error) {
	conn, err := utils.GetConnection(ctx, utils.DefaultConfig())
	if err != nil {
		return nil, fmt.Errorf("could not establish db connection: %w", err)
	}
	return connection{conn}, nil
}

type connection struct {
	conn *pgx.Conn
}

func (c connection) Close(ctx context.Context) error {
	return c.conn.Close(ctx)
}

func (c connection) Queries() Queries {
	return sqlc.New(c.conn)
}

func IdType(id int32) pgtype.Int4 {
	return pgtype.Int4{
		Int32: id,
		Valid: true,
	}
}
