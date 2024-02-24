package handler

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"scurvy10k/sql/db"
	"scurvy10k/src/config"
	"scurvy10k/src/utils"
	"strconv"
)

func GetDebt(c echo.Context) error {
	return c.String(200, "GetDebt!")
}

func AddDebt(c echo.Context) error {
	name := c.Param("player")
	if name == "" {
		return c.String(400, "Name is required!")
	}

	amount := c.Param("amount")
	if amount == "" {
		return c.String(400, "Amount is required!")
	}
	a, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		return c.String(400, "Amount must be an integer!")
	}

	conn, err := utils.GetConnection(config.NewConfig())
	if err != nil {
		return c.String(500, "Could not get db connection!")
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		_ = conn.Close(ctx)
	}(conn, context.Background())

	q := db.New(conn)
	pId, err := q.GetIdOfPlayer(context.Background(), name)
	if err != nil {
		return c.String(400, "Could not get player id!")
	}
	_, err = q.AddDebt(context.Background(), db.AddDebtParams{
		Amount:      a,
		Description: "",
		UserID: pgtype.Int4{
			Int32: pId,
			Valid: true,
		},
	})
	if err != nil {
		return c.String(500, "Could not add player!")
	}

	return c.String(200, fmt.Sprintf("added debt %s to player %s", amount, name))
}
