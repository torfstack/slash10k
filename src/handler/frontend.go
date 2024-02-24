package handler

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"scurvy10k/sql/db"
	"scurvy10k/src/config"
	"scurvy10k/src/utils"
	frontend "scurvy10k/templ"
	"strconv"
)

func ServeFrontend(c echo.Context) error {
	conn, err := utils.GetConnection(config.NewConfig())
	if err != nil {
		return c.String(500, "Could not get db connection!")
	}
	q := db.New(conn)
	sId, err := q.GetIdOfPlayer(context.Background(), "scurvy")
	if err != nil {
		return c.String(500, "Could not get player id!")
	}
	d, err := q.GetDebt(context.Background(), pgtype.Int4{
		Int32: sId,
		Valid: true,
	})
	if err != nil {
		return c.String(500, "Could not get debt!")
	}

	return frontend.
		Debt("scurvy", strconv.FormatInt(d.Amount, 10)).
		Render(context.Background(), c.Response())
}
