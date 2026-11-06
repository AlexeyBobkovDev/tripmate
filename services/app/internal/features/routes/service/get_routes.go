package routes_service

import (
	"context"
	"fmt"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func (s *RoutesService) GetRoutes(
	ctx context.Context,
	tripID int,
) ([]*domain.Route, error) {
	if tripID <= 0 {
		return nil, fmt.Errorf("ownerID must be greater than 0: %w", core_errors.ErrBadRequest)
	}
	return s.routesRepository.GetRoutes(ctx, tripID)
}
