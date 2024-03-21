package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"os"
)

const (
	MigrationsDir = "sql/migrations"
)

var (
	ErrMigrationDirNotExist = errors.New("migration directory 'sql/migrations' does not exist")
)

func Migrate(ctx context.Context, connectionString string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("could not set goose dialect 'postgres': %w", err)
	}

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return fmt.Errorf("could not open db connection: %w", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	if _, err = os.Stat(MigrationsDir); errors.Is(err, os.ErrNotExist) {
		return ErrMigrationDirNotExist
	}

	if err = goose.UpContext(ctx, db, MigrationsDir); err != nil {
		return fmt.Errorf("could not run goose up: %w", err)
	}

	return nil
}
