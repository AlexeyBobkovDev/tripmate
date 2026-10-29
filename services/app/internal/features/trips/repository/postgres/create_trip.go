package trips_postgres_repository

import (
	"context"
	"fmt"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

const createTripQuery = `
INSERT INTO app.trips (
    owner_id,
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
    time_zone
)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
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

func (r *TripsRepository) CreateTrip(
	ctx context.Context,
	ownerID int,
	trip *domain.TripCreate,
) (*domain.Trip, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var tripResponse domain.Trip
	err := r.pool.QueryRow(
		ctx,
		createTripQuery,
		ownerID,
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
		return nil, fmt.Errorf("create trip: %w", err)
	}

	return &tripResponse, nil
}
