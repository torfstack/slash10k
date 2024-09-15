package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"slash10k/internal/db"
	"slash10k/internal/models"
	sqlc "slash10k/sql/gen"
	frontend "slash10k/templ"
	"strconv"
	"strings"

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

func AllDebts(d db.Database) func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		conn, err := d.Connect(ctx)
		defer func(conn db.Connection, ctx context.Context) {
			_ = conn.Close(ctx)
		}(conn, ctx)
		if err != nil {
			return fmt.Errorf("could not get db connection: %w", err)
		}

		return allDebtsRespond(c, conn)
	}
}

func allDebtsRespond(c echo.Context, conn db.Connection) error {
	debts, err := allDebts(conn, c.Request().Context())
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
		return c.JSON(
			http.StatusOK, models.AllDebtsResponse{
				Debts: debts,
			},
		)
	}
	log.Warn().Str("accept", c.Request().Header.Get("Accept")).Msg("unsupported accept header")
	return c.String(400, "unsupported accept header")
}

func allDebts(conn db.Connection, ctx context.Context) ([]models.PlayerDebt, error) {
	dbDebts, err := conn.Queries().GetAllDebts(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get all player debts: %w", err)
	}
	var debts []models.PlayerDebt
	for _, debtRow := range dbDebts {
		debts = append(
			debts, models.PlayerDebt{
				Name:   debtRow.Name,
				Amount: fmt.Sprint(debtRow.Amount),
			},
		)
	}
	return debts, nil
}

type AddDebtParams struct {
	Description string `json:"description,omitempty" required:"false"`
}

func AddDebt(d db.Database) func(c echo.Context) error {
	return func(c echo.Context) error {
		name := c.Param("player")
		if name == "" {
			log.Error().Msg("name is required!")
			return c.String(400, "name is required!")
		}

		a := c.Param("amount")
		if a == "" {
			log.Error().Msg("amount is required!")
			return c.String(400, "amount is required!")
		}
		amount, err := strconv.ParseInt(a, 10, 64)
		if err != nil {
			log.Err(err).Msgf("amount must be an integer!, got %v", a)
			return c.String(400, "amount must be an integer!")
		}

		var params AddDebtParams
		if err = c.Bind(&params); err != nil {
			log.Err(err).Msg("could not bind params")
			return c.String(400, "could not bind params")
		}
		if len(params.Description) > 50 {
			log.Error().Msg("description must be less than 50 characters")
			return c.String(400, "description must be less than 50 characters")
		}

		ctx := c.Request().Context()
		conn, err := d.Connect(ctx)
		if err != nil {
			log.Err(err).Msg("could not get db connection!")
			return c.String(500, "could not get db connection!")
		}
		defer func(conn db.Connection, ctx context.Context) {
			_ = conn.Close(ctx)
		}(conn, ctx)

		err = addDebtToPlayer(ctx, conn, name, amount, params.Description)
		if err != nil {
			log.Err(err).Msg("could not add debt to player")
			return c.String(500, "could not add debt to player")
		}

		return allDebtsRespond(c, conn)
	}
}

func addDebtToPlayer(ctx context.Context, conn db.Connection, name string, amount int64, desc string) error {
	tx, err := conn.StartTransaction(ctx)
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}
	defer func(tx db.Transaction, ctx context.Context) {
		_ = tx.Commit(ctx)
	}(tx, ctx)

	queries := tx.Queries()

	pId, err := queries.GetIdOfPlayer(ctx, name)
	if err != nil {
		return fmt.Errorf("could not get player id for %s: %w", name, err)
	}
	currentDebt, err := queries.GetDebt(ctx, db.IdType(pId))
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
	err = queries.SetDebt(
		ctx, sqlc.SetDebtParams{
			Amount: currentDebt.Amount + amount,
			UserID: pgtype.Int4{
				Int32: pId,
				Valid: true,
			},
		},
	)

	if amount > 0 && desc != "" {
		_, err = queries.AddJournalEntry(
			ctx, sqlc.AddJournalEntryParams{
				Amount:      amount,
				Description: desc,
				UserID: pgtype.Int4{
					Int32: pId,
					Valid: true,
				},
			},
		)
		if err != nil {
			return fmt.Errorf("could not add journal entry for player (id:%v): %w", pId, err)
		}
	}
	if amount < 0 {
		var entries []sqlc.DebtJournal
		entries, err = queries.GetJournalEntries(ctx, db.IdType(pId))
		temp := amount
		for _, entry := range entries {
			temp += entry.Amount
			if temp <= 0 {
				err = queries.DeleteJournalEntry(ctx, entry.ID)
				if err != nil {
					return fmt.Errorf("could not delete journal entry for player (id:%v): %w", pId, err)
				}
			} else {
				_, err = queries.UpdateJournalEntry(
					ctx, sqlc.UpdateJournalEntryParams{
						Amount:      temp,
						Description: entry.Description,
						ID:          entry.ID,
					},
				)
				if err != nil {
					return fmt.Errorf("could not update journal entry for player (id:%v): %w", pId, err)
				}
			}
		}
	}

	if err != nil {
		return fmt.Errorf("could not set debt for player (id:%v): %w", pId, err)
	}
	return nil
}
