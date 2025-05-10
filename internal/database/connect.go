package database

import (
	"context"
	// "os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect возвращает *pgxpool.Pool c выключенным statement‑кэшем.
func Connect(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	// ← главное: просить pgx всегда использовать простой протокол
	cfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	// (необязательно, но полезно)
	cfg.MaxConns           = 20
	cfg.MinConns           = 2
	cfg.HealthCheckPeriod  = 30 * time.Second
	cfg.MaxConnIdleTime    = 5 * time.Minute
	cfg.MaxConnLifetime    = 30 * time.Minute

	return pgxpool.NewWithConfig(ctx, cfg)
}