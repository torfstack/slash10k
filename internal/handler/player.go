package handler

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"scurvy10k/internal/utils"
	"scurvy10k/sql/db"
)

func AddPlayer(c echo.Context) error {
	name := c.Param("name")
	if name == "" {
		return c.String(400, "Name is required!")
	}

	conn, err := utils.GetConnection(utils.DefaultConfig())
	if err != nil {
		return c.String(500, "Could not get db connection!")
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		_ = conn.Close(ctx)
	}(conn, context.Background())

	q := db.New(conn)

	count, err := q.NumberOfPlayers(context.Background())
	if err != nil {
		return c.String(500, "Could not get number of players!")
	}
	if count >= 100 {
		return c.String(400, "Max number of players reached!")
	}

	p, err := q.AddPlayer(context.Background(), name)
	if err != nil {
		return c.String(500, "Could not add player!")
	}
	_, err = q.SetDebt(context.Background(), db.SetDebtParams{
		Amount: 0,
		UserID: pgtype.Int4{
			Int32: p.ID,
			Valid: true,
		},
	})
	if err != nil {
		return c.String(500, "Could not add zero debt to player!")
	}

	return c.String(200, fmt.Sprintf("added player %v", p))
}

func DeletePlayer(c echo.Context) error {
	return c.String(200, "DeletePlayer!")
}
