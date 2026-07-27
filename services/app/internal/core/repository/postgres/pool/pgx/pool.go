package core_pgx_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	core_postgres_pool "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool"
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
	return pgxCommandTag{cmdTag}, mapErrors(err)
}

func (p *Pool) QueryRow(
	ctx context.Context,
	sql string,
	args ...any,
) core_postgres_pool.Row {
	row := p.Pool.QueryRow(ctx, sql, args...)
	return pgxRow{row}
}

type Option func(*Pool)

func WithOperationTimeout(opTimeout time.Duration) Option {
	return func(p *Pool) {
		p.opTimeout = opTimeout
	}
}

func NewPool(ctx context.Context, cfg any, opts ...Option) (*Pool, error) {
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
