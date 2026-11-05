package routes_service

import (
	"context"
	"fmt"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func (s *RoutesService) DeleteRoute(
	ctx context.Context,
	tripID int,
	routeID int,
) error {
	if tripID <= 0 {
		return fmt.Errorf("tripID must be greater than 0: %w", core_errors.ErrBadRequest)
	}
	if routeID <= 0 {
		return fmt.Errorf("routeID must be greater than 0: %w", core_errors.ErrBadRequest)
	}
	return s.routesRepository.DeleteRoute(ctx, tripID, routeID)
}
