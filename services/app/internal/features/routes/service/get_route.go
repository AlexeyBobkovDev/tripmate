package routes_service

import (
	"context"
	"fmt"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func (s *RoutesService) GetRoute(
	ctx context.Context,
	tripID int,
	routeID int,
) (*domain.Route, error) {
	if tripID <= 0 {
		return nil, fmt.Errorf("tripID must be greater than 0: %w", core_errors.ErrBadRequest)
	}
	if routeID <= 0 {
		return nil, fmt.Errorf("routeID must be greater than 0: %w", core_errors.ErrBadRequest)
	}
	return s.routesRepository.GetRoute(ctx, tripID, routeID)
}
