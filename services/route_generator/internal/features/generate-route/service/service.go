package generate_route_service

import (
	"context"

	"github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/domain"
)

type GenerateRouteRepository interface {
	FindCandidates(
		ctx context.Context,
		criterial FindCandidatesCriteria,
	) ([]domain.Route, error)
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
