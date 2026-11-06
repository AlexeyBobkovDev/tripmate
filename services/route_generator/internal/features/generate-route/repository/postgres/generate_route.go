package generate_route_postgres_repository

import (
	"context"

	"github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/domain"
)

func (r *GenerateRouteRepository) GenerateRoute(
	ctx context.Context,
	trip *domain.Trip,
) (*domain.Route, error) {
	return nil, nil
}
