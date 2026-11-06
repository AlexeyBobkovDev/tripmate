package generate_route_postgres_repository

import core_postgres_pool "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/repository/postgres/pool"

type GenerateRouteRepository struct {
	pool core_postgres_pool.Pool
}

func NewGenerateRouteRepository(
	pool core_postgres_pool.Pool,
) *GenerateRouteRepository {
	return &GenerateRouteRepository{
		pool: pool,
	}
}
