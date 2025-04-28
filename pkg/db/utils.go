package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func GetConnection(ctx context.Context, connectionString string) (*pgx.Conn, error) {
	return pgx.Connect(ctx, connectionString)
}

func ProcessOutbox(ctx context.Context, db Database, handler func(eventType string, payload map[string]interface{}) error) error {
	conn, err := db.Connect(ctx)
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}
	defer conn.Close(ctx)

	rows, err := conn.Queries().Query(ctx, "SELECT id, event_type, payload FROM outbox WHERE processed_at IS NULL")
	if err != nil {
		return fmt.Errorf("could not query outbox: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var eventType string
		var payload map[string]interface{}
		if err := rows.Scan(&id, &eventType, &payload); err != nil {
			return fmt.Errorf("could not scan outbox row: %w", err)
		}

		if err := handler(eventType, payload); err != nil {
			return fmt.Errorf("could not handle outbox event: %w", err)
		}

		if _, err := conn.Queries().Exec(ctx, "UPDATE outbox SET processed_at = now() WHERE id = $1", id); err != nil {
			return fmt.Errorf("could not mark outbox event as processed: %w", err)
		}
	}

	return nil
}
