package trips_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_postgres_pool_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/errors"
)

const getTripQuery = `
SELECT
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
	deleted_at
FROM app.trips
WHERE owner_id=$1
	AND id=$2
	AND deleted_at IS NULL;
`

func (r *TripsRepository) GetTrip(
	ctx context.Context,
	ownerID int,
	tripID int,
) (*domain.Trip, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var tripResponse domain.Trip
	err := r.pool.QueryRow(
		ctx,
		getTripQuery,
		ownerID,
		tripID,
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
		if errors.Is(err, core_postgres_pool_errors.ErrNoRows) {
			return nil, fmt.Errorf(
				"get tripID %d: %w",
				tripID,
				core_errors.NewError(
					errors.Join(err, core_errors.ErrNotFound),
					core_errors.CategoryNotFound,
				),
			)
		}
		return nil, fmt.Errorf("get trip: %w", err)
	}

	return &tripResponse, nil
}
