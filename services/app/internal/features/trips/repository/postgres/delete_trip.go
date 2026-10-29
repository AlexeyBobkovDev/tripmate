package trips_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

const deleteTripQuery = `
DELETE FROM app.trips
WHERE owner_id=$1
	AND id=$2
	AND deleted_at IS NULL;
`

func (s *TripsRepository) DeleteTrip(
	ctx context.Context,
	ownerID int,
	tripID int,
) error {
	ctx, cancel := context.WithTimeout(ctx, s.pool.OpTimeout())
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("create transaction to delete trip: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	cmdTag, err := tx.Exec(ctx, deleteTripQuery, ownerID, tripID)
	if err != nil {
		return fmt.Errorf("delete trip row: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("the trip %d does not exist: %w", tripID, core_errors.ErrNotFound)
	}

	return tx.Commit(ctx)
}
