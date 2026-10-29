package trips_postgres_repository

import (
	core_postgres_pool "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool"
)

type TripsRepository struct {
	pool core_postgres_pool.Pool
}

func NewTripsRepository(
	pool core_postgres_pool.Pool,
) *TripsRepository {
	return &TripsRepository{
		pool: pool,
	}
}
