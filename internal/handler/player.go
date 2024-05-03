package handler

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
	"slash10k/internal/db"
	"slash10k/internal/utils"
	"slash10k/sql/gen"
	"strings"
)

func AddPlayer(d db.Database) func(c echo.Context) error {
	return func(c echo.Context) error {
		name := c.Param("name")
		if name == "" {
			return c.String(http.StatusBadRequest, "name is required!")
		}

		conn, err := d.Connect(c.Request().Context(), utils.DefaultConfig().ConnectionString)
		if err != nil {
			log.Err(err).Msg("could not get db connection!")
			return c.String(http.StatusInternalServerError, "could not get db connection!")
		}
		defer func(conn db.Connection, ctx context.Context) {
			_ = conn.Close(ctx)
		}(conn, context.Background())

		count, err := conn.Queries().NumberOfPlayers(context.Background())
		if err != nil {
			log.Err(err).Msg("could not get number of players!")
			return c.String(http.StatusInternalServerError, "could not get number of players!")
		}
		if count >= 100 {
			return c.String(http.StatusBadRequest, "Max number of players reached!")
		}

		name = strings.ToLower(name)
		p, err := conn.Queries().AddPlayer(context.Background(), name)
		if err != nil {
			log.Err(err).Msg("could not add player!")
			return c.String(http.StatusInternalServerError, "Could not add player!")
		}
		err = conn.Queries().SetDebt(context.Background(), sqlc.SetDebtParams{
			Amount: 0,
			UserID: pgtype.Int4{
				Int32: p.ID,
				Valid: true,
			},
		})
		if err != nil {
			log.Err(err).Msg("could not add zero debt to player!")
			return c.String(http.StatusInternalServerError, "Could not add zero debt to player!")
		}

		log.Info().Msgf("added player %s with id %v", p.Name, p.ID)
		return c.String(http.StatusNoContent, fmt.Sprintf("added player %s with id %v", p.Name, p.ID))
	}
}

func DeletePlayer(d db.Database) func(c echo.Context) error {
	return func(c echo.Context) error {
		name := c.Param("name")
		if name == "" {
			return c.String(http.StatusBadRequest, "name is required!")
		}

		conn, err := d.Connect(c.Request().Context(), utils.DefaultConfig().ConnectionString)
		if err != nil {
			log.Err(err).Msg("could not get db connection!")
			return c.String(http.StatusInternalServerError, "could not get db connection!")
		}
		defer func(conn db.Connection, ctx context.Context) {
			_ = conn.Close(ctx)
		}(conn, context.Background())

		name = strings.ToLower(name)
		err = conn.Queries().DeletePlayer(context.Background(), name)
		if err != nil {
			log.Err(err).Msg("could not delete player!")
			return c.String(http.StatusInternalServerError, "Could not delete player!")
		}

		log.Info().Msgf("deleted player %s", name)
		return c.String(http.StatusNoContent, fmt.Sprintf("deleted player %s", name))
	}
}
