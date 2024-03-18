package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"scurvy10k/internal/utils"
	sqlc "scurvy10k/sql/gen"
)

type Database interface {
	Connect(ctx context.Context) (Connection, error)
}

type Connection interface {
	Close(ctx context.Context) error

	NumberOfPlayers(ctx context.Context) (int, error)
	AddPlayer(ctx context.Context, name string) (*sqlc.Player, error)
	GetIdOfPlayer(ctx context.Context, name string) (int32, error)

	GetAllDebts(ctx context.Context) ([]sqlc.AllPlayerDebtsRow, error)
	GetDebt(ctx context.Context, id int32) (*sqlc.Debt, error)
	SetDebt(ctx context.Context, params sqlc.SetDebtParams) error
	UpdateDebt(ctx context.Context, params sqlc.UpdateDebtParams) error

	GetBotSetup(ctx context.Context) (*sqlc.BotSetup, error)
	PutBotSetup(ctx context.Context, params sqlc.PutBotSetupParams) error
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
	c := connection{
		conn: conn,
	}
	return &c, nil
}

type connection struct {
	conn *pgx.Conn
}

func (c connection) NumberOfPlayers(ctx context.Context) (int, error) {
	numberOfPlayers, err := sqlc.New(c.conn).NumberOfPlayers(ctx)
	if err != nil {
		return 0, err
	}
	return int(numberOfPlayers), nil
}

func (c connection) AddPlayer(ctx context.Context, name string) (*sqlc.Player, error) {
	player, err := sqlc.New(c.conn).AddPlayer(ctx, name)
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func (c connection) GetIdOfPlayer(ctx context.Context, name string) (int32, error) {
	id, err := sqlc.New(c.conn).GetIdOfPlayer(ctx, name)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (c connection) Close(ctx context.Context) error {
	return c.conn.Close(ctx)
}

func (c connection) GetAllDebts(ctx context.Context) ([]sqlc.AllPlayerDebtsRow, error) {
	debts, err := sqlc.New(c.conn).AllPlayerDebts(ctx)
	if err != nil {
		return nil, err
	}
	return debts, nil
}

func (c connection) GetDebt(ctx context.Context, id int32) (*sqlc.Debt, error) {
	debt, err := sqlc.New(c.conn).GetDebt(ctx, idType(id))
	if err != nil {
		return nil, err
	}
	return &debt, nil
}

func (c connection) SetDebt(ctx context.Context, params sqlc.SetDebtParams) error {
	_, err := sqlc.New(c.conn).SetDebt(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

func (c connection) UpdateDebt(ctx context.Context, params sqlc.UpdateDebtParams) error {
	_, err := sqlc.New(c.conn).UpdateDebt(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

func (c connection) GetBotSetup(ctx context.Context) (*sqlc.BotSetup, error) {
	botSetup, err := sqlc.New(c.conn).GetBotSetup(ctx)
	if err != nil {
		return nil, err
	}
	return &botSetup, nil
}

func (c connection) PutBotSetup(ctx context.Context, params sqlc.PutBotSetupParams) error {
	_, err := sqlc.New(c.conn).PutBotSetup(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

func idType(id int32) pgtype.Int4 {
	return pgtype.Int4{
		Int32: id,
		Valid: true,
	}
}
