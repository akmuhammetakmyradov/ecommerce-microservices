package postgresql

import (
	"cart/internal/config"
	"context"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Closer interface {
	Close()
}

type Client interface {
	DB
	Closer
}

type PgClient struct {
	*pgxpool.Pool
}

var _ Client = (*PgClient)(nil)

func NewPostgres(ctx context.Context, cfg *config.Configs) (Client, error) {
	dbUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Postgres.UserName,
		url.QueryEscape(cfg.Postgres.Password),
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DbName,
		cfg.Postgres.Sslmode,
	)

	dbPool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Postgres: %w", err)
	}

	if err = dbPool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping Postgres: %w", err)
	}

	return &PgClient{Pool: dbPool}, nil
}
