package routes_postgres_repository

import (
	"time"

	core_postgres_pool "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool"
)

type RoutesRepository struct {
	pool core_postgres_pool.Pool
}

func NewRoutesRepository(
	pool core_postgres_pool.Pool,
) *RoutesRepository {
	return &RoutesRepository{
		pool: pool,
	}
}

func durationToSeconds(duration time.Duration) int64 {
	return int64(duration / time.Second)
}

func durationFromSeconds(seconds int64) time.Duration {
	return time.Duration(seconds) * time.Second
}
