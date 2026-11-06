package generate_route_service

import (
	"context"

	"github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/domain"
)

type GenerateRouteRepository interface {
	GenerateRoute(
		ctx context.Context,
		trip *domain.Trip,
	) (*domain.Route, error)
}
type GenerateRouteService struct {
	generateRouteRepository GenerateRouteRepository
}

func NewGenerateRouteService(
	generateRouteRepository GenerateRouteRepository,
) *GenerateRouteService {
	return &GenerateRouteService{
		generateRouteRepository: generateRouteRepository,
	}
}
