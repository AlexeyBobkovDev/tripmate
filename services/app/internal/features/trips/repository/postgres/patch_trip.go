package trips_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_postgres_pool_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/errors"
)

const patchTripQuery = `
UPDATE app.trips
SET
	version=version+1,
    name=$4,
    planning_area=$5,
    location_name=$6,
    location_latitude=$7,
    location_longitude=$8,
    start_date=$9,
    end_date=$10,
    amount_minor_units=$11,
    currency=$12,
    interests=$13,
    time_zone=$14
WHERE owner_id=$1 AND id=$2 AND version=$3 AND deleted_at IS NULL
RETURNING
	id,
	owner_id,
	version,
	name,
	planning_area,
	location_name,
	location_latitude,
	location_longitude,
	start_date,
	end_date,
	amount_minor_units,
	currency,
	interests,
	time_zone,
	created_at,
	updated_at,
	deleted_at;
`

const tripExistsQuery = `
SELECT EXISTS(
	SELECT 1
	FROM app.trips
	WHERE owner_id=$1 AND id=$2 AND deleted_at IS NULL
);
`

func (r *TripsRepository) PatchTrip(
	ctx context.Context,
	ownerID int,
	tripID int,
	trip *domain.Trip,
) (*domain.Trip, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var tripResponse domain.Trip
	err := r.pool.QueryRow(
		ctx,
		patchTripQuery,
		ownerID,
		tripID,
		trip.Version,
		trip.Name,
		trip.PlanningArea,
		trip.StartLocation.Name,
		trip.StartLocation.Latitude,
		trip.StartLocation.Longitude,
		trip.StartDate,
		trip.EndDate,
		trip.Budget.AmountMinorUnits,
		trip.Budget.Currency,
		trip.Interests,
		trip.TimeZone,
	).Scan(
		&tripResponse.ID,
		&tripResponse.OwnerID,
		&tripResponse.Version,
		&tripResponse.Name,
		&tripResponse.PlanningArea,
		&tripResponse.StartLocation.Name,
		&tripResponse.StartLocation.Latitude,
		&tripResponse.StartLocation.Longitude,
		&tripResponse.StartDate,
		&tripResponse.EndDate,
		&tripResponse.Budget.AmountMinorUnits,
		&tripResponse.Budget.Currency,
		&tripResponse.Interests,
		&tripResponse.TimeZone,
		&tripResponse.CreatedAt,
		&tripResponse.UpdatedAt,
		&tripResponse.DeletedAt,
	)
	if err != nil {
		if !errors.Is(err, core_postgres_pool_errors.ErrNoRows) {
			return nil, fmt.Errorf("update trip: %w", err)
		}

		var exists bool

		checkErr := r.pool.QueryRow(ctx, tripExistsQuery, ownerID, tripID).Scan(&exists)
		if checkErr != nil {
			return nil, fmt.Errorf(
				"check trip %d after patch matched no rows: %w",
				tripID,
				checkErr,
			)
		}
		if !exists {
			return nil, fmt.Errorf(
				"patch trip %d: %w",
				tripID,
				core_errors.ErrNotFound,
			)
		}

		return nil, fmt.Errorf(
			"patch trip %d: version conflict: %w",
			tripID,
			errors.Join(err, core_errors.ErrConflict),
		)
	}

	return &tripResponse, nil
}
