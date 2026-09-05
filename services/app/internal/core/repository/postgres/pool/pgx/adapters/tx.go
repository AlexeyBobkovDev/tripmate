package core_pgx_pool_adapters

import (
	"context"

	"github.com/jackc/pgx/v5"

	core_postgres_pool "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool"
)

type PgxTx struct {
	pgx.Tx
}

func (tx *PgxTx) Commit(ctx context.Context) error {
	return MapErrors(tx.Tx.Commit(ctx))
}

func (tx *PgxTx) Rollback(ctx context.Context) error {
	return MapErrors(tx.Tx.Rollback(ctx))
}

func (tx *PgxTx) Exec(
	ctx context.Context,
	sql string,
	arguments ...any,
) (core_postgres_pool.CommandTag, error) {
	cmdTag, err := tx.Tx.Exec(ctx, sql, arguments...)
	if err != nil {
		return nil, MapErrors(err)
	}

	return PgxCommandTag{CommandTag: cmdTag}, nil
}

func (tx *PgxTx) QueryRow(
	ctx context.Context,
	sql string,
	args ...any,
) core_postgres_pool.Row {
	row := tx.Tx.QueryRow(ctx, sql, args...)
	return PgxRow{Row: row}
}

func (tx *PgxTx) Query(ctx context.Context, sql string, args ...any) (core_postgres_pool.Rows, error) {
	rows, err := tx.Tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, MapErrors(err)
	}
	return PgxRows{Rows: rows}, nil
}
