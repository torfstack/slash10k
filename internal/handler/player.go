package handler

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"slash10k/internal/db"
	"slash10k/sql/gen"
)

func AddPlayer(d db.Database) func(c echo.Context) error {
	return func(c echo.Context) error {
		name := c.Param("name")
		if name == "" {
			return c.String(400, "name is required!")
		}

		conn, err := d.Connect(c.Request().Context())
		if err != nil {
			log.Err(err).Msg("could not get db connection!")
			return c.String(500, "could not get db connection!")
		}
		defer func(conn db.Connection, ctx context.Context) {
			_ = conn.Close(ctx)
		}(conn, context.Background())

		count, err := conn.Queries().NumberOfPlayers(context.Background())
		if err != nil {
			log.Err(err).Msg("could not get number of players!")
			return c.String(500, "could not get number of players!")
		}
		if count >= 100 {
			return c.String(400, "Max number of players reached!")
		}

		p, err := conn.Queries().AddPlayer(context.Background(), name)
		if err != nil {
			log.Err(err).Msg("could not add player!")
			return c.String(500, "Could not add player!")
		}
		_, err = conn.Queries().SetDebt(context.Background(), sqlc.SetDebtParams{
			Amount: 0,
			UserID: pgtype.Int4{
				Int32: p.ID,
				Valid: true,
			},
		})
		if err != nil {
			log.Err(err).Msg("could not add zero debt to player!")
			return c.String(500, "Could not add zero debt to player!")
		}

		log.Info().Msgf("added player %s with id %v", p.Name, p.ID)
		return c.String(200, fmt.Sprintf("added player %s with id %v", p.Name, p.ID))
	}
}

func DeletePlayer(c echo.Context) error {
	return c.String(200, "DeletePlayer!")
}
