package handler

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"scurvy10k/internal/utils"
	"scurvy10k/sql/db"
	frontend "scurvy10k/templ"
)

func ServeFrontend(c echo.Context) error {
	conn, err := utils.GetConnection(utils.DefaultConfig())
	if err != nil {
		return c.String(500, "Could not get db connection!")
	}
	q := db.New(conn)
	debts, err := q.AllPlayerDebts(context.Background())
	if err != nil {
		return c.String(500, "Could not get debt!")
	}
	log.Debug().Msgf("retrieved %v debts", len(debts))

	return frontend.
		Debt().
		Render(context.Background(), c.Response())
}
