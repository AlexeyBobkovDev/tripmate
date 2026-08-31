package core_pgx_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	core_postgres_pool "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool"
	core_pgx_pool_adapters "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/pgx/adapters"
)

type Pool struct {
	*pgxpool.Pool
	opTimeout time.Duration
}

func (p *Pool) Exec(
	ctx context.Context,
	sql string,
	arguments ...any,
) (core_postgres_pool.CommandTag, error) {
	cmdTag, err := p.Pool.Exec(ctx, sql, arguments...)
	if err != nil {
		return nil, core_pgx_pool_adapters.MapErrors(err)
	}
	return core_pgx_pool_adapters.PgxCommandTag{CommandTag: cmdTag}, nil
}

func (p *Pool) QueryRow(
	ctx context.Context,
	sql string,
	args ...any,
) core_postgres_pool.Row {
	row := p.Pool.QueryRow(ctx, sql, args...)
	return core_pgx_pool_adapters.PgxRow{Row: row}
}

func (p *Pool) Query(ctx context.Context, sql string, args ...any) (core_postgres_pool.Rows, error) {
	rows, err := p.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, core_pgx_pool_adapters.MapErrors(err)
	}
	return core_pgx_pool_adapters.PgxRows{Rows: rows}, nil
}

func (p *Pool) Begin(ctx context.Context) (core_postgres_pool.Tx, error) {
	tx, err := p.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, core_pgx_pool_adapters.MapErrors(err)
	}
	return &core_pgx_pool_adapters.PgxTx{Tx: tx}, nil
}

type Option func(*Pool)

func WithOperationTimeout(opTimeout time.Duration) Option {
	return func(p *Pool) {
		p.opTimeout = opTimeout
	}
}

func NewPool(ctx context.Context, cfg any, opts ...Option) (*Pool, error) {
	for _, opt := range opts {
		if opt == nil {
			return nil, fmt.Errorf("opt must not be nil")
		}
	}

	var (
		connUrl   string
		opTimeout *time.Duration
	)

	switch cfg := cfg.(type) {
	case string:
		connUrl = cfg
	case Config:
		connUrl = cfg.BuildDSN()
		opTimeout = &cfg.OpTimeout
	default:
		panic("invalid pool config type")
	}

	poolCfg, err := pgxpool.ParseConfig(connUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}
	poolCfg.HealthCheckPeriod = time.Minute
	poolCfg.MaxConns = 50
	poolCfg.MinConns = 10
	poolCfg.MaxConnIdleTime = 10 * time.Minute
	poolCfg.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create new postgres pool with the given config: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to created postgres pool: %w", err)
	}

	poolWithTimeout := &Pool{
		Pool: pool,
	}

	for _, opt := range opts {
		opt(poolWithTimeout)
	}

	if opTimeout != nil {
		poolWithTimeout.opTimeout = *opTimeout
	}

	return poolWithTimeout, nil
}
