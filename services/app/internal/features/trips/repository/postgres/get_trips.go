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

const getTripsQuery = `
SELECT
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
	deleted_at
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

	trips := make([]domain.Trip, 0)

	for tripRows.Next() {
		var tripResponse domain.Trip
		var startTime, endTime time.Time

		err := tripRows.Scan(
			&tripResponse.ID,
			&tripResponse.OwnerID,
			&tripResponse.Version,
			&tripResponse.Name,
			&tripResponse.PlanningArea,
			&tripResponse.StartCoordinates.Latitude,
			&tripResponse.StartCoordinates.Longitude,
			&tripResponse.StartDate,
			&tripResponse.EndDate,
			&startTime,
			&endTime,
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
		tripResponse.DailySchedule = domain.DailySchedule{
			StartTime: domain.LocalTime(startTime.Format("15:04")),
			EndTime:   domain.LocalTime(endTime.Format("15:04")),
		}

		mealRequirementRows, err := r.pool.Query(
			ctx,
			getTripMealRequirementsQuery,
			tripResponse.ID,
		)
		if err != nil {
			return nil, fmt.Errorf("get trip meal requirements: %w", err)
		}
		defer mealRequirementRows.Close()

		for mealRequirementRows.Next() {
			var (
				scannedMealRequirement           domain.MealRequirement
				scannedStartTime, scannedEndTime time.Time
			)

			err = mealRequirementRows.Scan(
				&scannedMealRequirement.Type,
				&scannedStartTime,
				&scannedEndTime,
				&scannedMealRequirement.Duration,
				&scannedMealRequirement.Budget.AmountMinorUnits,
				&scannedMealRequirement.Budget.Currency,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to create trip meal requirements: %w", err)
			}

			scannedMealRequirement.Window = domain.TimeWindow{
				Start: domain.LocalTime(scannedStartTime.Format("15:04")),
				End:   domain.LocalTime(scannedEndTime.Format("15:04")),
			}
			scannedMealRequirement.Duration *= time.Minute

			tripResponse.MealRequirements = append(tripResponse.MealRequirements, scannedMealRequirement)
		}
		if err := mealRequirementRows.Err(); err != nil {
			return nil, fmt.Errorf("failed to scan rows: %w", err)
		}

		trips = append(trips, tripResponse)
	}
	if err := tripRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan rows: %w", err)
	}

	return &trips, nil
}
