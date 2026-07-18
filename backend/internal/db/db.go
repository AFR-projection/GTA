package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	cfg, err := parsePoolConfig(databaseURL)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return pool, nil
}

func Migrate(databaseURL, migrationsDir string) error {
	cfg, err := parsePoolConfig(databaseURL)
	if err != nil {
		return err
	}
	sqlDB := stdlib.OpenDB(*cfg.ConnConfig)
	defer sqlDB.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	if err := goose.Up(sqlDB, migrationsDir); err != nil {
		return fmt.Errorf("migrate up: %w", err)
	}
	return nil
}

// OpenSQL is used by goose tooling helpers if needed.
func OpenSQL(databaseURL string) (*sql.DB, error) {
	cfg, err := parsePoolConfig(databaseURL)
	if err != nil {
		return nil, err
	}
	return stdlib.OpenDB(*cfg.ConnConfig), nil
}

// Neon (and many poolers) break with prepared statements.
// Force simple protocol so migrate/runtime stay stable.
func parsePoolConfig(databaseURL string) (*pgxpool.Config, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}
	cfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	return cfg, nil
}
