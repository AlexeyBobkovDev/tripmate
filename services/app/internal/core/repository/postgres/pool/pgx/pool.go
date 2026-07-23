package core_pgx_pool

import (
	"context"
	"fmt"
	"time"

	core_postgres_pool "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	*pgxpool.Pool
	opTimeout time.Duration
}

func (p *Pool) Exec(ctx context.Context, sql string, arguments ...any) (core_postgres_pool.CommandTag, error) {
	cmdTag, err := p.Pool.Exec(ctx, sql, arguments...)
	return pgxCommandTag{cmdTag}, mapErrors(err)
}

func NewPool(ctx context.Context, cfg Config) (*Pool, error) {
	connUrl := cfg.BuildDSN()
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

	return &Pool{
		Pool:      pool,
		opTimeout: cfg.OpTimeout,
	}, nil
}
