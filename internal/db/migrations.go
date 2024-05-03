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
	DefaultMigrationsDir = "sql/migrations"
)

type MigrateOpts func(*MigrateOptions)

type MigrateOptions struct {
	MigrationsDir string
}

func WithMigrationsDir(dir string) MigrateOpts {
	return func(opts *MigrateOptions) {
		opts.MigrationsDir = dir
	}
}

func Migrate(ctx context.Context, connectionString string, opts ...MigrateOpts) error {
	migrateOpts := &MigrateOptions{
		MigrationsDir: DefaultMigrationsDir,
	}
	for _, opt := range opts {
		opt(migrateOpts)
	}

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

	if _, err = os.Stat(migrateOpts.MigrationsDir); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("specified migration directory '%s' does not exist", migrateOpts.MigrationsDir)
	}

	if err = goose.UpContext(ctx, db, migrateOpts.MigrationsDir); err != nil {
		return fmt.Errorf("could not run goose up: %w", err)
	}

	return nil
}
