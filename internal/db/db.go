package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"slash10k/internal/utils"
	sqlc "slash10k/sql/gen"
)

//go:generate mockgen -destination=../mocks/db_mocks.go -package=mock_db slash10k/internal/db Database,Connection,Queries,Transaction

type Database interface {
	Connect(ctx context.Context) (Connection, error)
}

type Connection interface {
	Close(ctx context.Context) error
	StartTransaction(ctx context.Context) (Transaction, error)
	Queries() Queries
}

type Transaction interface {
	Commit(ctx context.Context) error
	Queries() Queries
}

type Queries interface {
	NumberOfPlayers(ctx context.Context) (int64, error)
	AddPlayer(ctx context.Context, name string) (sqlc.Player, error)
	DeletePlayer(ctx context.Context, name string) error
	GetIdOfPlayer(ctx context.Context, name string) (int32, error)

	GetAllDebts(ctx context.Context) ([]sqlc.GetAllDebtsRow, error)
	GetDebt(ctx context.Context, id pgtype.Int4) (sqlc.Debt, error)
	SetDebt(ctx context.Context, params sqlc.SetDebtParams) error

	AddJournalEntry(ctx context.Context, params sqlc.AddJournalEntryParams) (sqlc.DebtJournal, error)
	GetJournalEntries(ctx context.Context, id pgtype.Int4) ([]sqlc.DebtJournal, error)
	UpdateJournalEntry(ctx context.Context, params sqlc.UpdateJournalEntryParams) (sqlc.DebtJournal, error)
	DeleteJournalEntry(ctx context.Context, id int32) error

	GetBotSetup(ctx context.Context) (sqlc.BotSetup, error)
	PutBotSetup(ctx context.Context, params sqlc.PutBotSetupParams) (sqlc.BotSetup, error)
}

type database struct {
	connectionString string
}

func NewDatabase(connectionString string) *database {
	return &database{connectionString: connectionString}
}

func (d database) Connect(ctx context.Context) (Connection, error) {
	conn, err := utils.GetConnection(ctx, d.connectionString)
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

func (c connection) StartTransaction(ctx context.Context) (Transaction, error) {
	tx, err := c.conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not start transaction: %w", err)
	}
	return &transaction{tx: tx}, nil
}

func (c connection) Queries() Queries {
	return sqlc.New(c.conn)
}

type transaction struct {
	tx pgx.Tx
}

func (t transaction) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t transaction) Queries() Queries {
	return sqlc.New(t.tx)
}

func IdType(id int32) pgtype.Int4 {
	return pgtype.Int4{
		Int32: id,
		Valid: true,
	}
}
