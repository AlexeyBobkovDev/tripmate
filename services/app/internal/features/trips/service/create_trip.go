package trips_service

import (
	"context"
	"fmt"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func (s *TripsService) CreateTrip(
	ctx context.Context,
	ownerID int,
	trip *domain.TripCreate,
) (*domain.Trip, error) {
	if ownerID <= 0 {
		return nil, fmt.Errorf("ownerID must be greater than 0: %w", core_errors.ErrBadRequest)
	}
	if err := trip.Validate(); err != nil {
		return nil, err
	}

	tripResponse, err := s.tripsRepository.CreateTrip(ctx, ownerID, trip)
	if err != nil {
		return nil, fmt.Errorf("failed to create trip: %w", err)
	}
	return tripResponse, nil
}
