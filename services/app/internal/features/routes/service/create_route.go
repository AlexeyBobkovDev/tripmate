package routes_service

import (
	"context"
	"fmt"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func (s *RoutesService) CreateRoute(
	ctx context.Context,
	ownerID int,
	tripID int,
) (*domain.Route, error) {
	if ownerID <= 0 {
		return nil, fmt.Errorf("ownerID must be greater than 0: %w", core_errors.ErrBadRequest)
	}
	if tripID <= 0 {
		return nil, fmt.Errorf("tripID must be greater than 0: %w", core_errors.ErrBadRequest)
	}

	trip, err := s.tripsService.GetTrip(ctx, ownerID, tripID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve trip: %w", err)
	}

	generatedRoute, err := s.generateRouteClient.GenerateRoute(ctx, tripDomainToDTO(trip))
	if err != nil {
		return nil, fmt.Errorf("failed to generate route: %w", err)
	}

	routeResponse, err := s.routesRepository.CreateRoute(ctx, ownerID, trip, routeDomainFromDTO(generatedRoute))
	if err != nil {
		return nil, fmt.Errorf("failed to create route: %w", err)
	}
	return routeResponse, nil
}
