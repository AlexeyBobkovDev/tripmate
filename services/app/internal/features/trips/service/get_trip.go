package trips_service

import (
	"context"
	"fmt"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func (s *TripsService) GetTrip(
	ctx context.Context,
	ownerID int,
	tripID int,
) (*domain.Trip, error) {
	if ownerID <= 0 {
		return nil, fmt.Errorf("ownerID must be greater than 0: %w", core_errors.ErrBadRequest)
	}
	if tripID <= 0 {
		return nil, fmt.Errorf("tripID must be greater than 0: %w", core_errors.ErrBadRequest)
	}
	return s.tripsRepository.GetTrip(ctx, ownerID, tripID)
}
