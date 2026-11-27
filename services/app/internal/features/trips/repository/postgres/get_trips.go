package trips_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_postgres_pool_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/errors"
)

const getTripsQuery = `
SELECT id
FROM app.trips
WHERE owner_id=$1
	AND deleted_at IS NULL;
`

func (r *TripsRepository) GetTrips(
	ctx context.Context,
	ownerID int,
) (*[]domain.Trip, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tripRows, err := r.pool.Query(
		ctx,
		getTripsQuery,
		ownerID,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool_errors.ErrNoRows) {
			return nil, fmt.Errorf(
				"get ownerID %d: %w",
				ownerID,
				core_errors.NewError(
					errors.Join(err, core_errors.ErrNotFound),
					core_errors.CategoryNotFound,
				),
			)
		}
		return nil, fmt.Errorf("get trip: %w", err)
	}
	defer tripRows.Close()

	var tripsIDs []int
	for tripRows.Next() {
		var tripID int
		if err := tripRows.Scan(&tripID); err != nil {
			return nil, fmt.Errorf("scan trips ids: %w", err)
		}
		tripsIDs = append(tripsIDs, tripID)
	}
	if err := tripRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan trip ids: %w", err)
	}

	trips := make([]domain.Trip, len(tripsIDs))
	for tripInd, tripID := range tripsIDs {
		tripResponse, err := getTrip(ctx, r.pool, ownerID, tripID)
		if err != nil {
			return nil, err
		}

		trips[tripInd] = *tripResponse
	}
	return &trips, nil
}
