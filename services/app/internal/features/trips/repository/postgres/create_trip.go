package trips_postgres_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

const createTripQuery = `
INSERT INTO app.trips (
    owner_id,
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
    time_zone
)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
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

const createTripMealRequirementsQuery = `
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

func (r *TripsRepository) CreateTrip(
	ctx context.Context,
	ownerID int,
	trip *domain.TripCreate,
) (*domain.Trip, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateTrip transaction: %w", err)
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
		createTripQuery,
		ownerID,
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
		return nil, fmt.Errorf("create trip: %w", err)
	}
	tripResponse.DailySchedule = domain.DailySchedule{
		StartTime: domain.LocalTime(scannedStartTime.Format("15:04")),
		EndTime:   domain.LocalTime(scannedEndTime.Format("15:04")),
	}

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
			createTripMealRequirementsQuery,
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
			return nil, fmt.Errorf("failed to create trip meal requirements: %w", err)
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
