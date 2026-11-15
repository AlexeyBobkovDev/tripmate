package trips_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"time"

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
    location_latitude=$6,
    location_longitude=$7,
    start_date=$8,
    end_date=$9,
    start_time=$10,
    end_time=$11,
    amount_minor_units=$12,
    currency=$13,
    interests=$14,
    time_zone=$15
WHERE owner_id=$1 AND id=$2 AND version=$3 AND deleted_at IS NULL
RETURNING
	id,
	owner_id,
	version,
	name,
	planning_area,
	location_latitude,
	location_longitude,
	start_date,
	end_date,
	start_time,
	end_time,
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

const deleteTripMealRequirementsQuery = `
DELETE FROM app.meal_requirements
WHERE trip_id=$1;
`

const insertTripMealRequirementsQuery = `
INSERT INTO app.meal_requirements(
	trip_id,
	meal_type,
	start_time,
	end_time,
	visit_duration_minutes,
	amount_minor_units,
	currency
)
VALUES($1, $2, $3, $4, $5, $6, $7)
RETURNING meal_type,
	start_time,
	end_time,
	visit_duration_minutes,
	amount_minor_units,
	currency;
`

func (r *TripsRepository) PatchTrip(
	ctx context.Context,
	ownerID int,
	tripID int,
	trip *domain.Trip,
) (*domain.Trip, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create PatchTrip transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var tripResponse domain.Trip
	startTime, err := time.Parse("15:04", string(trip.DailySchedule.StartTime))
	if err != nil {
		return nil, core_errors.ErrInvalidTripStartTime
	}
	endTime, err := time.Parse("15:04", string(trip.DailySchedule.EndTime))
	if err != nil {
		return nil, core_errors.ErrInvalidTripEndTime
	}

	var scannedStartTime, scannedEndTime time.Time

	err = tx.QueryRow(
		ctx,
		patchTripQuery,
		ownerID,
		tripID,
		trip.Version,
		trip.Name,
		trip.PlanningArea,
		trip.StartCoordinates.Latitude,
		trip.StartCoordinates.Longitude,
		trip.StartDate,
		trip.EndDate,
		startTime,
		endTime,
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
		&tripResponse.StartCoordinates.Latitude,
		&tripResponse.StartCoordinates.Longitude,
		&tripResponse.StartDate,
		&tripResponse.EndDate,
		&scannedStartTime,
		&scannedEndTime,
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

	tripResponse.DailySchedule = domain.DailySchedule{
		StartTime: domain.LocalTime(startTime.Format("15:04")),
		EndTime:   domain.LocalTime(endTime.Format("15:04")),
	}

	// Delete all existing meal requirements for the trip
	_, err = tx.Exec(ctx, deleteTripMealRequirementsQuery, tripResponse.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete trip meal requirements: %w", err)
	}

	// Insert new meal requirements
	for _, mealRequirement := range trip.MealRequirements {
		var (
			scannedMealRequirement           domain.MealRequirement
			scannedStartTime, scannedEndTime time.Time
		)

		mealStartTime, err := time.Parse("15:04", string(mealRequirement.Window.Start))
		if err != nil {
			return nil, core_errors.ErrInvalidTripStartTime
		}
		mealEndTime, err := time.Parse("15:04", string(mealRequirement.Window.End))
		if err != nil {
			return nil, core_errors.ErrInvalidTripEndTime
		}

		err = tx.QueryRow(
			ctx,
			insertTripMealRequirementsQuery,
			tripResponse.ID,
			mealRequirement.Type,
			mealStartTime,
			mealEndTime,
			mealRequirement.Duration.Minutes(),
			mealRequirement.Budget.AmountMinorUnits,
			mealRequirement.Budget.Currency,
		).Scan(
			&scannedMealRequirement.Type,
			&scannedStartTime,
			&scannedEndTime,
			&scannedMealRequirement.Duration,
			&scannedMealRequirement.Budget.AmountMinorUnits,
			&scannedMealRequirement.Budget.Currency,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to patch trip's meal requirement: %w", err)
		}

		scannedMealRequirement.Window = domain.TimeWindow{
			Start: domain.LocalTime(scannedStartTime.Format("15:04")),
			End:   domain.LocalTime(scannedEndTime.Format("15:04")),
		}
		scannedMealRequirement.Duration *= time.Minute

		tripResponse.MealRequirements = append(tripResponse.MealRequirements, scannedMealRequirement)
	}

	_ = tx.Commit(ctx)
	return &tripResponse, nil
}
