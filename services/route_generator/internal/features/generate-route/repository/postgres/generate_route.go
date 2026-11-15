package generate_route_postgres_repository

import (
	"context"

	"github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/domain"
	generate_route_service "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/features/generate-route/service"
)

func (r *GenerateRouteRepository) FindCandidates(
	ctx context.Context,
	criterial generate_route_service.FindCandidatesCriteria,
) ([]domain.Route, error) {
	return nil, nil
}
