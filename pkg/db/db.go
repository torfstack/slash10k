package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	sqlc "slash10k/sql/gen"
)

//go:generate mockgen -destination=../mocks/db_mocks.go -package=mock_db slash10k/internal/db Database,Connection,Queries,Transaction

type Database interface {
	Connect(ctx context.Context) (Connection, error)
}

type Connection interface {
	Close(ctx context.Context)
	StartTransaction(ctx context.Context) (Transaction, error)
	Queries() Queries
}

type Transaction interface {
	Commit(ctx context.Context) error
	Queries() Queries
}

type Queries interface {
	NumberOfPlayers(ctx context.Context) (int64, error)
	AddPlayer(ctx context.Context, param sqlc.AddPlayerParams) (sqlc.Player, error)
	DeletePlayer(ctx context.Context, id int32) error
	GetIdOfPlayer(ctx context.Context, param sqlc.GetIdOfPlayerParams) (int32, error)
	GetPlayer(ctx context.Context, params sqlc.GetPlayerParams) ([]sqlc.GetPlayerRow, error)
	GetAllPlayers(ctx context.Context, guildId string) ([]sqlc.GetAllPlayersRow, error)
	DoesPlayerExist(ctx context.Context, params sqlc.DoesPlayerExistParams) (bool, error)

	SetDebt(ctx context.Context, params sqlc.SetDebtParams) error

	AddJournalEntry(ctx context.Context, params sqlc.AddJournalEntryParams) (sqlc.DebtJournal, error)
	GetJournalEntries(ctx context.Context, params int32) ([]sqlc.DebtJournal, error)
	UpdateJournalEntry(ctx context.Context, params sqlc.UpdateJournalEntryParams) (sqlc.DebtJournal, error)
	DeleteJournalEntry(ctx context.Context, id int32) error

	GetBotSetup(ctx context.Context, guildId string) (sqlc.BotSetup, error)
	DoesBotSetupExist(ctx context.Context, guildId string) (bool, error)
	PutBotSetup(ctx context.Context, params sqlc.PutBotSetupParams) (sqlc.BotSetup, error)
	DeleteBotSetup(ctx context.Context, guildId string) error
	GetAllBotSetups(ctx context.Context) ([]sqlc.BotSetup, error)
}

type database struct {
	connectionString string
}

func NewDatabase(connectionString string) Database {
	return &database{connectionString: connectionString}
}

func (d database) Connect(ctx context.Context) (Connection, error) {
	conn, err := GetConnection(ctx, d.connectionString)
	if err != nil {
		return nil, fmt.Errorf("could not establish db connection: %w", err)
	}
	return connection{conn, make([]transaction, 0)}, nil
}

type connection struct {
	conn *pgx.Conn
	txs  []transaction
}

func (c connection) Close(ctx context.Context) {
	for _, tx := range c.txs {
		if !tx.didCommit {
			_ = tx.tx.Rollback(ctx)
		}
	}
	_ = c.conn.Close(ctx)
}

func (c connection) StartTransaction(ctx context.Context) (Transaction, error) {
	tx, err := c.conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not start transaction: %w", err)
	}
	ts := transaction{tx: tx}
	c.txs = append(c.txs, ts)
	return ts, nil
}

func (c connection) Queries() Queries {
	return sqlc.New(c.conn)
}

type transaction struct {
	tx        pgx.Tx
	didCommit bool
}

func (t transaction) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t transaction) Queries() Queries {
	return sqlc.New(t.tx)
}
