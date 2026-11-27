package routes_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_postgres_pool "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool"
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
WHERE trip_id=$1 AND id=$2
ORDER BY id;
`

const getRouteDayQuery = `
SELECT id,
	date,
	day_number
FROM app.route_days
WHERE route_id=$1
ORDER BY day_number;
`

const getRouteStopQuery = `
SELECT id,
	place_id,
	position,
	name,
	description,
	categories,
	interests,
	address,
	website_url,
	phone,
	visit_duration_minutes,
	latitude,
	longitude,
	arrival_time,
	departure_time
FROM app.route_stops
WHERE route_day_id=$1
ORDER BY position;
`

const getRouteLegQuery = `
SELECT id,
	from_stop_id,
	to_stop_id,
	travel_mode,
	distance,
	duration
FROM app.route_legs
WHERE route_day_id=$1
ORDER BY from_stop_id;
`

const getOpeningHoursQuery = `
SELECT opens_at,
	closes_at
FROM app.opening_hours
WHERE route_stop_id=$1;
`

func (r *RoutesRepository) GetRoute(
	ctx context.Context,
	tripID int,
	routeID int,
) (*domain.Route, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	return getRoute(ctx, r.pool, tripID, routeID)
}

func getRoute(
	ctx context.Context,
	pool core_postgres_pool.Pool,
	tripID int,
	routeID int,
) (*domain.Route, error) {
	var routeResponse domain.Route

	if err := pool.QueryRow(
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
		if errors.Is(err, core_postgres_pool_errors.ErrNoRows) {
			return nil, fmt.Errorf(
				"get route %d: %w",
				routeID,
				core_errors.ErrNotFound,
			)
		}
		return nil, fmt.Errorf("failed to fetch route: %w", err)
	}

	routeDayRows, err := pool.Query(
		ctx,
		getRouteDayQuery,
		routeResponse.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch route day rows: %w", err)
	}
	defer routeDayRows.Close()

	for routeDayRows.Next() {
		var routeDayResponse domain.RouteDay
		if err := routeDayRows.Scan(
			&routeDayResponse.ID,
			&routeDayResponse.Date,
			&routeDayResponse.DayNumber,
		); err != nil {
			return nil, fmt.Errorf("failed to fetch route day: %w", err)
		}
		routeResponse.Days = append(routeResponse.Days, routeDayResponse)
	}
	if err := routeDayRows.Err(); err != nil {
		return nil, fmt.Errorf("error while fetching route day: %w", err)
	}

	for routeDaysRowsInd := range routeResponse.Days {
		routeDay, err := getRouteDay(ctx, pool, routeResponse.Days[routeDaysRowsInd])
		if err != nil {
			return nil, err
		}
		routeResponse.Days[routeDaysRowsInd] = routeDay
	}

	return &routeResponse, nil
}

func getRouteDay(
	ctx context.Context,
	pool core_postgres_pool.Pool,
	routeDay domain.RouteDay,
) (domain.RouteDay, error) {
	stops, err := getRouteStops(ctx, pool, routeDay.ID)
	if err != nil {
		return domain.RouteDay{}, fmt.Errorf("fetch route stops: %w", err)
	}

	legs, err := getRouteLegs(ctx, pool, routeDay.ID)
	if err != nil {
		return domain.RouteDay{}, fmt.Errorf("fetch route legs: %w", err)
	}

	routeDay.Stops = stops
	routeDay.Legs = legs

	return routeDay, nil
}

func getRouteStops(
	ctx context.Context,
	pool core_postgres_pool.Pool,
	routeDayID int,
) ([]domain.RouteStop, error) {
	routeStopRows, err := pool.Query(
		ctx,
		getRouteStopQuery,
		routeDayID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch route stop rows: %w", err)
	}
	defer routeStopRows.Close()

	var routeStopsResponses []domain.RouteStop

	for routeStopRows.Next() {
		var routeStopResponse domain.RouteStop

		if err := routeStopRows.Scan(
			&routeStopResponse.ID,
			&routeStopResponse.PlaceID,
			&routeStopResponse.Position,
			&routeStopResponse.Name,
			&routeStopResponse.Description,
			&routeStopResponse.Categories,
			&routeStopResponse.Interests,
			&routeStopResponse.Address,
			&routeStopResponse.WebsiteURL,
			&routeStopResponse.Phone,
			&routeStopResponse.VisitDurationMinutes,
			&routeStopResponse.Coordinates.Latitude,
			&routeStopResponse.Coordinates.Longitude,
			&routeStopResponse.ArrivalTime,
			&routeStopResponse.DepartureTime,
		); err != nil {
			return nil, fmt.Errorf("failed to fetch route stop: %w", err)
		}

		openingHours, err := getOpeningHours(ctx, pool, routeStopResponse.ID)
		if err != nil {
			return nil, err
		}
		routeStopResponse.OpeningHours = openingHours

		routeStopsResponses = append(routeStopsResponses, routeStopResponse)
	}
	if err := routeStopRows.Err(); err != nil {
		return nil, fmt.Errorf("error while fetching route day stops: %w", err)
	}

	return routeStopsResponses, nil
}

func getOpeningHours(
	ctx context.Context,
	pool core_postgres_pool.Pool,
	stopID int,
) ([]domain.OpeningHours, error) {
	openingHoursRows, err := pool.Query(
		ctx,
		getOpeningHoursQuery,
		stopID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch opening hours rows: %w", err)
	}
	defer openingHoursRows.Close()

	var openingHoursResponses []domain.OpeningHours

	for openingHoursRows.Next() {
		var openingHoursResponse domain.OpeningHours
		if err := openingHoursRows.Scan(
			&openingHoursResponse.OpensAt,
			&openingHoursResponse.ClosesAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan opening hours: %w", err)
		}
		openingHoursResponses = append(openingHoursResponses, openingHoursResponse)
	}
	if err := openingHoursRows.Err(); err != nil {
		return nil, fmt.Errorf("error while fetching route day stop opening hours: %w", err)
	}

	return openingHoursResponses, nil
}

func getRouteLegs(
	ctx context.Context,
	pool core_postgres_pool.Pool,
	routeDayID int,
) ([]domain.RouteLeg, error) {
	routeLegRows, err := pool.Query(
		ctx,
		getRouteLegQuery,
		routeDayID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch route leg rows: %w", err)
	}
	defer routeLegRows.Close()

	var routeLegResponses []domain.RouteLeg

	for routeLegRows.Next() {
		var routeLegResponse domain.RouteLeg

		var durationSeconds int64
		if err := routeLegRows.Scan(
			&routeLegResponse.ID,
			&routeLegResponse.FromStopID,
			&routeLegResponse.ToStopID,
			&routeLegResponse.TravelMode,
			&routeLegResponse.Distance,
			&durationSeconds,
		); err != nil {
			return nil, fmt.Errorf("failed to create route leg: %w", err)
		}
		routeLegResponse.Duration = durationFromSeconds(durationSeconds)

		routeLegResponses = append(routeLegResponses, routeLegResponse)
	}
	if err := routeLegRows.Err(); err != nil {
		return nil, fmt.Errorf("error while fetching route day legs: %w", err)
	}

	return routeLegResponses, nil
}
