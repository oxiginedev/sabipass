package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/oxiginedev/sabipass/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type DB struct {
	*bun.DB
	queryTimeout time.Duration
}

func NewDB(cfg *config.Config) (*DB, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.Database.Postgres.DSN)))

	db := bun.NewDB(sqldb, pgdialect.New())
	err := db.Ping()
	if err != nil {
		return nil, fmt.Errorf("postgres: failed to ping database: %w", err)
	}

	return &DB{db, cfg.Database.Postgres.QueryTimeout}, nil
}

func (db *DB) WithContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, db.queryTimeout)
}
