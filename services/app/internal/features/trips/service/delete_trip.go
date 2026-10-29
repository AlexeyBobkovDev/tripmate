package trips_service

import (
	"context"
	"fmt"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func (s *TripsService) DeleteTrip(
	ctx context.Context,
	ownerID int,
	tripID int,
) error {
	if ownerID <= 0 {
		return fmt.Errorf("ownerID must be greater than 0: %w", core_errors.ErrBadRequest)
	}
	if tripID <= 0 {
		return fmt.Errorf("tripID must be greater than 0: %w", core_errors.ErrBadRequest)
	}
	return s.tripsRepository.DeleteTrip(ctx, ownerID, tripID)
}
