package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (c *PgClient) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return c.Pool.Exec(ctx, sql, args...)
}

func (c *PgClient) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return c.Pool.Query(ctx, sql, args...)
}

func (c *PgClient) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return c.Pool.QueryRow(ctx, sql, args...)
}

func (c *PgClient) BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
	return c.Pool.BeginTx(ctx, opts)
}

func (c *PgClient) Close() {
	c.Pool.Close()
}
