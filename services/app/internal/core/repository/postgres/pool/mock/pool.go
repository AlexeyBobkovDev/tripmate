package mock

import (
	"context"

	core_postgres_pool "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool"
)

type Pool struct {
	ExecFunc func(
		ctx context.Context,
		sql string,
		arguments ...any,
	) (core_postgres_pool.CommandTag, error)

	QueryRowFunc func(
		ctx context.Context,
		sql string,
		args ...any,
	) core_postgres_pool.Row
}

func (p *Pool) Exec(
	ctx context.Context,
	sql string,
	arguments ...any,
) (core_postgres_pool.CommandTag, error) {
	return p.ExecFunc(ctx, sql, arguments...)
}

func (p *Pool) QueryRow(
	ctx context.Context,
	sql string,
	args ...any,
) core_postgres_pool.Row {
	return p.QueryRowFunc(ctx, sql, args...)
}
