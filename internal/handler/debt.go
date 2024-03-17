package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"scurvy10k/internal/models"
	"scurvy10k/internal/utils"
	sqlc "scurvy10k/sql/gen"
	frontend "scurvy10k/templ"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

func GetDebt(c echo.Context) error {
	return c.String(200, "GetDebt!")
}

var (
	ErrDebtTooHigh  = errors.New("debt cannot be more than 1_000_000")
	ErrDebtNegative = errors.New("debt cannot be negative")
)

func AllDebts(c echo.Context) error {
	debts, err := allDebts()
	if err != nil {
		log.Err(err).Msg("could not get all debts")
		return c.String(500, "could not get all debts")
	}
	switch c.Request().Header.Get("Accept") {
	case "text/html":
		b := strings.Builder{}
		for _, debt := range debts {
			err = frontend.DebtView(debt).Render(context.Background(), &b)
			if err != nil {
				log.Err(err).Msg("could not render debt")
				return c.String(500, "could not render debt")
			}
		}
		return c.HTML(http.StatusOK, b.String())
	case "application/json":
	case "":
		return c.JSON(http.StatusOK, models.AllDebtsResponse{
			Debts: debts,
		})
	}
	return c.String(400, "unsupported accept header")
}

func allDebts() ([]models.PlayerDebt, error) {
	conn, err := utils.GetConnection(utils.DefaultConfig())
	if err != nil {
		return nil, fmt.Errorf("could not get db connection: %w", err)
	}
	q := sqlc.New(conn)
	dbDebts, err := q.AllPlayerDebts(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not get all player debts: %w", err)
	}
	var debts []models.PlayerDebt
	for _, debtRow := range dbDebts {
		debts = append(debts, models.PlayerDebt{
			Name:   debtRow.Name,
			Amount: fmt.Sprint(debtRow.Amount),
		})
	}
	return debts, nil
}

func AddDebt(c echo.Context) error {
	name := c.Param("player")
	if name == "" {
		return c.String(400, "name is required!")
	}

	a := c.Param("amount")
	if a == "" {
		return c.String(400, "amount is required!")
	}
	amount, err := strconv.ParseInt(a, 10, 64)
	if err != nil {
		return c.String(400, "amount must be an integer!")
	}

	err = addDebtToPlayer(name, amount)
	if err != nil {
		log.Err(err).Msg("could not add debt to player")
		return c.String(500, "could not add debt to player")
	}

	return AllDebts(c)
}

func addDebtToPlayer(name string, amount int64) error {
	conn, err := utils.GetConnection(utils.DefaultConfig())
	if err != nil {
		return fmt.Errorf("could not get db connection: %w", err)
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		_ = conn.Close(ctx)
	}(conn, context.Background())

	q := sqlc.New(conn)
	pId, err := q.GetIdOfPlayer(context.Background(), name)
	if err != nil {
		return fmt.Errorf("could not get player id for %s: %w", name, err)
	}
	currentDebt, err := q.GetDebt(context.Background(), pgtype.Int4{
		Int32: pId,
		Valid: true,
	})
	if err != nil {
		return fmt.Errorf("could not get debt for player (id:%v): %w", pId, err)
	}

	newAmount := currentDebt.Amount + amount
	if newAmount < 0 {
		return ErrDebtNegative
	}
	if newAmount > 1_000_000 {
		return ErrDebtTooHigh
	}
	_, err = q.UpdateDebt(context.Background(), sqlc.UpdateDebtParams{
		Amount: currentDebt.Amount + amount,
		UserID: pgtype.Int4{
			Int32: pId,
			Valid: true,
		},
	})
	if err != nil {
		return fmt.Errorf("could not set debt for player (id:%v): %w", pId, err)
	}
	return nil
}
