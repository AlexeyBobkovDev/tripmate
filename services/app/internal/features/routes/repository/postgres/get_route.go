package routes_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_postgres_pool_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/errors"
)

const getRouteQuery = `
SELECT id,
	trip_id,
	trip_version,
	version,
	status,
	created_at,
	updated_at
FROM app.routes
WHERE trip_id=$1 AND id=$2;
`

const getRouteDayQuery = `
SELECT id,
	date,
	day_number
FROM app.route_days
WHERE route_id=$1;
`

const getRouteStopQuery = `
SELECT id,
	place_id,
	position,
	name,
	description,
	latitude,
	longitude,
	arrival_time,
	departure_time,
	visit_duration
FROM app.route_stops
WHERE route_day_id=$1;
`

const getRouteLegQuery = `
SELECT id,
	from_stop_id,
	to_stop_id,
	travel_mode,
	distance,
	duration
FROM app.route_legs
WHERE route_day_id=$1;
`

func (r *RoutesRepository) GetRoute(
	ctx context.Context,
	tripID int,
	routeID int,
) (*domain.Route, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("create transaction: %w", err)
	}
	defer func() { _ = tx.Commit(ctx) }()

	var routeResponse domain.Route

	if err := tx.QueryRow(
		ctx,
		getRouteQuery,
		tripID,
		routeID,
	).Scan(
		&routeResponse.ID,
		&routeResponse.TripID,
		&routeResponse.TripVersion,
		&routeResponse.Version,
		&routeResponse.Status,
		&routeResponse.CreatedAt,
		&routeResponse.UpdatedAt,
	); err != nil {
		_ = tx.Rollback(ctx)
		if errors.Is(err, core_postgres_pool_errors.ErrNoRows) {
			return nil, fmt.Errorf(
				"get route %d: %w",
				routeID,
				core_errors.ErrNotFound,
			)
		}
		return nil, fmt.Errorf("failed to fetch route: %w", err)
	}

	routeDayCounter := 0
	routeDayRows, err := tx.Query(
		ctx,
		getRouteDayQuery,
		routeResponse.ID,
	)
	if err != nil {
		_ = tx.Rollback(ctx)
		return nil, fmt.Errorf("failed to fetch route day rows: %w", err)
	}
	for routeDayRows.Next() {
		routeResponse.Days = append(routeResponse.Days, domain.RouteDay{})
		if err := routeDayRows.Scan(
			&routeResponse.Days[routeDayCounter].ID,
			&routeResponse.Days[routeDayCounter].Date,
			&routeResponse.Days[routeDayCounter].DayNumber,
		); err != nil {
			_ = tx.Rollback(ctx)
			return nil, fmt.Errorf("failed to fetch route day: %w", err)
		}
		routeDayCounter++
	}

	for routeDaysRowsInd := range routeResponse.Days {
		routeStopCounter := 0
		routeStopRows, err := tx.Query(
			ctx,
			getRouteStopQuery,
			routeResponse.Days[routeDaysRowsInd].ID,
		)
		if err != nil {
			_ = tx.Rollback(ctx)
			return nil, fmt.Errorf("failed to fetch route stop rows: %w", err)
		}
		for routeStopRows.Next() {
			routeResponse.Days[routeDaysRowsInd].Stops = append(routeResponse.Days[routeDaysRowsInd].Stops, domain.RouteStop{})

			var visitDurationSeconds int64
			if err := routeStopRows.Scan(
				&routeResponse.Days[routeDaysRowsInd].Stops[routeStopCounter].ID,
				&routeResponse.Days[routeDaysRowsInd].Stops[routeStopCounter].PlaceID,
				&routeResponse.Days[routeDaysRowsInd].Stops[routeStopCounter].Position,
				&routeResponse.Days[routeDaysRowsInd].Stops[routeStopCounter].Name,
				&routeResponse.Days[routeDaysRowsInd].Stops[routeStopCounter].Description,
				&routeResponse.Days[routeDaysRowsInd].Stops[routeStopCounter].Coordinates.Latitude,
				&routeResponse.Days[routeDaysRowsInd].Stops[routeStopCounter].Coordinates.Longitude,
				&routeResponse.Days[routeDaysRowsInd].Stops[routeStopCounter].ArrivalTime,
				&routeResponse.Days[routeDaysRowsInd].Stops[routeStopCounter].DepartureTime,
				&visitDurationSeconds,
			); err != nil {
				_ = tx.Rollback(ctx)
				return nil, fmt.Errorf("failed to fetch route stop: %w", err)
			}
			routeResponse.Days[routeDaysRowsInd].Stops[routeStopCounter].VisitDuration = durationFromSeconds(visitDurationSeconds)
			routeStopCounter++
		}

		routeLegCounter := 0
		routeLegRows, err := tx.Query(
			ctx,
			getRouteLegQuery,
			routeResponse.Days[routeDaysRowsInd].ID,
		)
		if err != nil {
			_ = tx.Rollback(ctx)
			return nil, fmt.Errorf("failed to fetch route leg rows: %w", err)
		}
		for routeLegRows.Next() {
			routeResponse.Days[routeDaysRowsInd].Legs = append(routeResponse.Days[routeDaysRowsInd].Legs, domain.RouteLeg{})

			var durationSeconds int64
			if err := routeLegRows.Scan(
				&routeResponse.Days[routeDaysRowsInd].Legs[routeLegCounter].ID,
				&routeResponse.Days[routeDaysRowsInd].Legs[routeLegCounter].FromStopID,
				&routeResponse.Days[routeDaysRowsInd].Legs[routeLegCounter].ToStopID,
				&routeResponse.Days[routeDaysRowsInd].Legs[routeLegCounter].TravelMode,
				&routeResponse.Days[routeDaysRowsInd].Legs[routeLegCounter].Distance,
				&durationSeconds,
			); err != nil {
				_ = tx.Rollback(ctx)
				return nil, fmt.Errorf("failed to create route leg: %w", err)
			}
			routeResponse.Days[routeDaysRowsInd].Legs[routeLegCounter].Duration = durationFromSeconds(durationSeconds)
			routeLegCounter++
		}
	}

	return &routeResponse, nil
}
