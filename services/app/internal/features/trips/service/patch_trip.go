package trips_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func (s *TripsService) PatchTrip(
	ctx context.Context,
	ownerID int,
	tripID int,
	tripPatch *domain.TripPatch,
) (*domain.Trip, error) {
	if ownerID <= 0 {
		return nil, fmt.Errorf("ownerID must be greater than 0: %w", core_errors.ErrBadRequest)
	}
	if tripID <= 0 {
		return nil, fmt.Errorf("tripID must be greater than zero: %w", core_errors.ErrBadRequest)
	}

	trip, err := s.tripsRepository.GetTrip(ctx, ownerID, tripID)
	if err != nil {
		return nil, fmt.Errorf("get trip: %w", err)
	}
	if trip.Version != tripPatch.Version {
		return nil, fmt.Errorf(
			"patch trip %d: expected version %d, current version %d: %w",
			tripID,
			tripPatch.Version,
			trip.Version,
			core_errors.ErrConflict,
		)
	}
	if err := trip.ApplyPatch(tripPatch); err != nil {
		return nil, fmt.Errorf("apply trip patch: %w", errors.Join(err, core_errors.ErrBadRequest))
	}

	trip, err = s.tripsRepository.PatchTrip(ctx, ownerID, tripID, trip)
	if err != nil {
		return nil, fmt.Errorf("patch trip: %w", err)
	}

	return trip, nil
}
