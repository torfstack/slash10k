package handler

import (
	"context"
	"errors"
	"fmt"
	"scurvy10k/sql/db"
	"scurvy10k/src/utils"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func GetDebt(c echo.Context) error {
	return c.String(200, "GetDebt!")
}

var (
	ErrNameNotSpecified   = errors.New("name is required")
	ErrAmountNotSpecified = errors.New("amount is required")
)

func AddDebt(c echo.Context) error {
	name, amount, err := nameAndAmount(c)
	if err != nil {
		return err
	}

	err = addDebtToPlayer(name, amount, c)
	if err != nil {
		return err
	}

	return c.String(200, fmt.Sprintf("added debt %v to player %s", amount, name))
}

func nameAndAmount(c echo.Context) (string, int64, error) {
	name := c.Param("player")
	if name == "" {
		_ = c.String(400, "Name is required!")
		return "", 0, ErrNameNotSpecified
	}

	amount := c.Param("amount")
	if amount == "" {
		_ = c.String(400, "Amount is required!")
		return "", 0, ErrAmountNotSpecified
	}
	a, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		return "", 0, c.String(400, "Amount must be an integer!")
	}
	return name, a, nil
}

func addDebtToPlayer(name string, amount int64, c echo.Context) error {
	conn, err := utils.GetConnection(utils.DefaultConfig())
	if err != nil {
		_ = c.String(500, "Could not get db connection!")
		return err
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		_ = conn.Close(ctx)
	}(conn, context.Background())

	q := db.New(conn)
	pId, err := q.GetIdOfPlayer(context.Background(), name)
	if err != nil {
		log.Error().Msgf("could not get player id for %s: %v", name, err)
		_ = c.String(400, "Could not get player id!")
		return err
	}
	currentDebt, err := q.GetDebt(context.Background(), pgtype.Int4{
		Int32: pId,
		Valid: true,
	})
	if err != nil {
		log.Error().Msgf("could not get debt for player %s(id:%v): %s", name, pId, err)
		_ = c.String(400, "Could not get player debt!")
		return err
	}
	newAmount := currentDebt.Amount + amount
	if newAmount < 0 {
		_ = c.String(400, "Debt cannot be negative!")
		return errors.New("debt cannot be negative")
	}
	_, err = q.UpdateDebt(context.Background(), db.UpdateDebtParams{
		Amount: currentDebt.Amount + amount,
		UserID: pgtype.Int4{
			Int32: pId,
			Valid: true,
		},
	})
	if err != nil {
		log.Error().Msgf("could not set debt for player %s(id:%v): %s", name, pId, err)
		_ = c.String(400, "Could not set player debt!")
		return err
	}
	return nil
}
