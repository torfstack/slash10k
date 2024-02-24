package handler

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"scurvy10k/sql/db"
	"scurvy10k/src/models"
	"scurvy10k/src/utils"
	frontend "scurvy10k/templ"
	"strconv"
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
		Debt(transformDebt(debts)).
		Render(context.Background(), c.Response())
}

func transformDebt(debts []db.AllPlayerDebtsRow) []models.PlayerDebt {
	var transformed []models.PlayerDebt
	for _, d := range debts {
		transformed = append(transformed, models.PlayerDebt{
			Name:   d.Name,
			Amount: strconv.FormatInt(d.Amount, 10),
		})
	}
	return transformed
}
