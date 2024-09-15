package handler

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"slash10k/internal/db"
	"slash10k/internal/models"
)

func GetJournalEntries(d db.Database) func(c echo.Context) error {
	return func(c echo.Context) error {
		name := c.Param("player")
		if name == "" {
			return c.String(400, "name is required!")
		}

		ctx := c.Request().Context()
		conn, err := d.Connect(ctx)
		defer func(conn db.Connection, ctx context.Context) {
			_ = conn.Close(ctx)
		}(conn, ctx)
		if err != nil {
			return fmt.Errorf("could not get db connection: %w", err)
		}

		id, err := conn.Queries().GetIdOfPlayer(ctx, name)
		if err != nil {
			return c.String(500, "could not get player id")
		}

		entries, err := conn.Queries().GetJournalEntries(ctx, db.IdType(id))
		if err != nil {
			return c.String(500, "could not get journal entries")
		}

		var journalEntries models.JournalEntries
		for _, entry := range entries {
			journalEntries.Entries = append(
				journalEntries.Entries, models.JournalEntry{
					Amount: int(entry.Amount),
					Reason: entry.Description,
					Date:   entry.Date.Time.Unix(),
				},
			)
		}

		return c.JSON(http.StatusOK, journalEntries)
	}
}
