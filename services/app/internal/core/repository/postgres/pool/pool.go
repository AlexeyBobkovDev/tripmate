package core_postgres_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	*pgxpool.Pool
	opTimeout time.Time
}

func NewPool(ctx context.Context, cfg Config) (*Pool, error) {
	connUrl := cfg.BuildDSN()
	poolCfg, err := pgxpool.ParseConfig(connUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config")
	}
	poolCfg.HealthCheckPeriod = time.Minute
	poolCfg.MaxConns = 50
	poolCfg.MinConns = 10
	poolCfg.MaxConnIdleTime = 10 * time.Minute
	poolCfg.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create new postgres pool with the given config")
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to created postgres pool")
	}

	return &Pool{
		Pool:      pool,
		opTimeout: cfg.OpTimeout,
	}, nil
}
